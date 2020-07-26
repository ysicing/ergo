#!/usr/bin/env bash

istioctl manifest generate --set profile=demo | grep "image:" | grep -v { | awk '{print $2}' | xargs -I {} docker pull {}

# grafana version
gv=$(istioctl manifest generate --set profile=demo | grep "image:" | grep grafana | awk -F: '{print $3}')
docker pull grafana/grafana:${gv}
docker tag grafana/grafana:${gv} registry.cn-beijing.aliyuncs.com/k7scn/grafana:${gv}
docker push registry.cn-beijing.aliyuncs.com/k7scn/grafana:${gv}

# kiali version
kv=$(istioctl manifest generate --set profile=demo | grep "image:" | grep kiali | awk -F: '{print $3}')
docker pull quay.io/kiali/kiali:${kv}
docker tag quay.io/kiali/kiali registry.cn-beijing.aliyuncs.com/k7scn/kiali:${kv}
docker push registry.cn-beijing.aliyuncs.com/k7scn/kiali:${kv}

# prometheus version
pv=$(istioctl manifest generate --set profile=demo | grep "image:" | grep prom | awk -F: '{print $3}')
docker pull docker.io/prom/prometheus:${pv}
docker tag docker.io/prom/prometheus:${pv} registry.cn-beijing.aliyuncs.com/k7scn/prometheus:${pv}
docker push registry.cn-beijing.aliyuncs.com/k7scn/prometheus:${pv}

# jaegertracing version
jv=$(istioctl manifest generate --set profile=demo | grep "image:" | grep jaeger | awk -F: '{print $3}'  | awk -F\" '{print $1}')
docker pull docker.io/jaegertracing/all-in-one:${jv}
docker tag docker.io/jaegertracing/all-in-one:${jv} registry.cn-beijing.aliyuncs.com/k7scn/jaegertracing-all-in-one:${jv}
docker push registry.cn-beijing.aliyuncs.com/k7scn/jaegertracing-all-in-one:${jv}

# istio version
iv=$(istioctl version | grep -v "pods")

istioctl manifest generate --set profile=demo | grep "image:" |grep -v { | grep istio | awk '{print $2}' | xargs -I {} docker pull {}
istioctl manifest generate --set profile=demo | grep "image:" |grep -v { | grep istio | awk '{print $2}' | awk -F/ '{print $NF}' | xargs -I {} docker tag docker.io/istio/{} registry.cn-beijing.aliyuncs.com/k7scn/istio-{}
istioctl manifest generate --set profile=demo | grep "image:" |grep -v { | grep istio | awk '{print $2}' | awk -F/ '{print $NF}' | xargs -I {} docker push registry.cn-beijing.aliyuncs.com/k7scn/istio-{}