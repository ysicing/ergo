#!/bin/bash

echo "fetch helm v3.9.0"
[ ! -f "./manifests/bin/helm-linux-amd64" ] && (
[ -f "/tmp/helm-linux-amd64.tar.gz" ] && rm -rf /tmp/helm-linux-amd64.tar.gz
[ -f "/tmp/helm-linux-arm64.tar.gz" ] && rm -rf /tmp/helm-linux-arm64.tar.gz
wget -q -O /tmp/helm-linux-amd64.tar.gz https://get.helm.sh/helm-v3.9.0-linux-amd64.tar.gz
wget -q -O /tmp/helm-linux-arm64.tar.gz https://get.helm.sh/helm-v3.9.0-linux-arm64.tar.gz
tar xzvfC /tmp/helm-linux-amd64.tar.gz /tmp
mv /tmp/linux-amd64/helm ./manifests/bin/helm-linux-amd64
tar xzvfC /tmp/helm-linux-arm64.tar.gz /tmp
mv /tmp/linux-arm64/helm ./manifests/bin/helm-linux-arm64
chmod +x ./manifests/bin/helm-linux-amd64 ./manifests/bin/helm-linux-arm64
rm -rf /tmp/linux-amd64 /tmp/linux-arm64
)


echo "fetch k3s v1.23.7+k3s1"
[ ! -f "./manifests/bin/k3s-linux-amd64" ] && (
wget -q -O ./manifests/bin/k3s-linux-amd64 https://github.com/k3s-io/k3s/releases/download/v1.23.7%2Bk3s1/k3s
wget -q -O ./manifests/bin/k3s-linux-arm64 https://github.com/k3s-io/k3s/releases/download/v1.23.7%2Bk3s1/k3s-arm64
chmod +x ./manifests/bin/k3s-linux-amd64 ./manifests/bin/k3s-linux-arm64
)


echo "fetch cilium v0.11.10"
[ ! -f "./manifests/bin/cilium-linux-amd64" ] && (
[ -f "/tmp/cilium-linux-amd64.tar.gz" ] && rm -rf /tmp/cilium-linux-amd64.tar.gz
[ -f "/tmp/cilium-linux-arm64.tar.gz" ] && rm -rf /tmp/cilium-linux-arm64.tar.gz
wget -q -O /tmp/cilium-linux-amd64.tar.gz https://github.com/cilium/cilium-cli/releases/download/v0.11.10/cilium-linux-amd64.tar.gz
tar xzvfC /tmp/cilium-linux-amd64.tar.gz /tmp
mv /tmp/cilium ./manifests/bin/cilium-linux-amd64
wget -q -O /tmp/cilium-linux-arm64.tar.gz https://github.com/cilium/cilium-cli/releases/download/v0.11.10/cilium-linux-arm64.tar.gz
tar xzvfC /tmp/cilium-linux-arm64.tar.gz /tmp
mv /tmp/cilium ./manifests/bin/cilium-linux-arm64
)

if [ ! -f "./manifests/bin/.upxdone" ]; then
  upx -9 manifests/bin/*
  touch ./manifests/bin/.upxdone
fi
exit 0
