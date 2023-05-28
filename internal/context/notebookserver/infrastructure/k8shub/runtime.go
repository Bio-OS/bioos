package k8shub

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"path"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/pkg/jupyterhub"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/storage"
	"github.com/Bio-OS/bioos/pkg/utils"
)

const (
	jupyterHomePath = "/home/jovyan"
)

var (
	mapGPUVendorResourceName = map[string]string{
		notebook.GPUVendorNvidia: "nvidia.com/gpu",
	}
)

type runtime struct {
	client             *jupyterhub.API
	hubAddress         string
	storageOpts        *storage.Options
	kubeClient         *kubernetes.Clientset
	storageClass       string
	namespace          string
	allowHomeExpansion bool
}

func NewRuntime(ctx context.Context, conf *notebook.JupyterhubConfig, storageOpts *storage.Options) (domain.Runtime, error) {
	if conf.Kubernetes == nil {
		return nil, fmt.Errorf("k8s hub required kubernetes config")
	}

	// new kube client
	kubeConfig, err := clientcmd.BuildConfigFromFlags(conf.Kubernetes.MasterURL, conf.Kubernetes.KubeconfigPath)
	if err != nil {
		log.Warnf("in cluster config init fail, try %s", clientcmd.RecommendedHomeFile)
		if kubeConfig, err = clientcmd.BuildConfigFromFlags(conf.Kubernetes.MasterURL, clientcmd.RecommendedHomeFile); err != nil {
			return nil, fmt.Errorf("get kubeconfig fail: %w", err)
		}
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("new kubernetes client fail: %w", err)
	}

	// check storage class
	allowHomeExpansion := false
	if conf.Kubernetes.StorageClass != "" {
		sc, err := kubeClient.StorageV1().StorageClasses().Get(ctx, conf.Kubernetes.StorageClass, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("get storageclass %s fail: %w", conf.Kubernetes.StorageClass, err)
		}
		if sc.AllowVolumeExpansion != nil {
			allowHomeExpansion = *sc.AllowVolumeExpansion
		}
	}

	// check PVC
	if storageOpts.FileSystem != nil && storageOpts.FileSystem.KubeResource != nil && storageOpts.FileSystem.KubeResource.PVCName != "" {
		if storageOpts.FileSystem.KubeResource.Namespace != conf.Kubernetes.Namespace {
			return nil, fmt.Errorf(
				"pvc must in jupyterhub namespace '%s', but now in '%s'",
				conf.Kubernetes.Namespace, storageOpts.FileSystem.KubeResource.Namespace)
		}
		err = checkSharePVC(ctx, kubeClient, storageOpts.FileSystem.KubeResource.Namespace, storageOpts.FileSystem.KubeResource.PVCName)
		if err != nil {
			return nil, err
		}
	}

	return &runtime{
		client:             jupyterhub.NewAPI(conf.Endpoint, conf.AdminToken, nil),
		hubAddress:         conf.Endpoint,
		storageOpts:        storageOpts,
		kubeClient:         kubeClient,
		storageClass:       conf.Kubernetes.StorageClass,
		namespace:          conf.Kubernetes.Namespace,
		allowHomeExpansion: allowHomeExpansion,
	}, nil
}

func (r *runtime) Create(ctx context.Context, srv *domain.NotebookServer) error {
	username := getHubUsername(srv)
	_, err := r.client.GetUser(ctx, username)
	if err != nil {
		if err != jupyterhub.ErrorNotFound {
			return fmt.Errorf("get hub user '%s' info fail: %w", username, err)
		}
		if err = r.client.CreateUser(ctx, username); err != nil {
			if err == jupyterhub.ErrorConflict {
				log.Warnf("jupyterhub user '%s' exist so no need create", username)
				return nil
			}
			return fmt.Errorf("create hub user '%s' fail: %w", username, err)
		}
	}
	if r.storageClass != "" {
		// TODO create dynamic home volume by myself
		_, err = r.kubeClient.CoreV1().PersistentVolumeClaims(r.namespace).Create(ctx, r.generatePVC(srv), metav1.CreateOptions{})
		if err != nil {
			if err2 := r.client.DeleteUser(ctx, username); err2 != nil {
				log.Errorf("revert jupyterhub useradd '%s' fail: %s", username, err2)
			}
			return fmt.Errorf("create home storage pvc fail: %w", err)
		}
	}
	return nil
}

