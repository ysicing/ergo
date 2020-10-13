## ergo

> 一个使用 Go 编写运维工具,尽量减少重复工作，同时降低维护脚本的成本

### 镜像使用

```bash
ysicing/ergo
```

### 二进制安装

可直接从 [release](https://github.com/ysicing/ergo/releases) 页下载预编译的二进制文件

### Mac OS安装

```bash
brew tap ysicing/tap
brew install ergo
```

## 命令支持 TODO

分类: 传统运维cli, 云原生运维cli, 云服务商cli

#### 传统运维cli

- [ ] debian系

```
# 新建debian vm
ergo vm new --mem 4096 --cpu 2 --num 2 --ip 10.0.0.0/24 # 内存，CPU，副本数, 默认IP端，建议使用默认的
# 初始化debian vm
ergo vm init --pass vagrant --ips 10.0.0.11 --ips 10.0.0.12
# 安装常用工具
ergo vm install --pass vagrant --ips 10.0.0.11 --ips 10.0.0.12 docker
# 执行shell
ergo vm exec --pass vagrant --ips 10.0.0.11 --ips 10.0.0.12 docker ps
```

#### 云原生运维cli

- [ ] 安装k8s 1.19.2

```
# 基于sealos 进行安装，只需要传master ip和worker ip以及节点password
# 初始化集群
ergo k8s --km 11.11.11.11
# 添加节点
ergo.go k8s --kw 11.11.11.12 --init=false
```

- [ ] helm安装

```
# 初始化
ergo helm init --ip 11.11.11.11 
# 安装
ergo helm install nginx-ingress-controller --ip 11.11.11.11 --pass vagrant
# 卸载
ergo helm install slb --ip 11.11.11.11 --pass vagrant -x
```

#### 云服务商cli

- [ ] 阿里云镜像仓库, ucloud镜像仓库