#!/usr/bin/env bash

version=$(apt-cache madison kubectl | grep ${1:-1.18} | head -1 | awk '{print $3}')

echo "更新kubeadm $version"

apt remove kubeadm -y
apt install -y  kubeadm=${version}

kubeadm config images pull
kubeadm config images list | awk -F/ '{print $2}' | xargs -I {} docker tag k8s.gcr.io/{} ysicing/{}
kubeadm config images list | awk -F/ '{print $2}' | xargs -I {} docker push ysicing/{}

kubeadm config images list | awk -F/ '{print $2}' | xargs -I {} docker tag k8s.gcr.io/{} registry.cn-beijing.aliyuncs.com/k7scn/{}
kubeadm config images list | awk -F/ '{print $2}' | xargs -I {} docker push registry.cn-beijing.aliyuncs.com/k7scn/{}