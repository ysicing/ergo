// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

const ali_kubernetes_cronhpa_controller = `
#!/bin/bash

kubectl apply -f https://gitee.com/ysicing/helminit/raw/master/yaml/ali_kubernetes_cronhpa_controller.yaml
`

const xali_kubernetes_cronhpa_controller = `
#!/bin/bash

kubectl delete -f https://gitee.com/ysicing/helminit/raw/master/yaml/ali_kubernetes_cronhpa_controller.yaml
`

const tekton = `
#!/bin/bash

kubectl apply -f https://gitee.com/ysicing/helminit/raw/master/yaml/tekton-releases-v0.17.1.yaml
`

const xtekton = `
#!/bin/bash

kubectl delete -f https://gitee.com/ysicing/helminit/raw/master/yaml/tekton-releases-v0.17.1.yaml
`

const metrics_server = `
#!/bin/bash

kubectl apply -f https://gitee.com/ysicing/helminit/raw/master/yaml/metrics-server-0.4.1.yaml
`

const xmetrics_server = `
#!/bin/bash

kubectl delete -f https://gitee.com/ysicing/helminit/raw/master/yaml/metrics-server-0.4.1.yaml
`

const kubernetes_dashboard = `
#!/bin/bash

kubectl apply -f https://gitee.com/ysicing/helminit/raw/master/yaml/k8s-dashboard.2.1.0.yaml

echo "登录节点使用kdtoken获取token"
echo "访问使用kubectl port-forward --namespace kubernetes-dashboard service/kubernetes-dashboard  10443:443"
`

const xkubernetes_dashboard = `
#!/bin/bash

kubectl delete -f https://gitee.com/ysicing/helminit/raw/master/yaml/k8s-dashboard.2.1.0.yaml
`
