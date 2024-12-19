# Laravel Project Creator

这是一个用Go语言编写的命令行工具，用于快速创建Laravel项目。

## 功能特点

- 创建指定版本的Laravel项目
- 自动使用当前目录名作为项目名（如果未指定）
- 支持查看所有可用的Laravel版本

## 安装

1. 确保你的系统已安装Go语言环境
2. 克隆本仓库
3. 在项目目录下运行：

Windows系统下编译：
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"; go build
```

Linux系统下编译：
```bash
GOOS=linux GOARCH=amd64 go build
```

## 使用方法

### 创建新项目

```bash
# 创建最新版本的Laravel项目
laravel-project create project_name

# 创建指定版本的Laravel项目
laravel-project create project_name 9

# 在当前目录创建项目（使用当前目录名作为项目名，默认使用最新版本）
laravel-project create
```

### 查看可用版本

```bash
laravel-project list
```

### 查看帮助

```bash
laravel-project help
```