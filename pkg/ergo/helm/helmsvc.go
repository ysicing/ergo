// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

// https://github.com/bitnami/charts/tree/master/bitnami/nginx-ingress-controller

const nginxIngressController = `
#!/bin/bash

helminit 

kubectl create ns nginx-ingress

helm upgrade -i nginx-ingress-controller -f %v/master/nginx-ingress-controller/values.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 7.4.1

# helm upgrade -i nginx-ingress-controller -f https://gitee.com/ysicing/helminit/raw/master/nginx-ingress-controller/values.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 7.4.1

# helm upgrade -i nginx-ingress-controller -f https://raw.githubusercontent.com/ysicing/helminit/master/nginx-ingress-controller/values.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 7.4.1

# helm upgrade -i flagger -f %v/master/flagger/values.yaml  flagger/flagger -n nginx-ingress --version 1.6.2
`

const xnginxIngressController = `
#!/bin/bash

helm delete nginx-ingress-controller -n nginx-ingress
helm delete flagger -n nginx-ingress
`

// https://github.com/bitnami/charts/tree/master/bitnami/metallb

const metallb = `
#!/bin/bash

helminit 

kubectl create ns metallb-system

helm upgrade -i metallb -f %v/master/metallb/values.yaml -n metallb-system bitnami/metallb --version 2.2.0
`

const xmetallb = `
#!/bin/bash

helm delete metallb -n metallb-system
`

const etcd = `
#!/bin/bash

helminit 

kubectl create ns ops

helm upgrade -i etcd -f %v/master/etcd/values.yaml -n ops bitnami/etcd --version 5.6.0
`

const xetcd = `
#!/bin/bash

helm delete etcd -n ops
`

const cm = `
#!/bin/bash

helminit 

kubectl create ns cert-manager
helminit || (
	helm repo add jetstack https://charts.jetstack.io
	helm repo update
)

helm upgrade -i cert-manager -n cert-manager -f %v/master/cert-manager/values.yaml --version v1.1.0 jetstack/cert-manager
`

const xcm = `
#!/bin/bash

helm delete cert-manager -n cert-manager
`
