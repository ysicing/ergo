// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"k8s.io/klog"
	"strings"
)

func K8sInstall() {
	i := &InstallConfig{
		Hosts:               Hosts,
		Masters:             Masters,
		Wokers:              Wokers,
		EnableIngress:       EnableIngress,
		EnableNfs:           EnableNfs,
		ExtendNfsAddr:       ExtendNfsAddr,
		EnableKuboard:       EnableKuboard,
		EnableMetricsServer: EnableMetricsServer,
		NfsPath:             NfsPath,
		DefaultSc:           DefaultSc,
		Master0:             strings.Split(Masters, "-")[0],
	}
	i.K8sInstall()
	if i.EnableIngress {
		i.IngressInstall()
	}
	if i.EnableNfs {
		i.NfsInstall()
		i.NfsDeploy()
	}

	if i.EnableMetricsServer {
		i.MetricsServerDeploy()
	}

	if i.EnableKuboard {
		i.KuboardInstall()
		i.KuboardDone()
	}
}

func (i *InstallConfig) K8sInstall() {

	var k8scmd string
	if len(Wokers) == 0 {
		k8scmd = fmt.Sprintf("docker run -v /root:/root -v /etc/kubernetes:/etc/kubernetes --rm -e MASTER_IP=%s -e PASS=%s  ysicing/k7s", Masters, SSHConfig.Password)
	} else {
		k8scmd = fmt.Sprintf("docker run -v /root:/root -v /etc/kubernetes:/etc/kubernetes --rm -e MASTER_IP=%s -e NODE_IP=%s -e PASS=%s  ysicing/k7s", Masters, Wokers, SSHConfig.Password)
	}

	SSHConfig.Cmd(i.Master0, k8scmd)
	SSHConfig.Cmd(i.Master0, "kubectl taint nodes --all node-role.kubernetes.io/master- || sleep 1") // 允许master节点调度
	SSHConfig.Cmd(i.Master0, "helminit")

}

