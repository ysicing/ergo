## Ergo CHANGELOG

- v1.0.7
    - 调整docker安装bip地址，由`169.254.1.0/24`调整为 `172.30.42.1/16`, <del>非腾讯云网络还是推荐使用 `169.254.0.0/16`</del>
    - `ops install`新增支持prom & grafana

- v1.0.8
    - 更新helm版本
    - 取消codegen mirror default
    - vm init,upcore 支持local方式

- v1.0.9
    - ops支持安装go
  
- v1.0.11
    - 新增支持k8s 1.18.14, 移除支持k8s 1.19.2,1.19.3
  
- v1.0.12
    - 修复helm按照
      - helm 更新nginxIngressController到7.0.9
      - helm 更新mlb更新到2.0.2
    - cloud dns支持