// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

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