func (i *InstallConfig) IngressInstall() {
	ncingcmd := fmt.Sprintf(`echo '%s' | kubectl apply -f -`, NcIngress)
	SSHConfig.Cmd(i.Master0, ncingcmd)
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

const NcIngress = `
apiVersion: v1
kind: Namespace
metadata:
  name: nginx-ingress

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx-ingress 
  namespace: nginx-ingress

---
apiVersion: v1
kind: Secret
metadata:
  name: default-server-secret
  namespace: nginx-ingress
type: Opaque
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN2akNDQWFZQ0NRREFPRjl0THNhWFhEQU5CZ2txaGtpRzl3MEJBUXNGQURBaE1SOHdIUVlEVlFRRERCWk8KUjBsT1dFbHVaM0psYzNORGIyNTBjbTlzYkdWeU1CNFhEVEU0TURreE1qRTRNRE16TlZvWERUSXpNRGt4TVRFNApNRE16TlZvd0lURWZNQjBHQTFVRUF3d1dUa2RKVGxoSmJtZHlaWE56UTI5dWRISnZiR3hsY2pDQ0FTSXdEUVlKCktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUwvN2hIUEtFWGRMdjNyaUM3QlBrMTNpWkt5eTlyQ08KR2xZUXYyK2EzUDF0azIrS3YwVGF5aGRCbDRrcnNUcTZzZm8vWUk1Y2Vhbkw4WGM3U1pyQkVRYm9EN2REbWs1Qgo4eDZLS2xHWU5IWlg0Rm5UZ0VPaStlM2ptTFFxRlBSY1kzVnNPazFFeUZBL0JnWlJVbkNHZUtGeERSN0tQdGhyCmtqSXVuektURXUyaDU4Tlp0S21ScUJHdDEwcTNRYzhZT3ExM2FnbmovUWRjc0ZYYTJnMjB1K1lYZDdoZ3krZksKWk4vVUkxQUQ0YzZyM1lma1ZWUmVHd1lxQVp1WXN2V0RKbW1GNWRwdEMzN011cDBPRUxVTExSakZJOTZXNXIwSAo1TmdPc25NWFJNV1hYVlpiNWRxT3R0SmRtS3FhZ25TZ1JQQVpQN2MwQjFQU2FqYzZjNGZRVXpNQ0F3RUFBVEFOCkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQWpLb2tRdGRPcEsrTzhibWVPc3lySmdJSXJycVFVY2ZOUitjb0hZVUoKdGhrYnhITFMzR3VBTWI5dm15VExPY2xxeC9aYzJPblEwMEJCLzlTb0swcitFZ1U2UlVrRWtWcitTTFA3NTdUWgozZWI4dmdPdEduMS9ienM3bzNBaS9kclkrcUI5Q2k1S3lPc3FHTG1US2xFaUtOYkcyR1ZyTWxjS0ZYQU80YTY3Cklnc1hzYktNbTQwV1U3cG9mcGltU1ZmaXFSdkV5YmN3N0NYODF6cFErUyt1eHRYK2VBZ3V0NHh3VlI5d2IyVXYKelhuZk9HbWhWNThDd1dIQnNKa0kxNXhaa2VUWXdSN0diaEFMSkZUUkk3dkhvQXprTWIzbjAxQjQyWjNrN3RXNQpJUDFmTlpIOFUvOWxiUHNoT21FRFZkdjF5ZytVRVJxbStGSis2R0oxeFJGcGZnPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBdi91RWM4b1JkMHUvZXVJTHNFK1RYZUprckxMMnNJNGFWaEMvYjVyYy9XMlRiNHEvClJOcktGMEdYaVN1eE9ycXgrajlnamx4NXFjdnhkenRKbXNFUkJ1Z1B0ME9hVGtIekhvb3FVWmcwZGxmZ1dkT0EKUTZMNTdlT1l0Q29VOUZ4amRXdzZUVVRJVUQ4R0JsRlNjSVo0b1hFTkhzbysyR3VTTWk2Zk1wTVM3YUhudzFtMApxWkdvRWEzWFNyZEJ6eGc2clhkcUNlUDlCMXl3VmRyYURiUzc1aGQzdUdETDU4cGszOVFqVUFQaHpxdmRoK1JWClZGNGJCaW9CbTVpeTlZTW1hWVhsMm0wTGZzeTZuUTRRdFFzdEdNVWozcGJtdlFmazJBNnljeGRFeFpkZFZsdmwKMm82MjBsMllxcHFDZEtCRThCay90elFIVTlKcU56cHpoOUJUTXdJREFRQUJBb0lCQVFDZklHbXowOHhRVmorNwpLZnZJUXQwQ0YzR2MxNld6eDhVNml4MHg4Mm15d1kxUUNlL3BzWE9LZlRxT1h1SENyUlp5TnUvZ2IvUUQ4bUFOCmxOMjRZTWl0TWRJODg5TEZoTkp3QU5OODJDeTczckM5bzVvUDlkazAvYzRIbjAzSkVYNzZ5QjgzQm9rR1FvYksKMjhMNk0rdHUzUmFqNjd6Vmc2d2szaEhrU0pXSzBwV1YrSjdrUkRWYmhDYUZhNk5nMUZNRWxhTlozVDhhUUtyQgpDUDNDeEFTdjYxWTk5TEI4KzNXWVFIK3NYaTVGM01pYVNBZ1BkQUk3WEh1dXFET1lvMU5PL0JoSGt1aVg2QnRtCnorNTZud2pZMy8yUytSRmNBc3JMTnIwMDJZZi9oY0IraVlDNzVWYmcydVd6WTY3TWdOTGQ5VW9RU3BDRkYrVm4KM0cyUnhybnhBb0dCQU40U3M0ZVlPU2huMVpQQjdhTUZsY0k2RHR2S2ErTGZTTXFyY2pOZjJlSEpZNnhubmxKdgpGenpGL2RiVWVTbWxSekR0WkdlcXZXaHFISy9iTjIyeWJhOU1WMDlRQ0JFTk5jNmtWajJTVHpUWkJVbEx4QzYrCk93Z0wyZHhKendWelU0VC84ajdHalRUN05BZVpFS2FvRHFyRG5BYWkyaW5oZU1JVWZHRXFGKzJyQW9HQkFOMVAKK0tZL0lsS3RWRzRKSklQNzBjUis3RmpyeXJpY05iWCtQVzUvOXFHaWxnY2grZ3l4b25BWlBpd2NpeDN3QVpGdwpaZC96ZFB2aTBkWEppc1BSZjRMazg5b2pCUmpiRmRmc2l5UmJYbyt3TFU4NUhRU2NGMnN5aUFPaTVBRHdVU0FkCm45YWFweUNweEFkREtERHdObit3ZFhtaTZ0OHRpSFRkK3RoVDhkaVpBb0dCQUt6Wis1bG9OOTBtYlF4VVh5YUwKMjFSUm9tMGJjcndsTmVCaWNFSmlzaEhYa2xpSVVxZ3hSZklNM2hhUVRUcklKZENFaHFsV01aV0xPb2I2NTNyZgo3aFlMSXM1ZUtka3o0aFRVdnpldm9TMHVXcm9CV2xOVHlGanIrSWhKZnZUc0hpOGdsU3FkbXgySkJhZUFVWUNXCndNdlQ4NmNLclNyNkQrZG8wS05FZzFsL0FvR0FlMkFVdHVFbFNqLzBmRzgrV3hHc1RFV1JqclRNUzRSUjhRWXQKeXdjdFA4aDZxTGxKUTRCWGxQU05rMXZLTmtOUkxIb2pZT2pCQTViYjhibXNVU1BlV09NNENoaFJ4QnlHbmR2eAphYkJDRkFwY0IvbEg4d1R0alVZYlN5T294ZGt5OEp0ek90ajJhS0FiZHd6NlArWDZDODhjZmxYVFo5MWpYL3RMCjF3TmRKS2tDZ1lCbyt0UzB5TzJ2SWFmK2UwSkN5TGhzVDQ5cTN3Zis2QWVqWGx2WDJ1VnRYejN5QTZnbXo5aCsKcDNlK2JMRUxwb3B0WFhNdUFRR0xhUkcrYlNNcjR5dERYbE5ZSndUeThXczNKY3dlSTdqZVp2b0ZpbmNvVlVIMwphdmxoTUVCRGYxSjltSDB5cDBwWUNaS2ROdHNvZEZtQktzVEtQMjJhTmtsVVhCS3gyZzR6cFE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-config
  namespace: nginx-ingress
data:
  server-names-hash-bucket-size: "1024"


---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: nginx-ingress
rules:
- apiGroups:
  - ""
  resources:
  - services
  - endpoints
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
  - update
  - create
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - list
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
- apiGroups:
  - extensions
  resources:
  - ingresses
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - "extensions"
  resources:
  - ingresses/status
  verbs:
  - update
- apiGroups:
  - k8s.nginx.org
  resources:
  - virtualservers
  - virtualserverroutes
  verbs:
  - list
  - watch
  - get

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: nginx-ingress
subjects:
- kind: ServiceAccount
  name: nginx-ingress
  namespace: nginx-ingress
roleRef:
  kind: ClusterRole
  name: nginx-ingress
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nginx-ingress
  namespace: nginx-ingress
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "9113"
spec:
  selector:
    matchLabels:
      app: nginx-ingress
  template:
    metadata:
      labels:
        app: nginx-ingress
    spec:
      serviceAccountName: nginx-ingress
      containers:
      - image: nginx/nginx-ingress:1.6.3
        name: nginx-ingress
        ports:
        - name: http
          containerPort: 80
          hostPort: 80
        - name: https
          containerPort: 443
          hostPort: 443
        - name: prometheus
          containerPort: 9113
        securityContext:
          allowPrivilegeEscalation: true
          runAsUser: 101 #nginx
          capabilities:
            drop:
            - ALL
            add:
            - NET_BIND_SERVICE
        env:
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        args:
          - -nginx-configmaps=$(POD_NAMESPACE)/nginx-config
          - -default-server-tls-secret=$(POD_NAMESPACE)/default-server-secret
         #- -v=3 # Enables extensive logging. Useful for troubleshooting.
         #- -report-ingress-status
         #- -external-service=nginx-ingress
         #- -enable-leader-election
          - -enable-prometheus-metrics
         #- -enable-custom-resources
`