func (r *runtime) Start(ctx context.Context, srv *domain.NotebookServer) error {
	options, err := r.newUserOptions(srv)
	if err != nil {
		return fmt.Errorf("invalid notebook server settings: %w", err)
	}
	username := getHubUsername(srv)
	user, err := r.client.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("get hub user '%s' info fail: %w", username, err)
	}
	servername := getServerName(srv)
	jupyterServer, ok := user.Servers[servername]
	if ok {
		status := judgeStatus(&jupyterServer)
		if status != domain.ServerStatusTerminated && status != domain.ServerStatusUnknown {
			log.Warnf("jupyter server '%s' is %s and can not start", servername, status)
			return nil
		}
	}
	// TODO check if disk size modified and storageclass allowVolumeExpansion must be true
	if err = r.client.StartServer(ctx, username, servername, options); err != nil {
		return fmt.Errorf("start server '%s' fail: %w", servername, err)
	}
	return nil
}

func (r *runtime) Stop(ctx context.Context, srv *domain.NotebookServer) error {
	username := getHubUsername(srv)
	user, err := r.client.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("get hub user '%s' info fail: %w", username, err)
	}

	servername := getServerName(srv)
	jupyterServer, ok := user.Servers[servername]
	if !ok {
		log.Warnf("jupyter server '%s' no exist", servername)
		return nil
	}
	status := judgeStatus(&jupyterServer)
	if status != domain.ServerStatusRunning {
		return fmt.Errorf("jupyter server '%s' is %s and can not stop", servername, status)
	}
	return r.client.StopServer(ctx, username, servername)
}

func (r *runtime) Delete(ctx context.Context, srv *domain.NotebookServer) error {
	username := getHubUsername(srv)
	servername := getServerName(srv)
	_, err := r.client.GetUser(ctx, username)
	if err != nil {
		if err == jupyterhub.ErrorNotFound {
			log.Warnf("jupyterhub user '%s' no found", username)
			return nil
		}
		return fmt.Errorf("get hub user '%s' info fail: %w", username, err)
	}
	if err = r.client.DeleteServer(ctx, username, servername); err == jupyterhub.ErrorNotFound {
		log.Warnf("jupyter server '%s' no exist", servername)
		return nil
	}
	// TODO make sure dynamic PVC deleted
	return err
}

func (r *runtime) GetStatus(ctx context.Context, srv *domain.NotebookServer) (*domain.Status, error) {
	username := getHubUsername(nil)
	user, err := r.client.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get hub user '%s' info fail: %w", username, err)
	}

	var res domain.Status
	servername := getServerName(srv)
	jupyterServer, ok := user.Servers[servername]
	if ok {
		res.Status = judgeStatus(&jupyterServer)
		if res.Status == domain.ServerStatusRunning {
			res.AccessURL = jupyterServer.URL
		} else if res.Status == domain.ServerStatusPending {
			u, _ := url.Parse(r.hubAddress)
			accessURL := u.JoinPath("hub", "spawn", username, servername)
			queryParam := url.Values{}
			queryParam.Set("next", jupyterServer.URL)
			accessURL.RawQuery = queryParam.Encode()
			res.AccessURL = accessURL.String()
		}
	} else {
		log.Warnf("jupyter server '%s' no found", servername)
		res.Status = domain.ServerStatusUnknown
	}
	return &res, nil
}

func (r *runtime) generatePVC(srv *domain.NotebookServer) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getHomePVCName(srv),
			Namespace:   r.namespace,
			Annotations: map[string]string{},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: &r.storageClass,
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: *resource.NewQuantity(srv.Settings.ResourceSize.Disk, resource.BinarySI),
				},
			},
		},
	}
}

func judgeStatus(srv *jupyterhub.Server) (status string) {
	if srv.Ready {
		status = domain.ServerStatusRunning
	} else if srv.Pending != nil {
		switch *srv.Pending {
		case "spawn":
			status = domain.ServerStatusPending
		case "stop":
			status = domain.ServerStatusTerminating
		default:
			log.Warnf("unknown jupyter server '%s' pending value '%s'", srv.Name, *srv.Pending)
			status = domain.ServerStatusUnknown
		}
	} else {
		status = domain.ServerStatusTerminated
	}
	return
}

