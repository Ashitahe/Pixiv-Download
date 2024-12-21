package download

import (
	"fmt"
	"net/url"
	"path"
	"time"

	"pixivDownload/pkg/collector"
	"pixivDownload/pkg/config"
	"pixivDownload/pkg/random"
	"pixivDownload/pkg/storage"

	"github.com/gocolly/colly"
)

func DownloadImage(url, savePath string) error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		c := collector.CreateCollector()
		var isError bool

		c.OnError(func(r *colly.Response, err error) {
			fmt.Printf("Download error (attempt %d/%d): %v\n", i+1, maxRetries, err)
			isError = true
		})

		c.OnResponse(func(r *colly.Response) {
			imgName := path.Base(r.Request.URL.Path)
			if err := storage.SaveFile(r.Body, savePath, imgName); err != nil {
				isError = true
			}
		})

		c.Visit(url)
		c.Wait()

		if !isError {
			return nil
		}

		// 如果不是最后一次重试，则等待一段时间后重试
		if i < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	return fmt.Errorf("failed to download image after %d attempts", maxRetries)
}

func getProxyHostDomain(originalHost string) string {
	proxyList := config.GlobalConfig.ProxyHosts
	if len(proxyList) == 0 {
		return originalHost
	}
	return random.GetRandomElement(proxyList)
}

func ProcessLink(linkStr string) string {
	parsed, _ := url.Parse(linkStr)
	originImgUrl := &url.URL{
		Scheme: "https",
		Host:   getProxyHostDomain(parsed.Host),
		Path:   parsed.Path,
	}
	return originImgUrl.String()
}
