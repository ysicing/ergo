#!/bin/bash

localDnsIp="169.254.20.10"
upstreanDnsIp=$(kubectl get svc -n kube-system | grep kube-dns | awk '{ print $3 }')

yaml=$(cat <<-END
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: node-local-dns
  namespace: kube-system
  labels:
    kubernetes.io/cluster-service: "true"
    addonmanager.kubernetes.io/mode: Reconcile
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: node-local-dns
  namespace: kube-system
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
data:
  Corefile: |
    cluster.local:53 {
        errors
        cache {
                success 9984 30
                denial 9984 5
        }
        reload
        loop
        bind $localDnsIp
        forward . $upstreanDnsIp {
                force_tcp
        }
        prometheus :9253
        health $localDnsIp:8080
        }
    in-addr.arpa:53 {
        errors
        cache 30
        reload
        loop
        bind $localDnsIp
        forward . $upstreanDnsIp {
                force_tcp
        }
        prometheus :9253
        }
    ip6.arpa:53 {
        errors
        cache 30
        reload
        loop
        bind $localDnsIp
        forward . $upstreanDnsIp {
                force_tcp
        }
        prometheus :9253
        }
    .:53 {
        errors
        cache 30
        reload
        loop
        bind $localDnsIp
        forward . /etc/resolv.conf {
                force_tcp
        }
        prometheus :9253
        }
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-local-dns
  namespace: kube-system
  labels:
    k8s-app: node-local-dns
    kubernetes.io/cluster-service: "true"
    addonmanager.kubernetes.io/mode: Reconcile
spec:
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 10%
  selector:
    matchLabels:
      k8s-app: node-local-dns
  template:
    metadata:
       labels:
          k8s-app: node-local-dns
    spec:
      priorityClassName: system-node-critical
      serviceAccountName: node-local-dns
      hostNetwork: true
      dnsPolicy: Default  # Don't use cluster DNS.
      tolerations:
      - key: "CriticalAddonsOnly"
        operator: "Exists"
      nodeSelector:
        beta.kubernetes.io/os: linux
      containers:
      - name: node-cache
        image: registry.cn-hangzhou.aliyuncs.com/acs/k8s-dns-node-cache:coredns-1.5.0
        resources:
          limits:
            memory: 30Mi
          requests:
            cpu: 25m
            memory: 5Mi
        args: [ "-localip", "$localDnsIp", "-conf", "/etc/coredns/Corefile" ]
        securityContext:
          privileged: true
        ports:
        - containerPort: 53
          name: dns
          protocol: UDP
        - containerPort: 53
          name: dns-tcp
          protocol: TCP
        - containerPort: 9253
          name: metrics
          protocol: TCP
        livenessProbe:
          httpGet:
            host: $localDnsIp
            path: /health
            port: 8080
          initialDelaySeconds: 60
          timeoutSeconds: 5
        volumeMounts:
        - mountPath: /run/xtables.lock
          name: xtables-lock
          readOnly: false
        - name: config-volume
          mountPath: /etc/coredns
      volumes:
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate
      - name: config-volume
        configMap:
          name: node-local-dns
          items:
            - key: Corefile
              path: Corefile
END
)

echo "$yaml" > nodelocaldns-ds.yaml
kubectl create -f nodelocaldns-ds.yaml