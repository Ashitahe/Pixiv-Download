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
