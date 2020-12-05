// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

// https://github.com/bitnami/charts/tree/master/bitnami/nginx-ingress-controller

const nginxIngressController = `
#!/bin/bash

helminit 

kubectl create ns nginx-ingress

helm upgrade -i nginx-ingress-controller -f {{ URL }}/master/nginx-ingress-controller/values.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 6.0.1

# helm upgrade -i nginx-ingress-controller -f https://gitee.com/ysicing/helminit/raw/master/nginx-ingress-controller/values.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 6.0.1

# helm upgrade -i nginx-ingress-controller -f https://raw.githubusercontent.com/ysicing/helminit/master/nginx-ingress-controller/values.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 6.0.1

# helm upgrade -i flagger -f {{ URL }}/master/flagger/values.yaml  flagger/flagger -n nginx-ingress --version 1.3.0
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

helm upgrade -i metallb -f {{ URL }}/master/metallb/values.yaml -n metallb-system bitnami/metallb --version 1.0.1
`

const xmetallb = `
#!/bin/bash

helm delete metallb -n metallb-system
`

const etcd = `
#!/bin/bash

helminit 

kubectl create ns ops

helm upgrade -i etcd -f {{ URL }}/master/etcd/values.yaml -n ops bitnami/etcd --version 5.3.0
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

helm upgrade -i cert-manager -n cert-manager -f {{ URL }}/master/cert-manager/values.yaml --version v1.0.3 jetstack/cert-manager
`

const xcm = `
#!/bin/bash

helm delete cert-manager -n cert-manager
`
