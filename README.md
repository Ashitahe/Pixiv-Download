中文 | [English](README_EN.md)

# PixivDownload

PixivDownload 是一个可以通过多种方式从 Pixiv 下载图片的工具，包括通过插画 ID、用户 ID 或从 JSON 文件下载。它还包括 DNS over HTTPS（DoH）测试功能。

## 功能特点

- 通过插画 ID 下载图片
- 通过用户 ID 下载图片
- 从 JSON 文件下载图片
- DNS over HTTPS（DoH）测试

## 如何使用

首先，确保你的网络可以连接到 Pixiv。然后，点击你想下载的图片。此时，你应该能在浏览器的地址栏看到插画的 ID，它看起来应该是这样的。![addressImage](https://article.biliimg.com/bfs/article/dbcb9f66dec8a99931a40df5ef8c1ff8b913104d.png)
注意到 ID 后，运行程序并选择输入插画 ID 的选项。你可以使用相同的方法获取插画家的 ID，并下载他们的所有作品。

## 更新说明 v0.1.0

### 新增功能

1. 配置文件系统

   - 自动在用户目录下创建 `~/.pixiv-download/config.json`
   - 支持代理域名列表配置
   - 支持直接粘贴浏览器 Cookie 字符串

2. 请求头优化

   - 添加完整的现代浏览器请求头
   - 支持最新的编码格式（包括 zstd）
   - 添加安全相关头部
   - 模拟 Chrome v131 浏览器特征

3. 下载性能优化
   - 增加连接池配置
   - 优化超时设置
   - 添加随机 User-Agent

### 配置文件示例

```json
{
  "proxy_hosts": [
    "pimg.rem.asia"
  ],
  "cookie": ""
}
```

## 安装

克隆仓库：

```sh
git clone https://github.com/yourusername/pixiv-download.git
cd pixiv-download
```

## 初始化 Go 模块：

```sh
go mod init pixivDownload
go mod tidy
```

## 使用方法

运行应用程序：

```sh
go run main.go
```

按照菜单选项下载图片或进行 DoH 测试：

```
1. 通过插画 ID 下载图片
2. 通过用户 ID 下载图片
3. 通过 JSON 文件下载图片
4. DoH 测试
```
