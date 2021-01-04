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

### Mac OS升级

```bash
brew upgrade
或者
ergo upgrade
```

## 命令支持 TODO

分类: 传统运维cli, 云原生运维cli, 云服务商cli

#### 传统运维cli

- [ ] debian系

```
# 新建debian vm
ergo vm new --mem 4096 --cpu 2 --num 2 --ip 11.11.11.0/24 # 内存，CPU，副本数, 默认IP端，建议使用默认的
# 初始化debian vm
ergo vm init --pass vagrant --ip 11.11.11.11 --ip 11.11.11.12
# 升级debian内核
ergo vm upcore --ip 11.11.11.11 --pk ~/.ssh/id_rsa

```

- [ ] ops


```
# ops install
## 法一， 通过参数方式
ergo.go ops install w --ip 11.11.11.11 --pk ~/.ssh/id_rsa
## 法二， 不传参数方式
ergo.go ops install --ip 11.11.11.11 --pk ~/.ssh/id_rsa
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select 安装包: 
    docker
    mysql
    etcd
    redis
↓ ▸ w

# ops exec 
ergo ops exec w  --ip 11.11.11.11 --pk ~/.ssh/id_rsa
```

#### 云原生运维cli

- [ ] 安装k8s 1.19.4/1.18.14

```
# 基于sealos 进行安装，只需要传master ip和worker ip以及节点password
# 初始化集群
ergo k8s init --km 11.11.11.11 --kv 1.19.4
# 添加节点
ergo.go k8s join --kw 11.11.11.12 --kv 1.19.4
```

- [ ] helm安装

```
# 列表
ergo helm list
# 初始化
ergo helm init --ip 11.11.11.11 
# 安装
ergo helm install nginx-ingress-controller --ip 11.11.11.11 --pass vagrant
# 卸载
ergo helm install slb --ip 11.11.11.11 --pass vagrant -x
```

#### 云服务商cli

- [ ] 阿里云镜像仓库, ucloud镜像仓库

- [ ] 阿里云dns解析

```bazaar
23:47 ➜  ergo cloud dns show godu.dev hk2
Using config file: /Users/ysicing/.config/ergo/config.yaml
A *.hk2.godu.dev ---> 127.0.0.1 *
A hk2.vps.godu.dev ---> 127.0.0.1 *

23:47 ➜  ergo cloud dns renew --domain hk2.vps.godu.dev --value 127.0.0.1
已存在记录
更新成功
```

#### 🎉🎉 参考其他开源项目

- [sealos](https://github.com/fanux/sealos) `k8s基于sealos安装部分`
- [zzz](https://github.com/sohaha/zzz) `codegen参考zzz init部分`


## 🎉🎉 赞助商

[![jetbrains](docs/jetbrains.svg)](https://www.jetbrains.com/?from=ergo)
