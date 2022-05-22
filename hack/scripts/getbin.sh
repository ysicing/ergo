#!/bin/bash

echo "fetch helm v3.9.0"
[ ! -f "./hack/bin/helm-linux-amd64" ] && (
[ -f "/tmp/helm-linux-amd64.tar.gz" ] && rm -rf /tmp/helm-linux-amd64.tar.gz
[ -f "/tmp/helm-linux-arm64.tar.gz" ] && rm -rf /tmp/helm-linux-arm64.tar.gz
wget -q -O /tmp/helm-linux-amd64.tar.gz https://get.helm.sh/helm-v3.9.0-linux-amd64.tar.gz
wget -q -O /tmp/helm-linux-arm64.tar.gz https://get.helm.sh/helm-v3.9.0-linux-arm64.tar.gz
tar xzvfC /tmp/helm-linux-amd64.tar.gz /tmp
mv /tmp/linux-amd64/helm ./hack/bin/helm-linux-amd64
tar xzvfC /tmp/helm-linux-arm64.tar.gz /tmp
mv /tmp/linux-arm64/helm ./hack/bin/helm-linux-arm64
chmod +x ./hack/bin/helm-linux-amd64 ./hack/bin/helm-linux-arm64
rm -rf /tmp/linux-amd64 /tmp/linux-arm64
)


echo "fetch k3s v1.23.6+k3s1"
[ ! -f "./hack/bin/k3s-linux-amd64" ] && (
wget -q -O ./hack/bin/k3s-linux-amd64 https://github.com/k3s-io/k3s/releases/download/v1.23.6%2Bk3s1/k3s
wget -q -O ./hack/bin/k3s-linux-arm64 https://github.com/k3s-io/k3s/releases/download/v1.23.6%2Bk3s1/k3s-arm64
chmod +x ./hack/bin/k3s-linux-amd64 ./hack/bin/k3s-linux-arm64
)


echo "fetch cilium v0.10.5"
[ ! -f "./hack/bin/cilium-linux-amd64" ] && (
[ -f "/tmp/cilium-linux-amd64.tar.gz" ] && rm -rf /tmp/cilium-linux-amd64.tar.gz
[ -f "/tmp/cilium-linux-arm64.tar.gz" ] && rm -rf /tmp/cilium-linux-arm64.tar.gz
wget -q -O /tmp/cilium-linux-amd64.tar.gz https://github.com/cilium/cilium-cli/releases/download/v0.10.5/cilium-linux-amd64.tar.gz
tar xzvfC /tmp/cilium-linux-amd64.tar.gz /tmp
mv /tmp/cilium ./hack/bin/cilium-linux-amd64
wget -q -O /tmp/cilium-linux-arm64.tar.gz https://github.com/cilium/cilium-cli/releases/download/v0.10.5/cilium-linux-arm64.tar.gz
tar xzvfC /tmp/cilium-linux-arm64.tar.gz /tmp
mv /tmp/cilium ./hack/bin/cilium-linux-arm64
)

if [ ! -f "./hack/bin/.upxdone" ]; then
  upx -9 hack/bin/*
  touch ./hack/bin/.upxdone
fi
exit 0
