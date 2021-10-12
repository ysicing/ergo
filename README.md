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

- [x] completion
- [x] debian
  - [x] `init` 初始化debian
  - [x] `upcore` 升级debian内核
- [x] ops
  - [x] `ps` 进程
  - [x] `nc` nc
  - [x] `exec` 执行命令
- [x] plugin
  - [x] `list` 列出ergo插件
- [x] repo
  - [x] `list` 列出支持的软件包
  - [x] `install` 安装软件包
    - [x] containerd
    - [x] mysql,redis,etcd,mongodb,consul,minio,postgresql,rabbitmq
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
```

#### 其他开源项目

> 感谢以下项目

- [sealos](https://github.com/fanux/sealos)
- [zzz](https://github.com/sohaha/zzz)
- [devspace](https://github.com/loft-sh/devspace)
- [CDK](https://github.com/cdk-team/CDK)
- [kubectl](https://k8s.io/kubectl)

## 🎉🎉 赞助商

[![jetbrains](docs/jetbrains.svg)](https://www.jetbrains.com/?from=ergo)
