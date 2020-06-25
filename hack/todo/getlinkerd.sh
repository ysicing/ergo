#!/usr/bin/env bash

# curl -sL https://run.linkerd.io/install | sh

linkerd install --ignore-cluster | grep image: | sed -e 's/^ *//' | sort | uniq | awk '{print $2}' | xargs -I {} docker pull {}
linkerd install --ignore-cluster | grep image: | sed -e 's/^ *//' | sort | uniq | awk '{print $2}' | grep "gcr.io/linkerd-io" | awk -F/ '{print $NF}' | xargs -I {} docker tag gcr.io/linkerd-io/{} ysicing/linkerd-io-{}
linkerd install --ignore-cluster | grep image: | sed -e 's/^ *//' | sort | uniq | awk '{print $2}' | grep "gcr.io/linkerd-io" | awk -F/ '{print $NF}' | xargs -I {} docker push ysicing/linkerd-io-{}

linkerd install --ignore-cluster | grep image: | sed -e 's/^ *//' | sort | uniq | awk '{print $2}' | grep "gcr.io/linkerd-io" | awk -F/ '{print $NF}' | xargs -I {} docker tag gcr.io/linkerd-io/{} registry.cn-beijing.aliyuncs.com/k7scn/linkerd-io-{}
linkerd install --ignore-cluster | grep image: | sed -e 's/^ *//' | sort | uniq | awk '{print $2}' | grep "gcr.io/linkerd-io" | awk -F/ '{print $NF}' | xargs -I {} docker push registry.cn-beijing.aliyuncs.com/k7scn/linkerd-io-{}