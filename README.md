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

- [x] version
- [x] upgrade
- [x] debian
  - [x] init
  - [x] upcore
  - [ ] 🎉lima
  - [ ] vagrant(2.x deprecated)
- [ ] `acr`/`tcr`/`gh`镜像管理

  

#### 🎉🎉 参考其他开源项目

- [sealos](https://github.com/fanux/sealos) `k8s基于sealos安装部分`
- [zzz](https://github.com/sohaha/zzz) `codegen参考zzz init部分`
- [devspace](https://github.com/loft-sh/devspace)

## 🎉🎉 赞助商

[![jetbrains](docs/jetbrains.svg)](https://www.jetbrains.com/?from=ergo)
