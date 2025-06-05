# feishu2md

[![Golang - feishu2md](https://img.shields.io/github/go-mod/go-version/wsine/feishu2md?color=%2376e1fe&logo=go)](https://go.dev/)
[![Unittest](https://github.com/Wsine/feishu2md/actions/workflows/unittest.yaml/badge.svg)](https://github.com/Wsine/feishu2md/actions/workflows/unittest.yaml)
[![Release](https://img.shields.io/github/v/release/wsine/feishu2md?color=orange&logo=github)](https://github.com/Wsine/feishu2md/releases)
[![Docker - feishu2md](https://img.shields.io/badge/Docker-feishu2md-2496ed?logo=docker&logoColor=white)](https://hub.docker.com/r/wwwsine/feishu2md)
[![Render - feishu2md](https://img.shields.io/badge/Render-feishu2md-4cfac9?logo=render&logoColor=white)](https://feishu2md.onrender.com)
![Last Review](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fbadge-last-review.wsine.workers.dev%2FWsine%2Ffeishu2md&query=%24.reviewed_at&label=last%20review)

这是一个将飞书文档转换为 Markdown 格式的工具，支持通过命令行或 Web 服务的方式使用。它可以下载单个飞书文档、批量下载文件夹中的文档或下载整个知识库的文档，并将其转换为 Markdown 格式。

**请看这里：招募有需求和有兴趣的开发者，共同探讨开发维护，有兴趣请联系。**

## 动机

[《一日一技 | 我开发的这款小工具，轻松助你将飞书文档转为 Markdown》](https://sspai.com/post/73386)

## 安装

### 二进制安装

从 [GitHub Releases](https://github.com/Wsine/feishu2md/releases) 下载适合您平台的预编译二进制文件，解压后放置在系统 PATH 路径中。

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/Wsine/feishu2md.git
cd feishu2md

# 构建命令行工具
make build

# 构建 Web 服务
make server

# Windows 下直接使用 go build 命令构建
go build -o feishu2md.exe ./cmd
```

### 使用 Docker

```bash
docker run -it --rm -p 8080:8080 -e FEISHU_APP_ID=<your id> -e FEISHU_APP_SECRET=<your secret> -e GIN_MODE=release wwwsine/feishu2md
```

## 配置

### 获取飞书 API 凭据

配置文件需要填写 APP ID 和 APP SECRET 信息，请参考 [飞书官方文档](https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/get-) 获取。推荐设置为

1. 访问[飞书开发者平台](https://open.feishu.cn/app)
2. 创建一个企业自建应用（个人版），信息随意填写
3. （重要）打开权限管理，开通以下必要的权限（可点击以下链接参考 API 调试台->权限配置字段）
   - [获取文档基本信息](https://open.feishu.cn/document/server-docs/docs/docs/docx-v1/document/get)，「查看新版文档」权限 `docx:document:readonly`
   - [获取文档所有块](https://open.feishu.cn/document/server-docs/docs/docs/docx-v1/document/list)，「查看新版文档」权限 `docx:document:readonly`
   - [下载素材](https://open.feishu.cn/document/server-docs/docs/drive-v1/media/download)，「下载云文档中的图片和附件」权限 `docs:document.media:download`
   - [获取文件夹中的文件清单](https://open.feishu.cn/document/server-docs/docs/drive-v1/folder/list)，「查看、评论、编辑和管理云空间中所有文件」权限 `drive:file:readonly`
   - [获取知识空间节点信息](https://open.feishu.cn/document/server-docs/docs/wiki-v2/space-node/get_node)，「查看知识库」权限 `wiki:wiki:readonly`
4. 打开凭证与基础信息，获取 App ID 和 App Secret

### 配置工具

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

## 使用方法

注意：飞书旧版文档的下载工具已决定不再维护，但分支 [v1_support](https://github.com/Wsine/feishu2md/tree/v1_support) 仍可使用，对应的归档为 [v1.4.0](https://github.com/Wsine/feishu2md/releases/tag/v1.4.0)，请知悉。

<details>
  <summary>命令行版本</summary>

  借助 Go 语言跨平台的特性，已编译好了主要平台的可执行文件，可以在 [Release](https://github.com/Wsine/feishu2md/releases) 中下载，并将相应平台的 feishu2md 可执行文件放置在 PATH 路径中即可。

  ### 查看帮助

  ```bash
  # 查看主命令帮助
  $ feishu2md -h
  NAME:
    feishu2md - Download feishu/larksuite document to markdown file

  USAGE:
    feishu2md [global options] command [command options] [arguments...]

  VERSION:
    v2-0e25fa5

  COMMANDS:
    config        Read config file or set field(s) if provided
    download, dl  Download feishu/larksuite document to markdown file
    help, h       Shows a list of commands or help for one command

  GLOBAL OPTIONS:
    --help, -h     show help (default: false)
    --version, -v  print the version (default: false)

  # 查看配置命令帮助
  $ feishu2md config -h
  NAME:
     feishu2md config - Read config file or set field(s) if provided

  USAGE:
     feishu2md config [command options] [arguments...]

  OPTIONS:
     --appId value      Set app id for the OPEN API
     --appSecret value  Set app secret for the OPEN API
     --help, -h         show help (default: false)

  # 查看下载命令帮助
  $ feishu2md dl -h
  NAME:
    feishu2md download - Download feishu/larksuite document to markdown file
 
  USAGE:
    feishu2md download [command options] <url>
 
  OPTIONS:
    --output value, -o value  Specify the output directory for the markdown files (default: "./")
    --dump                    Dump json response of the OPEN API (default: false)
    --batch                   Download all documents under a folder (default: false)
    --wiki                    Download all documents within the wiki. (default: false)
    --help, -h                show help (default: false)
  ```

  ### 下载单个文档

  ```bash
  # 下载单个文档到当前目录
  $ feishu2md dl "https://domain.feishu.cn/docx/docxtoken"

  # 下载单个文档到指定目录
  $ feishu2md dl -o output_directory "https://domain.feishu.cn/docx/docxtoken"

  # 下载单个文档并导出 API 响应的 JSON 数据
  $ feishu2md dl --dump "https://domain.feishu.cn/docx/docxtoken"

  # 使用命令行提供的 App ID 和 App Secret 下载文档
  $ feishu2md dl --appId <your_app_id> --appSecret <your_app_secret> "https://domain.feishu.cn/docx/docxtoken"
  ```

  文档链接可以通过 **分享 > 开启链接分享 > 互联网上获得链接的人可阅读 > 复制链接** 获得。

  ### 批量下载文件夹中的文档

  此功能暂时不支持Docker版本

  ```bash
  # 批量下载文件夹中的所有文档
  $ feishu2md dl --batch "https://domain.feishu.cn/drive/folder/foldertoken"

  # 批量下载文件夹中的所有文档到指定目录
  $ feishu2md dl --batch -o output_directory "https://domain.feishu.cn/drive/folder/foldertoken"
  ```

  文件夹链接可以通过 **分享 > 开启链接分享 > 互联网上获得链接的人可阅读 > 复制链接** 获得。

  ### 下载知识库中的文档

  ```bash
  # 下载知识库中的所有文档
  $ feishu2md dl --wiki "https://domain.feishu.cn/wiki/settings/123456789101112"

  # 下载知识库中的所有文档到指定目录
  $ feishu2md dl --wiki -o output_directory "https://domain.feishu.cn/wiki/settings/123456789101112"
  ```

  wiki settings链接可以通过打开知识库设置获得。

</details>

<details>
  <summary>Web 服务版本</summary>

  ### 使用 Docker 启动 Web 服务

  Docker 镜像：https://hub.docker.com/r/wwwsine/feishu2md

  Docker 命令：
  ```bash
  docker run -it --rm -p 8080:8080 -e FEISHU_APP_ID=<your id> -e FEISHU_APP_SECRET=<your secret> -e GIN_MODE=release wwwsine/feishu2md
  ```

  Docker Compose:

  ```yml
  # docker-compose.yml
  version: '3'
  services:
    feishu2md:
      image: wwwsine/feishu2md
      environment:
        FEISHU_APP_ID: <your id>
        FEISHU_APP_SECRET: <your secret>
        GIN_MODE: release
      ports:
        - "8080:8080"
  ```

  启动服务：
  ```bash
  docker compose up -d
  ```

  ### 手动启动 Web 服务

  ```bash
  # 使用环境变量设置 App ID 和 App Secret
  export FEISHU_APP_ID=<your_app_id>
  export FEISHU_APP_SECRET=<your_app_secret>

  # 启动 Web 服务
  ./feishu2md4web
  ```

  ### 使用 Web 界面

  1. 在浏览器中访问 `http://localhost:8080`
  2. 在输入框中粘贴飞书文档的 URL
  3. 点击下载按钮
  4. 根据文档是否包含图片，浏览器会下载 Markdown 文件或包含 Markdown 和图片的 ZIP 文件

  文档链接可以通过 **分享 > 开启链接分享 > 复制链接** 获得。
</details>

<details>
  <summary>在线版本</summary>

  我使用个人的测试 API Token 部署了一个 Unstable 版本在 Render 平台上，该版本不会保存任何的文档资料和图片在容器中，直接通过 HTTP 从**内存**中返回压缩包文件，但是 Render 平台的 Log 可能会记录一些 HTTP 信息。

  在版本仅供不在意隐私或懒于配置的用户临时使用，也可用于测试对比是否自己的 Token 权限配置有问题。Render 平台使用免费配额，仅有 512M 内存，不保证高可用性，信任链全靠开源代码，请自行斟酌。

  访问 https://feishu2md.onrender.com/ 粘贴文档链接即可，文档链接可以通过 **分享 > 开启链接分享 > 复制链接** 获得。
</details>

## 注意事项

### URL 格式

- 文档 URL 格式：`https://domain.feishu.cn/docx/docxtoken`
- 文件夹 URL 格式：`https://domain.feishu.cn/drive/folder/foldertoken`
- 知识库 URL 格式：`https://domain.feishu.cn/wiki/settings/wikitoken`

### 权限要求

确保您的飞书应用具有以下权限：

- 查看、评论和编辑文档
- 查看和下载云空间中的文件
- 访问和管理知识库

### 常见问题

#### 无法下载文档

- 检查 App ID 和 App Secret 是否正确
- 确认文档 URL 格式是否正确
- 确认您有权限访问该文档
- 检查飞书应用是否具有必要的权限

#### 图片无法显示

- 确认 `skip_img_download` 配置项为 `false`
- 检查图片目录是否存在且可写
- 确认 Markdown 查看器能够正确解析相对路径的图片

## 高级用法

### 自定义输出格式

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

### 与其他工具集成

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

## 感谢

- [chyroc/lark](https://github.com/chyroc/lark)
- [chyroc/lark_docs_md](https://github.com/chyroc/lark_docs_md)
