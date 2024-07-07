package search

import (
	"fmt"

	"pixivDownload/pkg/collector"
	"pixivDownload/pkg/download"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

func SearchByIllustId(id string, savePath string) error {
	c := collector.CreateCollector()
	reqUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/pages", id)

	var isDownloaded = true

	c.OnResponse(func(r *colly.Response) {
		imgObjList := gjson.Get(string(r.Body), "body.#.urls.original")
		imgObjList.ForEach(func(key, value gjson.Result) bool {
			link := download.ProcessLink(value.Str)
			if err := download.DownloadImage(link, savePath); err != nil {
				isDownloaded = false
			}
			return true
		})
	})

	c.Visit(reqUrl)
	c.Wait()

	if !isDownloaded {
		return fmt.Errorf("failed to download all images")
	}
	return nil
}

// 根据画师ID搜索并下载该画师的所有画品
func SearchByUid(uid string) error {
	c := collector.CreateCollector()
	userInfoUrl := fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/profile/all", uid)
	var isDownloaded = true

	c.OnResponse(func(r *colly.Response) {
		imgObjList := gjson.Get(string(r.Body), "body.illusts")
		imgObjList.ForEach(func(key, value gjson.Result) bool {
			imgId := key.Str
			link := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/pages", imgId)
			savePath := fmt.Sprintf("./downloads/%s/", imgId)
			err := DownloadAllImages(link, savePath)
			if err != nil {
				fmt.Printf("Failed to download images for illust ID %s: %v\n", imgId, err)
				isDownloaded = false
				return false // Stop iterating
			}
			return true // Continue iterating
		})
	})

	c.Visit(userInfoUrl)
	c.Wait()

	if !isDownloaded {
		return fmt.Errorf("failed to download all images")
	}
	return nil
}

func DownloadAllImages(illustURL string, savePath string) error {
	c := collector.CreateCollector()
	var isError bool

	c.OnResponse(func(r *colly.Response) {
		imgObjList := gjson.Get(string(r.Body), "body.#.urls.original")
		imgObjList.ForEach(func(key, value gjson.Result) bool {
			downloadLink := download.ProcessLink(value.Str)
			if err := download.DownloadImage(downloadLink, savePath); err != nil {
				isError = true
				return false // Stop iterating
			}
			return true // Continue iterating
		})
	})

	c.Visit(illustURL)
	c.Wait()

	if isError {
		return fmt.Errorf("error occurred while downloading images")
	}
	return nil
}
