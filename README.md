# About Bio-OS
Bio-OS is an open-source platform for genomics users.  It provides a workspace which encapsulates data, workflows, Notebooks, job history, etc. Bio-OS provides both  GUI and CLI(Command Line Interface) to quickly setup and run  workflows specified by languages like WDL and notebook as well.

## Bio-OS Concept
![](docs/static/bioos.png)

----
## Bio-OS Workspace
Workspace is a complete encapsulation of a bioinformatics research process, including data, environment, code, operational calculation procedures, results, and dashboard as an overview. It is the basic unit that realizes executable, transportable, reproducible, shareable and publishable scientific research and biological application.
![](docs/static/workspace.png)

----
## Architecture :
In Bio-OS, there are three parts in the architecture. 
### Control layer:
1. Bioos-web: the front-end component of Bio-OS.
2. Bioctl is the command-line component of Bio-OS. Users can access Bio-OS through bioctl.
3. OpenAPI is used to describe, produce, consume, and visualize RESTful web services. It's a specification standard for REST APIs that defines the structure and syntax.
4. Bioos-apiserver: It mainly realizes the management of workspace, workflow (conforming to WES specification), data, etc. The bioos-server will send requests to the clusters of each computing plane to implement the running workflow
5. IAM(not supported in Version 1.0.0): service that helps you securely control user access
### Workflow Layer
1. Cromwell: open-source cromwell version v85
2. Storage : store workflow input ,output data or log
### Notebook Layer
1. JupyterHub:  a multi-user Hub that spawns, manages, and proxies multiple instances of the single-user Jupyter notebook server.
2. Storage : Used to store data from jupyter notebook

 ![](docs/static/arch.jpg)

## Capabilities
There are four core capabilities which support Bio-OS.
- Data Management

  Effectively organize massive scale sample data and easily access public data

- Application Management

  Bring together data、workflow、Notebooks、work engine into Sharable、 Traceable and reproducible Workspace

- Resource Management

  Autoscaling through cloud-native which supports hybrid scheduling of multi-cloud heterogeneous resources

- Interactive Analysis

  Using Built-in Notebook to do tertiary analysis

## Getting Started
### Prerequisites
To run Bio-OS, the following prerequisites must be met:
- A computer with one of the following operating systems:
  - CentOS 7.9 +
  - Ubuntu 22.04 +
- Internet access
- Resource Required

Minimum 8 cpu and 16G memory

### Installation
There will be two ways to install Bio-OS .
1. Using Docker-compose

- cromwell deployment

Bio-OS V1.0.0 Version does not support cromwell containerized deployment，You need to install java-jre 11 version on the server and configure application.conf.
```text
include required(classpath("application"))
webservice {
  port = 8000
}
workflow-options {
  workflow-log-dir = /nfs/bioos-storage/cromwell-workflow-logs
  workflow-log-temporary = false
}
call-caching {
  enabled = true
  invalidate-bad-cache-results = true
}
database {
  profile = "slick.jdbc.MySQLProfile$"
  db {
    driver = "com.mysql.cj.jdbc.Driver"
    url = "jdbc:mysql://180.184.37.106:3306/wqw?rewriteBatchedStatements=true&useSSL=false"
    port = 3306
    user = "public"
    password = "Mirrors79"
    connectionTimeout = 5000
  }
}
backend {
  default = "Local"
  providers {
    Local {
      config {
        root = "/nfs/bioos-storage/cromwell-executions"
        filesystem {
          local {
           localization: [
                  "hard-link", "soft-link", "copy"
           ]

            caching {
              duplication-strategy: [
                "hard-link", "soft-link", "copy"
              ]
              hashing-strategy: "md5"
              check-sibling-md5: false
            }
          }
        }
      }
    }
  }
}
```
**Attention**: The workflow-log-dir configuration should be consistent with the apiserver reference storage configuration. Execute the following code to complete the cromwell local deployment, and the apiserver in the subsequent docker-compose.yaml will refer to the cromwell call.
```shell
java -jar -Dconfig.file=/root/cromwell/application.conf -DLOG_LEVEL=INFO -DLOG_MODE=standard /root/cromwell/cromwell.jar server
```
- Bioos local startup

For local container environment testing, you need to install the docker container environment. Our code base provides the'docker-compose.yaml 'file.

