// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

// https://github.com/bitnami/charts/tree/master/bitnami/nginx-ingress-controller

const nginxIngressController = `
#!/bin/bash

helminit 

kubectl create ns nginx-ingress

helm upgrade -i nginx-ingress-controller -f https://gitee.com/godu/helminit/raw/master/nginx-ingress-controller.5.6.12.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 5.6.12
`

const xnginxIngressController = `
#!/bin/bash

helm delete nginx-ingress-controller -n nginx-ingress
`

// https://github.com/bitnami/charts/tree/master/bitnami/metallb

const metallb = `
#!/bin/bash

helminit 

kubectl create ns metallb-system

helm upgrade -i metallb -f https://gitee.com/godu/helminit/raw/master/metallb.0.1.24.yaml -n metallb-system bitnami/metallb --version 0.1.24
`

const xmetallb = `
#!/bin/bash

helm delete metallb -n metallb-system
`

const etcd = `
#!/bin/bash

helminit 

kubectl create ns ops

helm upgrade -i etcd -f https://gitee.com/godu/helminit/raw/master/etcd-4.12.0.yaml -n ops bitnami/etcd --version 4.12.0
`

const xetcd = `
#!/bin/bash

helm delete etcd -n ops
`

