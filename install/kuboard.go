// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package install

import (
	"fmt"
	"github.com/wonderivan/logger"
)

func KuboardInstall() {
	i := InstallConfig{
		Master0: Hosts[0],
	}
	i.KuboardInstall()
	i.MetricsServerDeploy()
	i.KuboardDone()
}

func (i *InstallConfig) KuboardInstall() {
	kbcmd := fmt.Sprintf(`echo '%s' | kubectl apply -f -`, kuboardtpl)
	SSHConfig.Cmd(i.Master0, kbcmd)
}

func (i *InstallConfig) MetricsServerDeploy() {
	mscmd := fmt.Sprintf(`echo '%s' | kubectl apply -f -`, mstpl)
	SSHConfig.Cmd(i.Master0, mscmd)
}

func (i *InstallConfig) KuboardDone() {
	logger.Info("kuboard install success. visit: http://%s:32567", i.Master0)
	logger.Debug("token: ")
	SSHConfig.Cmd(i.Master0, "/usr/local/bin/kbtoken")
}

const kuboardtpl = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kuboard
  namespace: kube-system
  annotations:
    k8s.kuboard.cn/displayName: kuboard
    k8s.kuboard.cn/ingress: "true"
    k8s.kuboard.cn/service: NodePort
    k8s.kuboard.cn/workload: kuboard
  labels:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: kuboard
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s.kuboard.cn/layer: monitor
      k8s.kuboard.cn/name: kuboard
  template:
    metadata:
      labels:
        k8s.kuboard.cn/layer: monitor
        k8s.kuboard.cn/name: kuboard
    spec:
      containers:
      - name: kuboard
        image: eipwork/kuboard:beta
        imagePullPolicy: Always
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule

---
apiVersion: v1
kind: Service
metadata:
  name: kuboard
  namespace: kube-system
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    targetPort: 80
    nodePort: 32567
  selector:
    k8s.kuboard.cn/layer: monitor
    k8s.kuboard.cn/name: kuboard

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kuboard-user
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kuboard-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: kuboard-user
  namespace: kube-system

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kuboard-viewer
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kuboard-viewer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
- kind: ServiceAccount
  name: kuboard-viewer
  namespace: kube-system

---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: kuboard
  namespace: kube-system
  annotations:
    k8s.kuboard.cn/displayName: kuboard
    k8s.kuboard.cn/workload: kuboard
    nginx.org/websocket-services: "kuboard"
    nginx.com/sticky-cookie-services: "serviceName=kuboard srv_id expires=1h path=/"
spec:
  rules:
  - host: kuboard.yourdomain.com
    http:
      paths:
      - path: /
        backend:
          serviceName: kuboard
          servicePort: http
`

const mstpl = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:aggregated-metrics-reader
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
rules:
- apiGroups: ["metrics.k8s.io"]
  resources: ["pods", "nodes"]
  verbs: ["get", "list", "watch"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: metrics-server:system:auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: metrics-server-auth-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:metrics-server
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - nodes
  - nodes/stats
  - namespaces
  verbs:
  - get
  - list
  - watch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:metrics-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:metrics-server
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system

---
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: v1beta1.metrics.k8s.io
spec:
  service:
    name: metrics-server
    namespace: kube-system
  group: metrics.k8s.io
  version: v1beta1
  insecureSkipTLSVerify: true
  groupPriorityMinimum: 100
  versionPriority: 100

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: metrics-server
  namespace: kube-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: metrics-server
  namespace: kube-system
  labels:
    k8s-app: metrics-server
spec:
  selector:
    matchLabels:
      k8s-app: metrics-server
  template:
    metadata:
      name: metrics-server
      labels:
        k8s-app: metrics-server
    spec:
      serviceAccountName: metrics-server
      volumes:
      # mount in tmp so we can safely use from-scratch images and/or read-only containers
      - name: tmp-dir
        emptyDir: {}
      hostNetwork: true
      containers:
      - name: metrics-server
        image: registry.cn-hangzhou.aliyuncs.com/google_containers/metrics-server-amd64:v0.3.6
        # command:
        # - /metrics-server
        # - --kubelet-insecure-tls
        # - --kubelet-preferred-address-types=InternalIP
        args:
          - --cert-dir=/tmp
          - --secure-port=4443
          - --kubelet-insecure-tls=true
          - --kubelet-preferred-address-types=InternalIP,Hostname,InternalDNS,externalDNS
        ports:
        - name: main-port
          containerPort: 4443
          protocol: TCP
        securityContext:
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 1000
        imagePullPolicy: Always
        volumeMounts:
        - name: tmp-dir
          mountPath: /tmp
      nodeSelector:
        beta.kubernetes.io/os: linux

---
apiVersion: v1
kind: Service
metadata:
  name: metrics-server
  namespace: kube-system
  labels:
    kubernetes.io/name: "Metrics-server"
    kubernetes.io/cluster-service: "true"
spec:
  selector:
    k8s-app: metrics-server
  ports:
  - port: 443
    protocol: TCP
    targetPort: main-port
`
