# feishu2md 用户指南

## 简介

feishu2md 是一个将飞书文档转换为 Markdown 格式的工具，支持命令行和 Web 界面两种使用方式。本指南将详细介绍如何安装和使用该工具。

## 安装

### 方法一：下载预编译的二进制文件

1. 访问 [GitHub Releases](https://github.com/Wsine/feishu2md/releases) 页面
2. 下载适合您操作系统的二进制文件
3. 将下载的文件放置在系统 PATH 路径中

### 方法二：使用 Docker

```bash
# 直接运行
docker run -it --rm -p 8080:8080 -e FEISHU_APP_ID=<your id> -e FEISHU_APP_SECRET=<your secret> -e GIN_MODE=release wwwsine/feishu2md

# 或使用 Docker Compose
# 创建 docker-compose.yml 文件：
# version: '3'
# services:
#   feishu2md:
#     image: wwwsine/feishu2md
#     environment:
#       FEISHU_APP_ID: <your id>
#       FEISHU_APP_SECRET: <your secret>
#       GIN_MODE: release
#     ports:
#       - "8080:8080"

docker compose up -d
```

### 方法三：从源码编译

```bash
# 克隆仓库
git clone https://github.com/Wsine/feishu2md.git
cd feishu2md

# 编译
go build -o feishu2md ./cmd

# 安装到系统路径
go install ./cmd
```

## 获取飞书 API Token

在使用 feishu2md 之前，您需要获取飞书的 API Token（AppID 和 AppSecret）：

1. 进入飞书[开发者后台](https://open.feishu.cn/app)
2. 创建企业自建应用（个人版），信息随意填写
3. **重要**：打开权限管理，开通以下必要的权限：
   - 「查看新版文档」权限 `docx:document:readonly`
   - 「下载云文档中的图片和附件」权限 `docs:document.media:download`
   - 「查看、评论、编辑和管理云空间中所有文件」权限 `drive:file:readonly`
   - 「查看知识库」权限 `wiki:wiki:readonly`
4. 打开凭证与基础信息，获取 App ID 和 App Secret

## 命令行使用方式

### 配置

首次使用前，需要配置 AppID 和 AppSecret：

```bash
feishu2md config --appId <your_app_id> --appSecret <your_app_secret>
```

查看当前配置：

```bash
feishu2md config
```

### 下载单个文档

```bash
feishu2md dl "https://domain.feishu.cn/docx/docxtoken"
```

指定输出目录：

```bash
feishu2md dl -o output_directory "https://domain.feishu.cn/docx/docxtoken"
```

保存 API 响应的 JSON 数据（用于调试）：

```bash
feishu2md dl --dump "https://domain.feishu.cn/docx/docxtoken"
```

### 批量下载文件夹中的文档

```bash
feishu2md dl --batch -o output_directory "https://domain.feishu.cn/drive/folder/foldertoken"
```

### 批量下载知识库中的文档

```bash
feishu2md dl --wiki -o output_directory "https://domain.feishu.cn/wiki/settings/123456789101112"
```

## Web 界面使用方式

### 启动 Web 服务

如果您使用 Docker 或在线版本，Web 服务已经启动。如果您使用二进制文件或从源码编译，需要手动启动 Web 服务：

```bash
# 确保已配置 AppID 和 AppSecret
feishu2md config --appId <your_app_id> --appSecret <your_app_secret>

# 启动 Web 服务
go run ./web
```

### 使用 Web 界面

1. 打开浏览器，访问 `http://localhost:8080`（或您配置的其他地址）
2. 在输入框中粘贴飞书文档链接
3. 点击「下载」按钮
4. 浏览器将自动下载转换后的 Markdown 文件（以 ZIP 格式压缩）

## 配置文件详解

配置文件位于用户配置目录下的 `feishu2md/config.json`，包含以下配置项：

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

### 配置项说明

- **feishu**：飞书 API 配置
  - **app_id**：飞书应用的 AppID
  - **app_secret**：飞书应用的 AppSecret

- **output**：输出配置
  - **image_dir**：图片保存的目录，相对于 Markdown 文件
  - **title_as_filename**：是否使用文档标题作为文件名（默认为 false，使用文档 token）
  - **use_html_tags**：是否使用 HTML 标签（默认为 false，使用 Markdown 语法）
  - **skip_img_download**：是否跳过图片下载（默认为 false，下载图片）

## 支持的文档元素

 feishu2md 支持转换以下飞书文档元素：

- 标题（一级到九级）
- 段落文本
- 加粗、斜体、删除线、下划线
- 有序列表和无序列表
- 代码块（支持多种语言的语法高亮）
- 表格
- 图片
- 引用
- 数学公式
- 待办事项（任务列表）
- 分割线
- 链接

## 常见问题

### 1. 无法获取文档内容

**问题**：使用工具时提示无法获取文档内容。

**解决方案**：
- 确认您已正确配置 AppID 和 AppSecret
- 确认您已开通所有必要的权限
- 确认文档链接是有效的，且您有权限访问该文档
- 确认文档链接是通过「分享 > 开启链接分享 > 复制链接」获得的

### 2. 图片无法显示

**问题**：转换后的 Markdown 文件中的图片无法显示。

**解决方案**：
- 确认您已开通「下载云文档中的图片和附件」权限
- 确认您没有设置 `skip_img_download` 为 true
- 检查图片文件是否已下载到指定的 `image_dir` 目录中

### 3. 批量下载失败

**问题**：批量下载文件夹或知识库时失败。

**解决方案**：
- 确认您已开通「查看、评论、编辑和管理云空间中所有文件」和「查看知识库」权限
- 确认文件夹或知识库链接是有效的，且您有权限访问
- 对于大型文件夹或知识库，可能需要更长的处理时间，请耐心等待

## 限制和注意事项

- 飞书旧版文档（docs）不再支持，请使用新版文档（docx）
- 批量下载功能暂不支持 Docker 版本
- 在线版本不保存任何文档资料和图片在容器中，但 Render 平台的 Log 可能会记录一些 HTTP 信息
- 转换后的 Markdown 文件可能与原文档的排版略有不同，特别是对于复杂的表格和布局

## 贡献和反馈

如果您发现任何问题或有改进建议，欢迎在 [GitHub Issues](https://github.com/Wsine/feishu2md/issues) 中提出，或提交 Pull Request。