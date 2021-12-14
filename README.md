# ergo

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

compatibility:

- [x] 100% support `Debian 10+`
- [ ] for macOS some features are available
- [ ] for CentOS some features are available

## ergo能干什么 / What does Ergo do?

- 将常用脚本或者公有云操作抽象成cli命令, 简化工作
- 灵活的自定义插件管理工具,像使用`helm repo`方式管理插件

## Install

### Binary

Downloaded from [release](https://github.com/ysicing/ergo/releases) pre-compiled binaries

### macOS Install

```bash
brew tap ysicing/tap
brew install ergo
```

### Running with Docker

```bash
ysicing/ergo
```

### Debian Install

```bash
echo "deb [trusted=yes] https://debian.ysicing.me/ /" | sudo tee /etc/apt/sources.list.d/ergo.list
apt update
# 避免与源里其他ergo混淆,deb包为opsergo
apt-get install -y opsergo
ergo version
```

### Building From Source

ergo is currently using go v1.16 or above. In order to build ergo from source you must:

```bash
# Clone the repo
# Build and run the executable
make build && ./dist/ergo_darwin_amd64 
```

### Upgrade

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

## Support

具体参见[文档](./docs/index.md)

### China Mainland users

> 默认github相关资源使用ghproxy代理，可使用`export NO_MIRROR=6wa6wa`不使用代理加速地址

### ergo plugin

> 默认支持`ergo-`插件, 类似krew

```bash
# 列出插件
ergo plugin list
[warn]   Unable to read directory "/Users/ysicing/bin" from your PATH: open /Users/ysicing/bin: no such file or directory. Skipping...
The following compatible plugins are available:
[info]   doge /usr/local/bin/ergo-doge
[info]   hello /Users/ysicing/.ergo/bin/ergo-hello

# ergo-doge插件
cat /usr/local/bin/ergo-doge                                   
#!/bin/bash
echo $@

# 使用
ergo doge haha  
haha

# 插件仓库列表
ergo repo list
[info]   上次变更时间: 2021-10-13 15:37:18.782145 +0800 CST
NAME      URL                                                           
default   https://raw.githubusercontent.com/ysicing/ergo-plugin/master/default.yaml

# 列出远程插件
ergo plugin ls-remote 
[done] √ 索引全部更新完成
[done] √ 加载完成.
repo          	name 	version  	homepage                           	desc                                            	url
default-plugin	helm 	v3.7.1   	https://helm.sh                    	The Kubernetes Package Manager                  	https://get.helm.sh/helm-v3.7.1-linux-amd64.tar.gz
```

#### Issue

- Q: docker compose命令不识别
  - A: 需要使用compose v2版本 [配置文档](https://github.com/docker/compose#linux)

#### 其他开源项目

> 感谢以下项目

- [loft-sh/devspace](https://github.com/loft-sh/devspace)
- [cdk-team/CDK](https://github.com/cdk-team/CDK)
- [kubernetes/kubectl](https://github.com/kubernetes/kubernetes)
- [helm/helm](https://github.com/helm/helm)

## 🎉🎉 Sponsors

[![jetbrains](docs/jetbrains.svg)](https://www.jetbrains.com/?from=ergo)
