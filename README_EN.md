[中文](README.md) | English

# Pixiv Download

Pixiv Download is a tool for downloading images from Pixiv using various methods such as by illustration ID, user ID, or from a JSON file. It also includes a DNS over HTTPS (DoH) test feature.

## Features

- Download images by illustration ID
- Download images by user ID
- Download images from a JSON file
- DNS over HTTPS (DoH) test

## How to use

First, ensure that your network can connect to Pixiv. Then, click on the image that you wish to download. At this point, you should be able to see the illustration's ID in the browser's address bar. It will look something like this ![addressImage](https://article.biliimg.com/bfs/article/dbcb9f66dec8a99931a40df5ef8c1ff8b913104d.png). After noticing the ID, run the program and choose the option to enter the illustration's ID.You can use the same method to obtain the ID of an illustrator and download all of their artworks.

## Update v0.1.0

### New Features

1. Configuration file system

   - Automatically create `~/.pixiv-download/config.json` in the user directory
   - Support for proxy domain name list configuration
   - Support for pasting browser Cookie strings directly

2. Request header optimization

   - Add complete modern browser request headers
   - Support for the latest encoding format (including zstd)
   - Add security-related headers
   - Simulate Chrome v131 browser characteristics

3. Download performance optimization

   - Increase connection pool configuration
   - Optimize timeout settings
   - Add random User-Agent

### Configuration file example

```json
{
  "proxy_hosts": [
    "pimg.rem.asia"
  ],
  "cookie": ""
}
```

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/pixiv-download.git
   cd pixiv-download
   ```

2. Initialize the Go module:
   ```sh
   go mod init pixivDownload
   go mod tidy
   ```

## Usage

1. Run the application:

   ```sh
   go run main.go
   ```

2. Follow the menu options to download images or perform a DoH test:
   ```
   1. Download image by illust's id
   2. Download image by uid
   3. Download image by json file
   4. Doh test
   ```
