# feishu2md 项目架构文档

## 项目概述

feishu2md 是一个将飞书文档转换为 Markdown 格式的工具，使用 Go 语言实现。该工具支持命令行和 Web 界面两种使用方式，可以下载单个文档、批量下载文件夹中的文档或下载整个知识库的文档。

## 主要功能模块

### 1. 配置管理模块

- **功能**：管理应用的配置信息，包括飞书 API 的 AppID 和 AppSecret，以及输出相关的配置
- **核心文件**：`core/config.go`
- **主要结构**：
  - `Config`：包含飞书配置和输出配置
  - `FeishuConfig`：存储 AppID 和 AppSecret
  - `OutputConfig`：控制输出行为，如图片目录、是否使用文档标题作为文件名等

### 2. 客户端模块

- **功能**：与飞书 API 交互，获取文档内容、下载图片等
- **核心文件**：`core/client.go`
- **主要方法**：
  - `GetDocxContent`：获取文档内容
  - `DownloadImage`：下载文档中的图片
  - `GetDriveFolderFileList`：获取文件夹中的文件列表
  - `GetWikiNodeList`：获取知识库节点列表

### 3. 解析器模块

- **功能**：将飞书文档的内容解析为 Markdown 格式
- **核心文件**：`core/parser.go`
- **主要方法**：
  - `ParseDocxContent`：解析文档内容
  - `ParseDocxBlock`：解析文档块
  - 各种特定块类型的解析方法：如标题、列表、代码块、表格等

### 4. 命令行接口模块

- **功能**：提供命令行界面，处理用户输入的命令
- **核心文件**：`cmd/main.go`、`cmd/config.go`、`cmd/download.go`
- **主要命令**：
  - `config`：配置 AppID 和 AppSecret
  - `download`/`dl`：下载文档为 Markdown

### 5. Web 界面模块

- **功能**：提供 Web 界面，方便用户在浏览器中使用
- **核心文件**：`web/main.go`、`web/download.go`
- **主要功能**：通过 Web 界面上传文档链接，下载为 Markdown

## 执行流程

### 命令行模式执行流程

```mermaid
flowchart TD
    A[开始] --> B[解析命令行参数]
    B --> C{命令类型}
    C -->|config| D[读取/创建配置文件]
    C -->|download| E[处理下载命令]
    E --> F{下载类型}
    F -->|单文档| G[下载单个文档]
    F -->|批量文件夹| H[下载文件夹中的所有文档]
    F -->|知识库| I[下载知识库中的所有文档]
    G --> J[获取文档内容]
    H --> J
    I --> J
    J --> K[解析文档内容为Markdown]
    K --> L[下载文档中的图片]
    L --> M[保存Markdown文件]
    M --> N[结束]
    D --> N
```

### 单文档下载流程

```mermaid
flowchart TD
    A[开始下载文档] --> B[验证文档URL]
    B --> C[获取文档Token]
    C --> D{文档类型}
    D -->|wiki| E[获取Wiki节点信息]
    E --> F[更新文档类型和Token]
    D -->|docx| G[获取文档内容]
    F --> G
    G --> H[解析文档内容]
    H --> I{是否跳过图片下载}
    I -->|否| J[下载文档中的图片]
    I -->|是| K[格式化Markdown]
    J --> K
    K --> L[保存到文件]
    L --> M[结束]
```

### 批量下载流程

```mermaid
flowchart TD
    A[开始批量下载] --> B{下载类型}
    B -->|文件夹| C[验证文件夹URL]
    B -->|知识库| D[验证知识库URL]
    C --> E[获取文件夹Token]
    D --> F[获取知识库ID]
    E --> G[递归处理文件夹]
    F --> H[递归处理知识库节点]
    G --> I[获取文件夹文件列表]
    H --> J[获取知识库节点列表]
    I --> K{文件类型}
    J --> L{节点类型}
    K -->|文件夹| G
    K -->|文档| M[并发下载文档]
    L -->|有子节点| H
    L -->|文档| M
    M --> N[结束]
```

### 文档解析流程

```mermaid
flowchart TD
    A[开始解析] --> B[构建块映射]
    B --> C[获取入口块]
    C --> D[递归解析块]
    D --> E{块类型}
    E -->|页面| F[解析页面块]
    E -->|文本| G[解析文本块]
    E -->|标题| H[解析标题块]
    E -->|列表| I[解析列表块]
    E -->|代码| J[解析代码块]
    E -->|表格| K[解析表格块]
    E -->|图片| L[解析图片块]
    E -->|其他| M[解析其他块]
    F --> N[解析子块]
    G --> O[解析文本元素]
    H --> O
    I --> O
    J --> O
    K --> O
    L --> P[收集图片Token]
    M --> O
    N --> Q[合并结果]
    O --> Q
    P --> Q
    Q --> R[返回Markdown文本]
    R --> S[结束]
```

## 配置文件结构

```json
{
  "feishu": {
    "app_id": "your_app_id",
    "app_secret": "your_app_secret"
  },
  "output": {
    "image_dir": "static",
    "title_as_filename": false,
    "use_html_tags": false,
    "skip_img_download": false
  }
}
```

## 项目依赖

- [chyroc/lark](https://github.com/chyroc/lark)：飞书 API 客户端
- [urfave/cli](https://github.com/urfave/cli)：命令行界面框架
- [88250/lute](https://github.com/88250/lute)：Markdown 处理器
- [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter)：表格渲染