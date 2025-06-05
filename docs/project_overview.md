# feishu2md 项目概述

## 项目简介

feishu2md 是一个用于将飞书文档转换为 Markdown 格式的工具，使用 Go 语言实现。该工具支持通过命令行或 Web 服务的方式，下载单个飞书文档、批量下载文件夹中的文档或下载整个知识库的文档，并将其转换为 Markdown 格式。

## 主要功能模块

根据代码分析，feishu2md 项目主要包含以下功能模块：

### 1. 命令行接口模块 (cmd)

- **main.go**: 定义命令行工具的入口点和主要命令结构
- **config.go**: 处理配置文件的读取、创建和更新
- **download.go**: 实现文档下载的核心功能，包括单文档下载、批量下载和知识库下载

### 2. 核心功能模块 (core)

- **client.go**: 封装与飞书 API 的交互，提供文档获取、图片下载等功能
- **config.go**: 定义配置结构和配置文件操作
- **parser.go**: 解析飞书文档内容，将其转换为 Markdown 格式

### 3. Web 服务模块 (web)

- **main.go**: Web 服务的入口点，设置路由和模板
- **download.go**: 处理 Web 界面的文档下载请求

### 4. 工具函数模块 (utils)

- **url.go**: 提供 URL 验证和处理功能
- **common.go**: 提供通用工具函数

## 执行流程

### 命令行模式执行流程

```mermaid
flowchart TD
    A[开始] --> B[解析命令行参数]
    B --> C{命令类型}
    C -->|config| D[处理配置命令]
    C -->|download| E[处理下载命令]
    
    D --> D1[获取配置文件路径]
    D1 --> D2{配置文件存在?}
    D2 -->|是| D3[读取现有配置]
    D2 -->|否| D4[创建默认配置]
    D3 --> D5{提供了新参数?}
    D5 -->|是| D6[更新配置]
    D5 -->|否| D7[显示当前配置]
    D4 --> D7
    D6 --> D7
    D7 --> Z[结束]
    
    E --> E1[获取配置文件]
    E1 --> E2{配置文件存在?}
    E2 -->|是| E3[读取现有配置]
    E2 -->|否| E4[创建默认配置]
    E3 --> E5{命令行提供了凭据?}
    E5 -->|是| E6[更新配置]
    E5 -->|否| E7[使用现有配置]
    E4 --> E7
    E6 --> E7
    
    E7 --> E8[创建客户端]
    E8 --> E9{下载类型}
    E9 -->|单文档| E10[下载单个文档]
    E9 -->|批量| E11[下载文件夹中所有文档]
    E9 -->|知识库| E12[下载知识库中所有文档]
    
    E10 --> E13[验证文档URL]
    E13 --> E14[获取文档内容]
    E14 --> E15[解析为Markdown]
    E15 --> E16[下载文档中的图片]
    E16 --> E17[保存Markdown文件]
    E17 --> Z
    
    E11 --> E18[验证文件夹URL]
    E18 --> E19[递归获取文件夹内容]
    E19 --> E20[并发下载文档]
    E20 --> Z
    
    E12 --> E21[验证知识库URL]
    E21 --> E22[获取知识库节点列表]
    E22 --> E23[递归下载知识库节点]
    E23 --> Z
```

### Web 服务模式执行流程

```mermaid
flowchart TD
    A[开始] --> B[初始化Web服务]
    B --> C[设置路由]
    C --> D[等待用户请求]
    D --> E[接收下载请求]
    E --> F[验证文档URL]
    F --> G[创建客户端]
    G --> H[获取文档内容]
    H --> I[解析为Markdown]
    I --> J[下载文档中的图片]
    J --> K{有图片?}
    K -->|是| L[创建ZIP包含Markdown和图片]
    K -->|否| M[直接返回Markdown文件]
    L --> N[返回ZIP文件]
    M --> O[结束]
    N --> O
```

## 核心数据流

```mermaid
flowchart LR
    A[用户输入] --> B[命令解析]
    B --> C[配置管理]
    C --> D[客户端创建]
    D --> E[飞书API交互]
    E --> F[文档内容获取]
    F --> G[Markdown解析]
    G --> H[图片下载]
    H --> I[文件保存]
    I --> J[输出结果]
```

## 模块依赖关系

```mermaid
flowchart TD
    A[cmd/main.go] --> B[cmd/config.go]
    A --> C[cmd/download.go]
    B --> D[core/config.go]
    C --> D
    C --> E[core/client.go]
    C --> F[core/parser.go]
    C --> G[utils/url.go]
    E --> H[chyroc/lark API]
    F --> I[文档解析逻辑]
    
    J[web/main.go] --> K[web/download.go]
    K --> E
    K --> F
    K --> G
```

## 配置文件结构

```mermaid
classDiagram
    class Config {
        +FeishuConfig Feishu
        +OutputConfig Output
    }
    
    class FeishuConfig {
        +string AppId
        +string AppSecret
    }
    
    class OutputConfig {
        +string ImageDir
        +bool TitleAsFilename
        +bool UseHTMLTags
        +bool SkipImgDownload
        +bool Delta
    }
    
    Config "1" *-- "1" FeishuConfig
    Config "1" *-- "1" OutputConfig
```