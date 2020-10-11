// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

// https://github.com/bitnami/charts/tree/master/bitnami/nginx-ingress-controller

const nginxIngressController = `
#!/bin/bash

helminit 

kubectl create ns nginx-ingress

helm upgrade -i nginx-ingress-controller -f https://gitee.com/godu/helminit/raw/master/nginx-ingress-controller.5.6.11.yaml -n nginx-ingress bitnami/nginx-ingress-controller --version 5.6.11
`