func (r *runtime) newUserOptions(srv *domain.NotebookServer) (*userOptions, error) {
	cpu := int64(math.Ceil(srv.Settings.ResourceSize.CPU))
	res := &userOptions{
		Profile: profile{
			KubespawnerOverride: kubespawnerOverride{
				Image:                   srv.Settings.DockerImage,
				CPULimit:                cpu,
				CPUGuarantee:            cpu,
				MemLimit:                srv.Settings.ResourceSize.Memory,
				MemGuarantee:            srv.Settings.ResourceSize.Memory,
				ExtraResourceLimits:     map[string]string{}, // must be empty dict {} to insure override
				ExtraResourceGuarantees: map[string]string{},
				NodeSelector:            map[string]string{},
				Volumes:                 []corev1.Volume{},
				VolumeMounts:            []corev1.VolumeMount{},
			},
		},
	}

	if srv.Settings.ResourceSize.GPU != nil {
		gpu := srv.Settings.ResourceSize.GPU
		if gpu.Card != math.Trunc(gpu.Card) {
			return nil, fmt.Errorf("GPU card %f must be integer in k8s", gpu.Card)
		}
		card := strconv.FormatInt(int64(math.Ceil(gpu.Card)), 10)
		resource := mapGPUVendorResourceName[gpu.Vendor()]
		res.Profile.KubespawnerOverride.ExtraResourceLimits[resource] = card
		res.Profile.KubespawnerOverride.ExtraResourceGuarantees[resource] = card
		if len(srv.Settings.NodeSelector) > 0 {
			res.Profile.KubespawnerOverride.NodeSelector = srv.Settings.NodeSelector
		}
	}

	if err := r.setKubespawnerOverrideVolume(srv, &res.Profile.KubespawnerOverride); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *runtime) setKubespawnerOverrideVolume(srv *domain.NotebookServer, override *kubespawnerOverride) error {
	pvcNameSet := map[string]string{} // pvc name -> volume name in pod
	getVolumeName := func(pvcName, name string) string {
		v, ok := pvcNameSet[pvcName]
		if ok {
			return v
		}
		pvcNameSet[pvcName] = name
		return name
	}

	// set home dir
	const homeVolumeName = "home"
	if r.storageClass != "" {
		override.VolumeMounts = append(override.VolumeMounts, corev1.VolumeMount{
			Name:      getVolumeName(getHomePVCName(srv), homeVolumeName),
			MountPath: jupyterHomePath,
		})
	} else {
		override.VolumeMounts = append(override.VolumeMounts, corev1.VolumeMount{
			Name:      getVolumeName(r.storageOpts.FileSystem.KubeResource.PVCName, homeVolumeName),
			MountPath: jupyterHomePath,
			SubPath:   "notebookserver/" + srv.ID,
		})
	}

	// set extra volume
	for _, volume := range srv.Volumes {
		if volume.Type == domain.VolumeTypeNFS {
			if r.storageOpts.FileSystem == nil || r.storageOpts.FileSystem.KubeResource == nil {
				return fmt.Errorf("nfs k8s resource no specified in storage options")
			}
			subpath, ok := utils.GetSubPath(r.storageOpts.FileSystem.RootPath, volume.Source)
			if !ok {
				return fmt.Errorf("nfs path '%s' no specified in storage options", volume.Source)
			}
			override.VolumeMounts = append(override.VolumeMounts, corev1.VolumeMount{
				Name:      getVolumeName(r.storageOpts.FileSystem.KubeResource.PVCName, volume.Name),
				MountPath: path.Join(jupyterHomePath, volume.MountRelativePath),
				SubPath:   subpath,
			})
		} else {
			return fmt.Errorf("volume type %s not support", volume.Type)
		}
	}

	for k, v := range pvcNameSet {
		override.Volumes = append(override.Volumes, corev1.Volume{
			Name: v,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: k,
				},
			},
		})
	}
	return nil
}

func getHomePVCName(srv *domain.NotebookServer) string {
	return getHubUsername(srv) + "-" + srv.ID
}

func getHubUsername(srv *domain.NotebookServer) string {
	// TODO fix it when multi user developed
	return "nobody"
}

func getServerName(srv *domain.NotebookServer) string {
	return srv.ID
}

func checkSharePVC(ctx context.Context, kubeClient *kubernetes.Clientset, namespace, name string) error {
	pvc, err := kubeClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get pvc %s/%s fail: %w", namespace, name, err)
	}
	if pvc.Status.Phase != corev1.ClaimBound {
		return fmt.Errorf("pvc %s/%s phase %s, must be Bound", namespace, name, pvc.Status.Phase)
	}
	hasReadWriteMany := false
	for _, m := range pvc.Status.AccessModes {
		if m == corev1.ReadWriteMany {
			hasReadWriteMany = true
		}
	}
	if !hasReadWriteMany {
		log.Warnf("pvc %s/%s is not ReadWriteMany", namespace, name)
	}
	return nil
}
