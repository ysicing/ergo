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