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

兼容性:

- [x] 100%兼容`Debian 11+`系
- [ ] macOS部分功能可用

## ergo能干什么 / What does Ergo do?

- 将常用脚本或者公有云操作抽象成cli命令, 简化工作
- 灵活的自定义插件管理工具,像使用`helm repo`方式管理插件

## 安装使用

### 二进制安装

可直接从 [release](https://github.com/ysicing/ergo/releases) 下载预编译的二进制文件

### macOS安装

```bash
brew tap ysicing/tap
brew install ergo
```

### macOS升级

```bash
brew upgrade
或者
ergo upgrade
```

### 镜像使用

```bash
ysicing/ergo
```

### Debian使用

```bash
echo "deb [trusted=yes] https://debian.ysicing.me/ /" | sudo tee /etc/apt/sources.list.d/ergo.list
apt update
# 避免与源里其他ergo混淆,deb包为opsergo
apt-get install -y opsergo
ergo version
```

## 命令支持

具体参见[文档](./docs/index.md)

### 下载说明

> 默认github相关资源使用ghproxy代理，可使用`export NO_MIRROR=6wa6wa`不使用代理加速地址

### ergo插件

> 默认支持`ergo-`插件

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
[done] √ sync done.
Repo    NAME            URL                                                                                                             Desc                                                    Available
default tgsend-linux    https://github.techoc.workers.dev/https://github.com/mritd/tgsend/releases/download/v1.0.1/tgsend_linux_amd64   一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等 false    
default tgsend-darwin   https://github.techoc.workers.dev/https://github.com/mritd/tgsend/releases/download/v1.0.1/tgsend_darwin_amd64  一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等 true   
```

#### 已知问题

- Q: docker compose命令不识别
  - A: 需要使用compose v2版本 [配置文档](https://github.com/docker/compose#linux)

#### 其他开源项目

> 感谢以下项目

- [loft-sh/devspace](https://github.com/loft-sh/devspace)
- [cdk-team/CDK](https://github.com/cdk-team/CDK)
- [kubernetes/kubectl](https://github.com/kubernetes/kubernetes)
- [helm/helm](https://github.com/helm/helm)

## 🎉🎉 赞助商

[![jetbrains](docs/jetbrains.svg)](https://www.jetbrains.com/?from=ergo)
