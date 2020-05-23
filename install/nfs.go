// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"k8s.io/klog"
)

func NfsInstall() {
	i := &InstallConfig{
		Master0:       Hosts[0],
		EnableNfs:     EnableNfs,
		ExtendNfsAddr: ExtendNfsAddr,
		NfsPath:       NfsPath,
		DefaultSc:     DefaultSc,
	}
	if i.EnableNfs {
		i.NfsInstall()
		i.NfsDeploy()
	}
}

func (i *InstallConfig) NfsInstall() {
	if i.Master0 == i.ExtendNfsAddr || len(i.ExtendNfsAddr) == 0 {
		klog.Info("install nfs on ", i.Master0)
		nfsinstallprecmd := fmt.Sprintf(`echo '%s' > /tmp/nfs.install`, i.Template(installnfs))
		SSHConfig.Cmd(i.Master0, nfsinstallprecmd)
		nfsinstallcmd := fmt.Sprintf("bash -x /tmp/nfs.install")
		SSHConfig.Cmd(i.Master0, nfsinstallcmd)
	}
}

func (i *InstallConfig) NfsDeploy() {
	klog.Info("deploy nfs to k8s ", i.Master0)
	scingcmd := fmt.Sprintf(`echo '%s' | kubectl apply -f -`, i.Template(ScDefault))
	SSHConfig.Cmd(i.Master0, scingcmd)
}

const installnfs = `
install_debian() {
# 安装
apt update
apt install -y nfs-kernel-server
# 配置
mkdir {{.NfsPath}}
echo "{{.NfsPath}} *(insecure,rw,sync,no_root_squash,no_subtree_check)" > /etc/exports
# 启动nfs
systemctl enable rpcbind
systemctl enable nfs-server
systemctl start rpcbind
systemctl start nfs-server
exportfs -r
# 测试
showmount -e 127.0.0.1
}

install_centos() {
# 安装nfs
yum install -y nfs-utils
# 配置共享目录
mkdir {{.NfsPath}}
echo "{{.NfsPath}} *(insecure,rw,sync,no_root_squash)" > /etc/exports
# 启动nfs
systemctl enable rpcbind
systemctl enable nfs-server

systemctl start rpcbind
systemctl start nfs-server
exportfs -r
# 测试
showmount -e 127.0.0.1
}
which apt && install_debian || install_centos
`

const ScDefault = `
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: eip-nfs-client-provisioner
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: eip-nfs-client-provisioner-runner
  namespace: kube-system
rules:
  - apiGroups:
      - ""
    resources:
      - persistentvolumes
    verbs:
      - get
      - list
      - watch
      - create
      - delete
  - apiGroups:
      - ""
    resources:
      - persistentvolumeclaims
    verbs:
      - get
      - list
      - watch
      - update
  - apiGroups:
      - storage.k8s.io
    resources:
      - storageclasses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - update
      - patch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: eip-run-nfs-client-provisioner
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: eip-nfs-client-provisioner-runner
subjects:
  - kind: ServiceAccount
    name: eip-nfs-client-provisioner
    namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: eip-leader-locking-nfs-client-provisioner
  namespace: kube-system
rules:
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - get
      - list
      - watch
      - create
      - update
      - patch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: eip-leader-locking-nfs-client-provisioner
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: eip-leader-locking-nfs-client-provisioner
subjects:
  - kind: ServiceAccount
    name: eip-nfs-client-provisioner
    namespace: kube-system

---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nfs-service
  name: nfs-service
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nfs-service
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: nfs-service
    spec:
      containers:
        - env:
            - name: PROVISIONER_NAME
              value: nfs-provisioner
            - name: NFS_SERVER
              value: {{.ExtendNfsAddr}}
            - name: NFS_PATH
              value: {{.NfsPath}}
          image: ysicing/nfs-client-provisioner
          name: nfs-client-provisioner
          volumeMounts:
            - mountPath: /persistentvolumes
              name: nfs-client-root
      serviceAccountName: eip-nfs-client-provisioner
      volumes:
        - name: nfs-client-root
          nfs:
            path: {{.NfsPath}}
            server: {{.ExtendNfsAddr}}
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    k8s.eip.work/storageType: nfs_client_provisioner
    storageclass.kubernetes.io/is-default-class: "true"
  name: {{.DefaultSc}}
parameters:
  archiveOnDelete: "false"
provisioner: nfs-provisioner
reclaimPolicy: Delete
volumeBindingMode: Immediate
`
