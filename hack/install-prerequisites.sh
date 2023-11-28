# This script in all prerequisites for installing bioos using helm
# Reference: https://bio-os.gitbook.io/userguide/bu-shu/getting-set-up/bu-shu-bioos
# Usage: `hack/install-prerequisites.sh`.

minikube start \
    --image-mirror-country=cn \
    --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers \
    --force \
    --addons ingress \
    --cni flannel \
    --force \
    --nodes 1 \
    --ports 80,443 \
    --service-cluster-ip-range '10.255.0.0/16' \
    --subnet '10.254.0.0/16' # 可执行修改

# set alias for kubectl
echo 'alias kubectl="minikube kubectl -- "' >> ~/.bashrc
source ~/.bashrc

# install Helm
wget https://get.helm.sh/helm-v3.12.0-linux-amd64.tar.gz
tar xvf helm-v3.12.0-linux-amd64.tar.gz --strip-components=1 -C /usr/local/bin

# install NFS Server
function installNFSServer() {
    apt update && apt install -y nfs-kernel-server nfs-common
}

# Config NFS Server
function configNFS() {
    NFS_PATH=/nfs
    mkdir -p ${NFS_PATH}/bioos-storage && chmod -R 777 ${NFS_PATH}
    chown -R nobody:nogroup ${NFS_PATH}
    echo "${NFS_PATH} *(insecure,rw,sync,root_squash,no_subtree_check,all_squash)" >> /etc/exports
    systemctl restart nfs-server
    systemctl restart rpcbind
    exportfs -arv
    showmount -e localhost
}

# Set up auto-start on boot.
function autoStart() {
    systemctl enable nfs-server
    systemctl enable rpcbind
}

# install, config and start NFS server
function main_nf() {
    installNFSServer
    configNFS
    autoStart
}

main_nfs

# 添加 nfs-csi helm repo
helm repo add csi-driver-nfs https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts

# 创建 NFS Secret
kubectl -n kube-system create secret generic mount-options --from-literal mountOptions="nfsvers=3,hard"

# install nfs-csi helm chart
helm install csi-driver-nfs csi-driver-nfs/csi-driver-nfs \
        --namespace kube-system \
        --version v4.4.0 \
        --set image.nfs.repository=dyrnq/nfsplugin \
        --set image.csiProvisioner.repository=dyrnq/csi-provisioner \
        --set image.csiSnapshotter.repository=dyrnq/csi-snapshotter \
        --set image.livenessProbe.repository=dyrnq/livenessprobe \
        --set image.nodeDriverRegistrar.repository=dyrnq/csi-node-driver-registrar \
        --set image.externalSnapshotter.repository=dyrnq/snapshot-controller
