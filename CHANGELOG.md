## Ergo CHANGELOG

- v1.0.22
  - 优化vm init

- v1.0.21
  - 优化vm init

- v1.0.20
  - 支持重装腾讯云 云服务器 & 轻量应用服务器

- v1.0.19
  - vm 升级内核默认安装`wireguard`

- v1.0.18
  - k8s升级到1.18.20
  - go升级到1.16.x

- v1.0.17
  - 优化upcore

- v1.0.16
  - 优化arm架构初始化和升级内核

- v1.0.13
  - k8s子命令支持master调度
  - k8s默认支持安装1.18.15
  - helm版本升级
    - nginx-ingress-controller更新到7.4.1(flagger更新到1.6.2)
    - mlb(metallb)更新到2.2.0
    - etcd更新到5.6.0
    - cert-manager更新到v1.1.0

- v1.0.12
  - 修复helm按照
    - helm 更新nginxIngressController到7.0.9
    - helm 更新mlb到2.0.2
  - cloud dns支持

- v1.0.11
  - 新增支持k8s 1.18.14, 移除支持k8s 1.19.2,1.19.3

- v1.0.9
  - ops支持安装go

- v1.0.8
  - 更新helm版本
  - 取消codegen mirror default
  - vm init,upcore 支持local方式

- v1.0.7
  - 调整docker安装bip地址，由`169.254.1.0/24`调整为 `172.30.42.1/16`, <del>非腾讯云网络还是推荐使用 `169.254.0.0/16`</del>
  - `ops install`新增支持prom & grafana
