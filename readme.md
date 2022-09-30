# 二狗 ergo

[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ysicing_ergo&metric=ncloc)](https://sonarcloud.io/dashboard?id=ysicing_ergo)
![GitHub Workflow Status (event)](https://img.shields.io/github/workflow/status/ysicing/ergo/tag?style=flat-square)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/ysicing/ergo?filename=go.mod&style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/ysicing/ergo?style=flat-square)
![GitHub all releases](https://img.shields.io/github/downloads/ysicing/ergo/total?style=flat-square)
![GitHub](https://img.shields.io/github/license/ysicing/ergo?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/ysicing/ergo)](https://goreportcard.com/report/github.com/ysicing/ergo)
[![Releases](https://img.shields.io/github/release-pre/ysicing/ergo.svg)](https://github.com/ysicing/ergo/releases)
[![docs](https://img.shields.io/badge/docs-done-green)](https://ysicing.github.io/ergo/)

> 一款使用 Go 编写的轻量运维工具集,尽量减少重复工作，同时降低维护脚本的成本

兼容性:

- [x] `100%` 支持 `Debian 11`
- [ ] 部分功能在macOS系上测试通过

## ergo能干什么

- 将常用脚本或者公有云操作抽象成cli命令, 简化工作
- 灵活的自定义插件管理工具,像使用`helm repo`方式管理插件
- 与大猫云平台集成

## 安装

### 二进制安装

从 [Github Release](https://github.com/ysicing/ergo/releases) 下载已经编译好的二进制文件:

### macOS安装

- 支持brew方式

```bash
brew tap ysicing/tap
brew install ergo
```

- 支持容器Docker

```bash
ysicing/ergo
```

### Debian系安装

```bash
echo "deb [trusted=yes] https://debian.ysicing.me/ /" | sudo tee /etc/apt/sources.list.d/ergo.list
apt update
# 避免与源里其他ergo混淆,deb包为opsergo
apt-get install -y opsergo
ergo version
```

### 源码编译安装

- 支持go v1.18+

```bash
# Clone the repo
# Build and run the executable
make build && ./dist/ergo_darwin_amd64 
```

### 升级

```bash
# macOS
brew upgrade
# apt / debian
apt-get update
apt-get --only-upgrade install opsergo
# default
ergo upgrade
# other
ergo ops wget https://github.com/ysicing/ergo/releases/latest/download/ergo_linux_amd64
/root/.ergo/tmp/ergo_linux_amd64 experimental install
```

## 文档

具体参见[文档](./docs/index.md)

### 中国大陆用户

> 默认github相关资源使用ghproxy代理，可使用`export NO_MIRROR=6wa6wa`不使用代理加速地址

### 特性-插件

> 默认支持`ergo-`插件, 类似krew

```bash
# 列出已安装插件
ergo addons list
repo    name       version
ysicing docker     latest
ysicing dockercfg  latest
ysicing go         1.18.1
ysicing etcd       3.5
ysicing etcdctl    3.5.3
ysicing mysql      5.7
ysicing postgresql 14

# ergo-doge插件
cat /usr/local/bin/ergo-doge
#!/bin/bash
echo $@

# 使用
ergo doge haha
haha

# 插件仓库列表
ergo repo list
[info]   上次变更时间: 2022-04-26 00:03:13.617004838 +0800 CST
name    path                                                                        source
ysicing https://github.com/ysicing/ergo-index/releases/latest/download/default.yaml remote

# 列出远程插件
ergo addons search
Repo    Name
ysicing autok3s
ysicing cilium
```

#### 其他开源项目

> 感谢以下项目

- [loft-sh/devspace](https://github.com/loft-sh/devspace)
- [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes)
- [helm/helm](https://github.com/helm/helm)
- [cilium/cilium-cli](https://github.com/cilium/cilium-cli)

## 🎉🎉 赞助商

[![jetbrains](docs/jetbrains.svg)](https://www.jetbrains.com/?from=ergo)

## 📊 Stats

![Alt](https://repobeats.axiom.co/api/embed/7067f86501e4c17c2f638dcc419df0a047b01208.svg "Repobeats analytics image")
