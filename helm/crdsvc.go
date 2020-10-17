// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

const ali_kubernetes_cronhpa_controller = `
#!/bin/bash

kubectl apply -f https://gitee.com/godu/helminit/raw/master/ali_kubernetes_cronhpa_controller.yaml
`

const xali_kubernetes_cronhpa_controller = `
#!/bin/bash

kubectl delete -f https://gitee.com/godu/helminit/raw/master/ali_kubernetes_cronhpa_controller.yaml
`

const tekton = `
#!/bin/bash

kubectl apply -f https://gitee.com/godu/helminit/raw/master/tekton-releases-v0.17.1.yaml
`

const xtekton = `
#!/bin/bash

kubectl delete -f https://gitee.com/godu/helminit/raw/master/tekton-releases-v0.17.1.yaml
`

const metrics_server = `
#!/bin/bash

kubectl apply -f https://gitee.com/godu/helminit/raw/master/metrics-server-0.3.7.yaml
`

const xmetrics_server = `
#!/bin/bash

kubectl delete -f https://gitee.com/godu/helminit/raw/master/metrics-server-0.3.7.yaml
`