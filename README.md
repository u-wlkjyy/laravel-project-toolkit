# Laravel Project Creator

这是一个用Go语言编写的命令行工具，用于快速创建Laravel项目。

## 功能特点

- 创建指定版本的Laravel项目
- 自动使用当前目录名作为项目名（如果未指定）
- 支持查看所有可用的Laravel版本
- 支持优雅终止进程（Ctrl+C）

## 安装

### 方式一：下载预编译版本

1. 访问 [Releases](../../releases) 页面
2. 下载最新版本的 `laravel-project.exe`
3. 将可执行文件放入系统 PATH 环境变量包含的目录中

### 方式二：从源码编译

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

## 依赖要求

- Go 1.20或更高版本（仅编译时需要）
- Composer
- PHP 7.4或更高版本