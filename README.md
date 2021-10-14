[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ysicing_ergo&metric=ncloc)](https://sonarcloud.io/dashboard?id=ysicing_ergo)

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

## 命令支持

- [x] code 初始化项目
- [x] completion
- [x] debian
  - [x] `init` 初始化debian
  - [x] `upcore` 升级debian内核
- [x] ops
  - [x] `ps` 进程
  - [x] `nc` nc
  - [x] `exec` 执行命令
- [x] plugin
  - [x] `install` 安装插件
  - [x] `list` 列出ergo插件
  - [x] `ls-remote` 列出远程插件
  - [x] `repo` 插件仓库管理, 类似helm仓库
     - [x] `add` 添加插件仓库
     - [x] `list` 列出插件仓库列表
     - [x] `del` 移除插件仓库
     - [x] `update` 更新插件索引
- [x] repo
  - [x] `list` 列出支持的软件包
  - [x] `install` 安装软件包
    - [x] `containerd`
    - [x] `mysql`等
  - [x] `dump` dump安装脚本 
- [x] upgrade
- [x] version

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
ergo plugin repo list
[info]   上次变更时间: 2021-10-13 15:37:18.782145 +0800 CST
NAME      URL                                                           
default   https://raw.githubusercontent.com/ysicing/ergo-plugin/master/default.yaml

# 列出远程插件
ergo plugin ls-remote 
[done] √ "local"已经更新索引: /Users/ysicing/.ergo/.config/default.index.yaml
[done] √ sync done.
Repo    NAME            URL                                                                                                             Desc                                                    Available
default tgsend-linux    https://github.techoc.workers.dev/https://github.com/mritd/tgsend/releases/download/v1.0.1/tgsend_linux_amd64   一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等 false    
default tgsend-darwin   https://github.techoc.workers.dev/https://github.com/mritd/tgsend/releases/download/v1.0.1/tgsend_darwin_amd64  一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等 true   
```

#### 其他开源项目

> 感谢以下项目

- [loft-sh/devspace](https://github.com/loft-sh/devspace)
- [cdk-team/CDK](https://github.com/cdk-team/CDK)
- [kubernetes/kubectl](https://github.com/kubernetes/kubernetes)
- [helm/helm](https://github.com/helm/helm)

## 🎉🎉 赞助商

[![jetbrains](docs/jetbrains.svg)](https://www.jetbrains.com/?from=ergo)
