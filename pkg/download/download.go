package download

import (
	"fmt"
	"net/url"
	"path"

	"pixivDownload/pkg/collector"
	"pixivDownload/pkg/random"
	"pixivDownload/pkg/storage"

	"github.com/gocolly/colly"
)

func DownloadImage(url, savePath string) error {
	c := collector.CreateCollector()

	var isError bool

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Download error: %v\n", err)
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

	if isError {
		return fmt.Errorf("download image error")
	}
	return nil
}

func getProxyHostDomain() string {
	proxyList := []string{"pimg.rem.asia"}
	return random.GetRandomElement(proxyList)
}

func ProcessLink(linkStr string) string {
	parsed, _ := url.Parse(linkStr)
	originImgUrl := &url.URL{
		Scheme: "https",
		Host:   getProxyHostDomain(),
		Path:   parsed.Path,
	}
	return originImgUrl.String()
}
