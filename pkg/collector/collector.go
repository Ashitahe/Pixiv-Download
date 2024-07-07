package collector

import (
	"crypto/tls"
	"net/http"
	"runtime"
	"time"

	"github.com/gocolly/colly"
)

func CreateCollector() *colly.Collector {
	c := colly.NewCollector(colly.Async(true))

	c.Limit(&colly.LimitRule{
		Parallelism: runtime.NumCPU(),
		RandomDelay: 5 * time.Second,
	})

	tr := &http.Transport{
		TLSHandshakeTimeout: 1 * time.Minute,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	c.WithTransport(tr)

	return c
}
