package collector

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"pixivDownload/pkg/config"

	"github.com/gocolly/colly"
)

func CreateCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent(getRandomUserAgent()),
	)

	// 设置 cookies
	cookies := config.ParseCookies(config.GlobalConfig.Cookie)
	for name, value := range cookies {
		c.SetCookies(".pixiv.net", []*http.Cookie{
			{
				Name:   name,
				Value:  value,
				Domain: ".pixiv.net",
			},
		})
	}

	c.Limit(&colly.LimitRule{
		Parallelism: runtime.NumCPU(),
		RandomDelay: 5 * time.Second,
	})

	tr := &http.Transport{
		TLSHandshakeTimeout: 1 * time.Minute,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
	}

	c.SetRequestTimeout(5 * time.Minute)

	c.WithTransport(tr)

	c.OnRequest(func(r *colly.Request) {
		// 基本请求头
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br, zstd")

		// 设置 referer
		r.Headers.Set("Referer", "https://www.pixiv.net/")

		// 缓存控制
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Pragma", "no-cache")
		
		// 安全相关
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		
		// Chrome 特定标识
		r.Headers.Set("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("sec-ch-ua-platform", "macOS")
		
		// Fetch 元数据
		r.Headers.Set("Priority", "u=0, i")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
	})

	// TODO: 使用代理
	// proxies := []string{
	// 	"http://proxy1.example.com:8080",
	// 	"http://proxy2.example.com:8080",
	// 	// ... 更多代理
	// }
	// rp, err := proxy.RoundRobinProxySwitcher(proxies...)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.SetProxyFunc(rp)

	return c
}

func getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2.1 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 14; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Mobile Safari/537.36",
	}

	return userAgents[rand.Intn(len(userAgents))]
}
