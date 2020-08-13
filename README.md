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

## 命令支持

#### 初始化类

```bash
ergo init vm # 创建vm虚拟机
ergo init debian # 初始化debian
```

#### 安装类

```bash
ergo install docker/go/tools # 安装docker,go,tools
ergo install k8s
```

具体参数可以参考docs部分
