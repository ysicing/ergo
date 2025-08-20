# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在该代码仓库中工作时提供指导。

## 项目概述

Ergo（二狗）是一个用 Go 编写的轻量级 DevOps 工具集，旨在减少重复工作并降低脚本维护成本。它是一个 CLI 工具，将常用脚本和公有云操作抽象为 CLI 命令，具备插件管理能力。

## 核心架构

### 核心组件
- **主入口**：`ergo.go` 包含带启动初始化的 main 函数
- **命令结构**：基于 Cobra 的 CLI 框架，命令位于 `cmd/` 目录
- **插件系统**：支持 `ergo-` 前缀插件，类似 kubectl 插件（类 krew）
- **Factory 模式**：`internal/pkg/util/factory/` 提供客户端创建和日志记录接口
- **Kubernetes 集成**：通过 `k8s.io/kubectl` 依赖嵌入 kubectl 功能

### 目录结构
- `cmd/`：核心 CLI 命令（root、kubectl、version、upgrade）
- `internal/`：内部包，包括工具、日志、SSH 和 Kubernetes 客户端
- `pkg/`：公共包，用于配置、下载和工具
- `hack/`：构建脚本、Docker 设置和发布自动化
- `docs/`：通过 cobra 文档生成的自动生成文档

## 开发命令

### 构建
```bash
# 为所有平台构建
make build

# 本地开发构建（macOS amd64）
make local

# 本地安装
make install

# 清理构建产物
make clean
```

### 代码质量
```bash
# 格式化代码
make fmt

# 运行 linter
make lint

# 完整质量检查（版权 + fmt + lint）
make default
```

### 文档
```bash
# 生成文档
make doc
```

### 测试
```bash
# 运行测试
go test ./...

# 运行特定测试文件
go test ./pkg/downloader/downloader_test.go
```

### 发布与分发
```bash
# 构建 Docker 镜像
make docker

# 创建 GitHub release
make release

# 创建预发布版本
make pre-release

# 构建 deb 包
make deb

# 使用 goreleaser 测试
make snapshot
```

## 项目特定模式

### 版本管理
- 版本存储在 `version.txt` 中
- 通过 Makefile 中的 LDFLAGS 注入构建元数据，包括版本、构建日期和 git commit
- 通过构建标志维护 Kubernetes 客户端版本兼容性

### 插件架构
- 插件遵循 `ergo-<name>` 命名约定
- 通过 PATH 和自定义二进制目录（`~/.ergo/bin`）进行插件发现
- 插件处理器支持命令链和参数转发

### 日志系统
- `internal/pkg/util/log/` 中基于 Factory 的日志模式
- 多种 logger 实现（stdout、file、stream、union）
- 通过全局标志控制调试和静默模式

### 配置
- `~/.ergo/config.yaml` 中基于 YAML 的配置
- 默认配置创建和加载模式

## 构建配置

### Go 模块设置
- 模块：`github.com/ysicing/ergo`
- Go 版本：1.21+
- 主要依赖：Kubernetes 客户端库、Cobra CLI、SSH 工具

### 跨平台支持
- 支持：darwin/amd64、darwin/arm64、linux/amd64、linux/arm64、windows/amd64
- 使用 gox 进行交叉编译
- 插件执行中的 Windows 特定处理

## 重要说明

- 项目对某些依赖使用自定义分支（参见 go.mod 中的 `replace` 指令）
- 插件系统需要仔细的 PATH 管理来进行插件发现
- 通过特定构建标志维护 Kubernetes 版本兼容性
- 中文文档表明主要用户群体在中国，需要考虑代理支持