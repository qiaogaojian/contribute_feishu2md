# feishu2md 用户指南

## 1. 简介

feishu2md 是一个将飞书文档转换为 Markdown 格式的工具，支持通过命令行或 Web 服务的方式使用。它可以下载单个飞书文档、批量下载文件夹中的文档或下载整个知识库的文档，并将其转换为 Markdown 格式。

## 2. 安装

### 2.1 二进制安装

从 [GitHub Releases](https://github.com/Wsine/feishu2md/releases) 下载适合您平台的预编译二进制文件，解压后放置在系统 PATH 路径中。

### 2.2 从源码构建

```bash
# 克隆仓库
git clone https://github.com/Wsine/feishu2md.git
cd feishu2md

# 构建命令行工具
make build

# 构建 Web 服务
make server
```

### 2.3 使用 Docker

```bash
docker run -it --rm -p 8080:8080 -e FEISHU_APP_ID=<your id> -e FEISHU_APP_SECRET=<your secret> -e GIN_MODE=release wwwsine/feishu2md
```

## 3. 配置

### 3.1 获取飞书 API 凭据

1. 访问[飞书开发者平台](https://open.feishu.cn/app)
2. 创建一个企业自建应用（信息可以任意填写）
3. 发布应用（无需等待审核通过）
4. 在应用页面中，找到「凭证与基础信息」，获取 App ID 和 App Secret

### 3.2 配置工具

#### 命令行方式

```bash
# 生成配置文件并设置 App ID 和 App Secret
feishu2md config --appId <your_app_id> --appSecret <your_app_secret>

# 查看当前配置
feishu2md config
```

#### 手动编辑配置文件

配置文件位置：
- Windows: `%AppData%/feishu2md/config.json`
- Linux: `$XDG_CONFIG_HOME/feishu2md/config.json` 或 `~/.config/feishu2md/config.json`
- Mac: `$XDG_CONFIG_HOME/feishu2md/config.json` 或 `~/.config/feishu2md/config.json`

配置文件格式：

```json
{
  "feishu": {
    "app_id": "飞书应用的App ID",
    "app_secret": "飞书应用的App Secret"
  },
  "output": {
    "image_dir": "static",       // 图片保存目录
    "title_as_filename": true,   // 使用文档标题作为文件名
    "use_html_tags": false,      // 使用HTML标签而非Markdown语法
    "skip_img_download": false,  // 跳过图片下载
    "delta": true               // 增量下载，跳过已存在的文件
  }
}
```

## 4. 使用方法

### 4.1 命令行使用

#### 查看帮助

```bash
# 查看主命令帮助
feishu2md --help

# 查看配置命令帮助
feishu2md config --help

# 查看下载命令帮助
feishu2md dl --help
```

#### 下载单个文档

```bash
# 下载单个文档到当前目录
feishu2md dl "https://domain.feishu.cn/docx/docxtoken"

# 下载单个文档到指定目录
feishu2md dl -o output_directory "https://domain.feishu.cn/docx/docxtoken"

# 下载单个文档并导出 API 响应的 JSON 数据
feishu2md dl --dump "https://domain.feishu.cn/docx/docxtoken"

# 使用命令行提供的 App ID 和 App Secret 下载文档
feishu2md dl --appId <your_app_id> --appSecret <your_app_secret> "https://domain.feishu.cn/docx/docxtoken"
```

#### 批量下载文件夹中的文档

```bash
# 批量下载文件夹中的所有文档
feishu2md dl --batch "https://domain.feishu.cn/drive/folder/foldertoken"

# 批量下载文件夹中的所有文档到指定目录
feishu2md dl --batch -o output_directory "https://domain.feishu.cn/drive/folder/foldertoken"
```

#### 下载知识库中的文档

```bash
# 下载知识库中的所有文档
feishu2md dl --wiki "https://domain.feishu.cn/wiki/settings/123456789101112"

# 下载知识库中的所有文档到指定目录
feishu2md dl --wiki -o output_directory "https://domain.feishu.cn/wiki/settings/123456789101112"
```

### 4.2 Web 服务使用

#### 启动 Web 服务

```bash
# 使用环境变量设置 App ID 和 App Secret
export FEISHU_APP_ID=<your_app_id>
export FEISHU_APP_SECRET=<your_app_secret>

# 启动 Web 服务
./feishu2md4web
```

或者使用 Docker：

```bash
docker run -it --rm -p 8080:8080 -e FEISHU_APP_ID=<your id> -e FEISHU_APP_SECRET=<your secret> -e GIN_MODE=release wwwsine/feishu2md
```

#### 使用 Web 界面

1. 在浏览器中访问 `http://localhost:8080`
2. 在输入框中粘贴飞书文档的 URL
3. 点击下载按钮
4. 根据文档是否包含图片，浏览器会下载 Markdown 文件或包含 Markdown 和图片的 ZIP 文件

## 5. 注意事项

### 5.1 URL 格式

- 文档 URL 格式：`https://domain.feishu.cn/docx/docxtoken`
- 文件夹 URL 格式：`https://domain.feishu.cn/drive/folder/foldertoken`
- 知识库 URL 格式：`https://domain.feishu.cn/wiki/settings/wikitoken`

### 5.2 权限要求

确保您的飞书应用具有以下权限：

- 查看、评论和编辑文档
- 查看和下载云空间中的文件
- 访问和管理知识库

### 5.3 常见问题

#### 无法下载文档

- 检查 App ID 和 App Secret 是否正确
- 确认文档 URL 格式是否正确
- 确认您有权限访问该文档
- 检查飞书应用是否具有必要的权限

#### 图片无法显示

- 确认 `skip_img_download` 配置项为 `false`
- 检查图片目录是否存在且可写
- 确认 Markdown 查看器能够正确解析相对路径的图片

## 6. 高级用法

### 6.1 自定义输出格式

编辑配置文件中的 `output` 部分：

```json
"output": {
  "image_dir": "custom_images",  // 自定义图片目录
  "title_as_filename": true,     // 使用文档标题作为文件名
  "use_html_tags": true,        // 使用HTML标签而非Markdown语法
  "skip_img_download": false,   // 是否跳过图片下载
  "delta": false                // 是否跳过已存在的文件
}
```

### 6.2 与其他工具集成

可以将 feishu2md 集成到自动化工作流中，例如：

```bash
# 下载文档并使用 pandoc 转换为其他格式
feishu2md dl "https://domain.feishu.cn/docx/docxtoken" -o temp/
pandoc -f markdown -t html temp/docxtoken.md -o output.html

# 批量下载知识库并推送到 Git 仓库
feishu2md dl --wiki -o docs/ "https://domain.feishu.cn/wiki/settings/wikitoken"
git add docs/
git commit -m "Update documentation"
git push
```