Local environment testing can be quickly completed with the following commands.
```shell
# 进入开源 bioos 目录,执行命令
docker-compose up -d
```
**Note**: Local deployment will involve the problem of jupyterhub dynamically obtaining tokens. You can refer to the online deployment part. To update the token, you need to restart the bioos-apiserver container.

2. Online Deployment

Online deployment requires the preparation of a container orchestration engine, and we provide helm deployment packages to quickly complete bioos deployment.

- Deployment based on Kubernetes environment

If there is no local Kubernetes environment, you can try to deploy minikube. For details, please refer to the [official website](https://minikube.sigs.k8s.io/docs/start/).

When installing the container runtime, you have two options:
- Docker - The  most common container runtime environment
  - Deployment method
  ```shell
  sudo apt update && apt install  docker.io -y
  ```
- Nerdctl - Open source and open operating environment, a perfect alternative to the cncf community
  - Deployment method
  ```shell
    wget https://github.com/containerd/nerdctl/releases/download/v1.4.0/nerdctl-full-1.4.0-linux-amd64.tar.gz
    tar zxf nerdctl-full-1.4.0-linux-amd64.tar.gz -C /usr/local/
    cp /usr/local/lib/systemd/system/*.service /etc/systemd/system/
    ln -s /usr/local/bin/nerdctl /usr/bin/docker
    systemctl enable buildkit containerd
    systemctl restart buildkit containerd
  ```
Kubernetes environment also has many deployment methods, common deployment forms such as minikube/microk8s/kubeadm/kubespray, considering the simplicity and ease of use, here take minikube as an example to build the local environment, production environment recommends choosing kubespray for high availability deployment.
```shell
# 安装基础依赖
sudo apt update && apt install -y conntrack
sudo sysctl fs.protected_regular=0 #重启之后，记得执行一下，避免 minikube /tmp/juju-mkf6a06118463380f4d96c12aced04598f450743: permission denied 类似报错
mkdir -p /etc/containerd && containerd config default >  /etc/containerd/config.toml
wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.27.0/crictl-v1.27.0-linux-amd64.tar.gz
tar xvf crictl-v1.27.0-linux-amd64.tar.gz -C /usr/bin/
# 安装 minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
# 使用国内源安装 kubernetes
minikube start \
--container-runtime="containerd" \
--image-mirror-country=cn \
--image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers \
--driver=none \
--apiserver-ips='180.184.43.61' # 可选配置，如果需要外网访问，需要这里配置主机的外网IP
```
Deployment log information:
```shell
root@registry:/home/vagrant# minikube start --container-runtime="containerd" --image-mirror-country=cn --driver=none --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containeres
😄  minikube v1.30.1 on Ubuntu 22.04 (amd64)
✨  Using the none driver based on existing profile
❗  Using the 'containerd' runtime with the 'none' driver is an untested configuration!
❗  Using the 'containerd' runtime with the 'none' driver is an untested configuration!
👍  Starting control plane node minikube in cluster minikube
🔄  Restarting existing none bare metal machine for "minikube" ...
ℹ️  OS release is Ubuntu 22.04.1 LTS
📦  Preparing Kubernetes v1.26.3 on containerd 1.6.12-0ubuntu1 ...
▪ kubelet.resolv-conf=/run/systemd/resolve/resolv.conf
> kubectl.sha256:  64 B / 64 B [-------------------------] 100.00% ? p/s 0s
> kubeadm.sha256:  64 B / 64 B [-------------------------] 100.00% ? p/s 0s
> kubelet.sha256:  64 B / 64 B [-------------------------] 100.00% ? p/s 0s
> kubectl:  45.81 MiB / 45.81 MiB [-------------] 100.00% 5.61 MiB p/s 8.4s
> kubeadm:  44.61 MiB / 44.61 MiB [--------------] 100.00% 2.27 MiB p/s 20s
> kubelet:  115.65 MiB / 115.65 MiB [------------] 100.00% 4.39 MiB p/s 27s
▪ Generating certificates and keys ...
▪ Booting up control plane ...
▪ Configuring RBAC rules ...
🔗  Configuring bridge CNI (Container Networking Interface) ...
🤹  Configuring local host environment ...

❗  The 'none' driver is designed for experts who need to integrate with an existing VM
💡  Most users should use the newer 'docker' driver instead, which does not require root!
📘  For more information, see: https://minikube.sigs.k8s.io/docs/reference/drivers/none/

❗  kubectl and minikube configuration will be stored in /root
❗  To use kubectl or minikube commands as your own user, you may need to relocate them. For example, to overwrite your own settings, run:

    ▪ sudo mv /root/.kube /root/.minikube $HOME
    ▪ sudo chown -R $USER $HOME/.kube $HOME/.minikube

💡  This can also be done automatically by setting the env var CHANGE_MINIKUBE_NONE_USER=true
▪ Using image registry.cn-hangzhou.aliyuncs.com/google_containers/storage-provisioner:v5
🔎  Verifying Kubernetes components...
🌟  Enabled addons: default-storageclass, storage-provisioner
💡  kubectl not found. If you need it, try: 'minikube kubectl -- get pods -A'
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```
After installing minikube, we still need to do some small configuration, such as configuring network plugins
```shell
# 追加 kubectl  快捷方式
echo 'alias kubectl="minikube kubectl --"' >> ~/.bashrc
source ~/.bashrc
# 配置 calico cni 网络
mkdir -p /opt/cni/bin && wget https://github.com/containernetworking/plugins/releases/download/v1.3.0/cni-plugins-linux-amd64-v1.3.0.tgz
tar -xvf cni-plugins-linux-amd64-v1.3.0.tgz -C /opt/cni/bin
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.1/manifests/calico.yaml
```
Environmental inspection, the final deployment situation is as follows:
```shell
root@registry:/home/vagrant# kubectl get po -A
NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE
kube-system   calico-kube-controllers-5857bf8d58-2p25d   1/1     Running   0          7m46s
kube-system   calico-node-fk6vd                          1/1     Running   0          7m46s
kube-system   coredns-567c556887-8r8cx                   1/1     Running   0          11m
kube-system   etcd-registry                              1/1     Running   2          11m
kube-system   kube-apiserver-registry                    1/1     Running   2          12m
kube-system   kube-controller-manager-registry           1/1     Running   2          11m
kube-system   kube-proxy-vt2ks                           1/1     Running   0          11m
kube-system   kube-scheduler-registry                    1/1     Running   2          11m
kube-system   storage-provisioner                        1/1     Running   0          11m
```
Since Bioos needs to use NAS storage, we choose to use the NFS storage solution here. For convenient access, we also need to install an Ingress controller to provide network access support.

**Install the ingress controller**

Ingress deployment can refer to the official website, and you can also refer to the nginx-ingress deployment configuration of kubespray.
```shell
# 安装 ingress 控制器(国内可能镜像下载会有问题)
#kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.1/deploy/static/provider/cloud/deploy.yaml
# 如果镜像下载有问题，可以使用下面的代码替代
# kubectl apply -f https://raw.githubusercontent.com/markthink/helm-charts/main/StorageClass/base/yaml/ingress-nginx.yaml
git clone https://github.com/markthink/helm-charts
cd helm-charts && kubectl apply -f ingress_nginx/
```
Environmental inspection, the final deployment situation is as follows:
```shell
root@registry:/home/vagrant# kubectl get po -A
NAMESPACE       NAME                                       READY   STATUS    RESTARTS   AGE
ingress-nginx   ingress-nginx-controller-f6c55fdc8-x5dtz   0/1     Running   0          40s
kube-system     calico-kube-controllers-5857bf8d58-2p25d   1/1     Running   0          15m
kube-system     calico-node-fk6vd                          1/1     Running   0          15m
kube-system     coredns-567c556887-8r8cx                   1/1     Running   0          19m
kube-system     etcd-registry                              1/1     Running   2          19m
kube-system     kube-apiserver-registry                    1/1     Running   2          19m
kube-system     kube-controller-manager-registry           1/1     Running   2          19m
kube-system     kube-proxy-vt2ks                           1/1     Running   0          19m
kube-system     kube-scheduler-registry                    1/1     Running   2          19m
kube-system     storage-provisioner                        1/1     Running   0          19m
root@registry:/home/vagrant# kubectl get ingressclass -A
NAME    CONTROLLER             PARAMETERS   AGE
nginx   k8s.io/ingress-nginx   <none>       47s
```
**Install NFS Storage Service**

Use helm to install the nfs storage service (please refer to the remarks for Nfs-server deployment)
Since you can't download foreign (registry.k8s.io) images in China, you need to prepare the helm configuration file values.yaml:
```yaml
image:
  nfs:
    repository: dyrnq/nfsplugin
    tag: v4.2.0
    pullPolicy: IfNotPresent
  csiProvisioner:
    repository: dyrnq/csi-provisioner
    tag: v3.3.0
    pullPolicy: IfNotPresent
  livenessProbe:
    repository: dyrnq/livenessprobe
    tag: v2.8.0
    pullPolicy: IfNotPresent
  nodeDriverRegistrar:
    repository: dyrnq/csi-node-driver-registrar
    tag: v2.6.2
    pullPolicy: IfNotPresent
```
```shell
# 安装 helm 部署包
wget https://get.helm.sh/helm-v3.12.0-linux-amd64.tar.gz
tar xvf helm-v3.12.0-linux-amd64.tar.gz \
--strip-components=1 -C /usr/local/bin
# 默认配置一般不需要修改
# https://github.com/kubernetes-csi/csi-driver-nfs/blob/master/charts/v4.2.0/csi-driver-nfs/values.yaml
# 安装 NFS CSI 存储驱动
helm repo add csi-driver-nfs https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts
helm install csi-driver-nfs csi-driver-nfs/csi-driver-nfs --namespace kube-system --version v4.2.0 -f values.yaml
```
Use Kustomize to process related variables, [Code link](https://github.com/markthink/helm-charts/tree/main/StorageClass)：
```shell
root@opensource-bioos:~/helm-charts/storageclass# tree .
.
├── base
│   ├── kustomization.yaml
│   └── yaml
│       ├── ingress-nginx.yaml
│       └── sc.yaml
└── overlays
    ├── dev
    │   └── kustomization.yaml
    └── prod
        └── kustomization.yaml

5 directories, 5 files
```
The configuration file directory structure is as above, base is the basic resource configuration file, overlays the related variables according to different environment configuration, enter the directory and execute the following command to view the generation result. Before generating, you need to modify the kustomization.yaml file

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
# https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/patches/
# https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/patchesjson6902/
resources:
# - yaml/ingress-nginx.yaml
- yaml/sc.yaml

patches:
- patch: |-
  - op: replace
    path: /parameters/server
    value: 192.168.46.300  #这里填写 nfs server 地址
  - op: replace
    path: /parameters/share
    value: /nfs #这里填写 nfs server 根目录地址
    target:
    kind: StorageClass
    name: nfs-csi

- patch: |-
  - op: replace
    path: /spec/csi/volumeAttributes/server
    value: 192.168.46.300
  - op: replace
    path: /spec/csi/volumeAttributes/share
    value: /nfs
  - op: replace
    path: /spec/csi/volumeHandle
    value: 192.168.46.300#nfs#bioos-storage#
  - op: replace
    path: /spec/capacity/storage
    value: 50Gi
    target:
    kind: PersistentVolume
    name: bioos-storage

- patch: |-
  - op: replace
    path: /spec/resources/requests/storage
    value: 50Gi  #这里配置存储块的大小
    target:
    kind: PersistentVolumeClaim
    name: bioos-storage-pvc
```
Fill in the correct server related information and execute the command to check that it is correct
```shell
# kubectl kustomize <kustomization_directory>
cd overlays/dev && kubectl kustomize .
```
Execute the following command to complete the resource creation.
```shell
# kubectl apply -k <kustomization_directory>
cd overlays/dev && kubectl apply -k .
```

**Install Bioos Service**
- Install mysql service

Execute the following command to complete the mysql deployment
```shell
# 注意要保证 nfs-server 共享目录有 777 权限
# https://github.com/bitnami/charts/blob/main/bitnami/mysql/values.yaml
helm install mysql \
--set auth.rootPassword="admin",auth.database=bioos,auth.username=admin,auth.password=admin,global.storageClass=nfs-csi,primary.persistence.size=50Gi \
oci://registry-1.docker.io/bitnamicharts/mysql
```
- Prepare Jupyterhub/cromwell

Bioos uses helm packaging, which is currently divided into four sub-packages. It should be noted that after Jupyterhub is deployed, open the browser to obtain new tokens to continue deploying bioos-server and web services.

The corresponding Helm installation command is as follows:
```shell
# 添加 helm 仓库
helm repo add bioos https://markthink.github.io/helm-charts
helm search repo bioos
```
Search the helm repository using the command
```shell
helm search repo bioos
NAME                                 CHART VERSION        APP VERSION        DESCRIPTION
bioos/bioos                          0.1.0                v1                 BioOS UI 前端
bioos/cromwell                       0.1.0                1.0.0              A Helm chart for cromwell with local backend
bioos/jupyterhub                     2.0.0                3.0.0              Multi-user Jupyter installation
```
Install cromwell to update the values.yaml configuration file:
```yaml
## platformConfig contains information about the environment on which the Chart is being installed
## These values are expected be updated during the platform installation stage
platformConfig:
  ## Container registry for all images involved in the chart
  registryDomain: docker.io
  ## Container repository for platform components
  registryRepository: broadinstitute
  ## Platform-wide image pull policy
  imagePullPolicy: Always
  ## ImagePullSecret name for all images involved in the chart
  imagePullSecret: ""

labels: {}

podLabels: {}

podAnnotations: {}

## Ref: https://kubernetes.io/docs/user-guide/node-selection/
nodeSelector: {}

## Ref: https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
tolerations: []

image:
  name: cromwell
  tag: 85
replicaCount: 1
priorityClassName: ""
securityContext:
  enabled: false
  runAsUser: 65534 # nobody
  runAsGroup: 65534 # nobody
  fsGroup: 65534 # nobody
hardAntiAffinity: false
livenessProbe:
  enabled: true
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 10
  failureThreshold: 3
  successThreshold: 1
readinessProbe:
  enabled: true
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 10
  failureThreshold: 3
  successThreshold: 1
startupProbe:
  enabled: true
  initialDelaySeconds: 60
  periodSeconds: 30
  timeoutSeconds: 10
  failureThreshold: 10
  successThreshold: 1
resources:
  limits: {}
  requests:
    cpu: 100m
    memory: 512Mi
config:
  dir: /etc/cromwell
  file: config
service:
  type: ClusterIP
  httpPort: 8000
  # clusterIP: None
  # loadBalancerIP:
  # loadBalancerSourceRanges:
  # - 10.10.10.0/24
  annotations: {}
  labels: {}
persistence:
  enabled: true
  internal: false
  pvcName: bioos-storage-pvc
  storageClass: "nfs-csi"
  pvcSize: 50Gi
  accessMode: ReadWriteMany
log:
  path: /cromwell-workflow-logs

basePath: /nfs/bioos-storage

executionPath: /cromwell-executions

db:
  local:
    enabled: false
    path: /cromwell-db
  mysql:
    enabled: true
    host: "mysql"
    port: 3306
    name: "bioos"
    username: "root"
    password: "admin"

```
install cromwell
```shell
helm install cromwell bioos/cromwell -f values.yaml
```
Setup jupyterhub configure file values.yaml：
```yaml
imagePullSecrets: []
hub:
  baseUrl: "/jupyterhub/"
  allowNamedServers: true
  config:
    NologinAuthenticator:
      username: nobody
    JupyterHub:
      admin_access: true
      authenticator_class: myauthenticator.NologinAuthenticator
      tornado_settings:
        slow_spawn_timeout: 0
        headers:
          Access-Control-Allow-Origin: "*"
          Content-Security-Policy: "frame-ancestors 'self' http://localhost"
    Spawner:
      args:
        - "--NotebookApp.allow_origin=*"
        - '--NotebookApp.tornado_settings={"headers":{"Content-Security-Policy": "frame-ancestors ''self'' http://localhost"}}'
  image:
    name: bioos/jupyterhub
    tag: "v1.0"
  networkPolicy:
    enabled: false
  db:
    type: mysql
    url: mysql+pymysql://root:admin@mysql:3306/bioos

singleuser:
  networkTools:
    image:
      name: jupyterhub/k8s-network-tools
  image:
    name: jupyterhub/k8s-singleuser-sample
  storage:
    type: none
  cpu:
    limit: 1
    guarantee: 1
  memory:
    limit: 1G
    guarantee: 1G

scheduling:
  userScheduler:
    enabled: false
  userPlaceholder:
    enabled: false

prePuller:
  hook:
    enabled: false
  continuous:
    enabled: false

```
Execute the following commands to install jupyterhub：
```shell
helm install jupyterhub bioos/jupyterhub -f values.yaml
```
Open a browser to visit: http://serverIP/jupyterhub, and apply for a new token.

**Install Bioos Service**

You need to update the configuration of bioos. The specific path is as follows:
- You need to update the configuration of bioos. The specific path is as follows: bioos/values.yaml
  
Modify the configuration and execute the command as follows
```yaml 
image:
  repository: docker.io
  project: bioos
  web_name: web:v1.2
  apiserver_name: apiserver:v1.2
  pullPolicy: IfNotPresent

imagePullSecrets: []
replicaCount: 1
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  create: true
  annotations: {}
  name: ""

rbac:
  create: true

podAnnotations: {}
podSecurityContext: {}

securityContext: {}

service:
  type: ClusterIP
  web_port: 80
  api_http_port: 8888
  api_grpc_port: 50051

ingress:
  enabled: true
  className: ""
  annotations:
    {
    # https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#custom-max-body-size
    nginx.ingress.kubernetes.io/proxy-body-size: "64m"
    }
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: ""
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls: []

resources: {}

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}

notebook:
  images:
    - name: datascience
      version: "1.0"
      description: "this is official image"
      image: jupyter/datascience-notebook:hub-3.0.0
  resources:
    - cpu: 1
      memory: 1Gi
      disk: 20Gi
    - cpu: 2
      memory: 4Gi
      disk: 20Gi

mysql:
  hostname: mysql
  database: bioos
  username: root
  password: admin

storage:
  pvc: bioos-storage-pvc
  mountPath: /app/conf/fs

wes:
  endpoint: http://180.184.43.61:8000

jupyterhub:
  endpoint: http://180.184.43.61/jupyterhub
  adminToken: 6026738d798c495aa01c7831048539d9
```
Update config and run the following command
```shell
helm install demo bioos/bioos -f values.yaml
```
**Environmental inspection**

Execute the command'kubectl get po, pv, pvc, ing 'to view the application deployment status. Bioos uses two pieces of storage, one for Mysql and one for bioos storage, and provides several ingress rules to facilitate subpath reference services. The command status should look like this:
```shell
# kubectl get po,pv,pvc,ing
NAME                            READY   STATUS    RESTARTS   AGE
pod/apiserver-bd4f89b4c-46vgp   1/1     Running   0          147m
pod/hub-78c55bf4c8-tr8v8        1/1     Running   0          149m
pod/jupyter-nobody              1/1     Running   0          17m
pod/mysql-77b4f7f77c-b9smq      1/1     Running   0          150m
pod/web-7ff595699-rnmkf         1/1     Running   0          147m

NAME                                                        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                       STORAGECLASS   REASON   AGE
persistentvolume/bioos-storage                              10Gi       RWO            Delete           Bound    default/bioos-storage-pvc   nfs-csi                 150m
persistentvolume/pvc-6b827f42-82e5-40b2-8e74-2c56fb8d836e   10Gi       RWO            Delete           Bound    default/mysql-pv-claim      nfs-csi                 150m

NAME                                      STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistentvolumeclaim/bioos-storage-pvc   Bound    bioos-storage                              10Gi       RWO            nfs-csi        150m
persistentvolumeclaim/mysql-pv-claim      Bound    pvc-6b827f42-82e5-40b2-8e74-2c56fb8d836e   10Gi       RWO            nfs-csi        150m

NAME                                                                       CLASS   HOSTS   ADDRESS        PORTS   AGE
ingress.networking.k8s.io/demo-bioos                                       nginx   *       10.211.55.13   80      147m
ingress.networking.k8s.io/jupyter--2fjupyterhub-2f-route                   nginx   *       10.211.55.13   80      149m
ingress.networking.k8s.io/jupyter--2fjupyterhub-2fuser-2fnobody-2f-route   nginx   *       10.211.55.13   80      17m
```
Open a browser and visit http://serverIP/workspace

This completes the bioos installation.

## Tutorial
### Bio-OS user guide
Please look at the Bio-OS user guide.
### Bio-OS CLI
Please look at the CLI tutorial.

## License
This project is licensed under the Apache-2.0 License.
