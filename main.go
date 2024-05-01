package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

type DnsResponse struct {
	Answer []struct {
		Data string `json:"data"`
	} `json:"Answer"`
}

func getDnsResponse(domain string) string {

	if domain == "" {
		domain = "www.pixiv.net"
	}

	urls := []string{
		"https://1.0.0.1/dns-query",
		"https://1.1.1.1/dns-query",
		"https://doh.dns.sb/dns-query",
		"https://cloudflare-dns.com/dns-query",
	}

	client := &http.Client{Timeout: 3 * time.Second}

	var dnsResp DnsResponse

	for _, url := range urls {
		req, _ := http.NewRequest("GET", url, nil)
		q := req.URL.Query()
		q.Add("name", domain)
		q.Add("type", "A")
		q.Add("do", "false")
		q.Add("cd", "false")
		req.URL.RawQuery = q.Encode()
		req.Header.Add("Accept", "application/dns-json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Unable to establish connection to", url)
			continue
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&dnsResp)
		if err != nil {
			fmt.Println("Unable to get according hostname info from", url)
			continue
		}

		if len(dnsResp.Answer) > 0 {
			fmt.Println("IP address:", dnsResp.Answer[0].Data, "from", dnsResp.Answer)
		} else {
			fmt.Println("No IP address found from", url, dnsResp)
		}
	}

	return dnsResp.Answer[0].Data
}

func getRandomElement(list []string) string {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}

func getProxyHostDomain() string {
	proxyList := []string{"pimg.rem.asia"}
	return getRandomElement(proxyList)
}

func createCollector() *colly.Collector {
	c := colly.NewCollector(colly.Async(true))

	c.Limit(&colly.LimitRule{
		Parallelism: runtime.NumCPU(),
		RandomDelay: 5 * time.Second,
	})

	// Create a Transport that skips SNI verification
	tr := &http.Transport{
		TLSHandshakeTimeout: 1 * time.Minute,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "",
		},
	}

	// Set the Collector's Transport
	c.WithTransport(tr)

	return c
}

func saveFile(saveContent []byte, savePath string, fileName string) error {
	// create directory if not exists.
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		fmt.Println("Error: can't create directory, using current directory instead:", err)
		savePath = "./"
	}

	fullPath := filepath.Join(savePath, fileName)
	output, err := os.Create(fullPath) // create a file to save file. say..pic.jpg.
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", fullPath, err)
		return err
	}
	defer output.Close()

	if _, err = io.Copy(output, bytes.NewReader(saveContent)); err != nil {
		fmt.Printf("Failed to write to file %s: %v\n", fullPath, err)
		return err
	}

	fmt.Printf("%s downloaded\n", fileName)

	return nil
}

func downloadImage(url string, savePath string) error {

	imgVisitor := createCollector()

	isError := false

	handleDownloadImgError := func(r *colly.Response, err error) {
		fmt.Println("Download illust error", err)
		isError = true
	}

	handleDownloadImgSuccess := func(resp *colly.Response) {
		imgName := path.Base(resp.Request.URL.Path)
		saveFile(resp.Body, savePath, imgName)
	}

	imgVisitor.OnError(handleDownloadImgError)

	imgVisitor.OnResponse(handleDownloadImgSuccess)

	imgVisitor.Visit(url)

	imgVisitor.Wait()

	if isError {
		return errors.New("download illust error")
	} else {
		return nil
	}
}

func processLink(linkStr string) string {
	parsed, _ := url.Parse(linkStr)
	originImgUrl := &url.URL{
		Scheme: "https",
		Host:   getProxyHostDomain(),
		Path:   parsed.Path,
	}
	return originImgUrl.String()
}

// 根据插画ID搜索并下载
func searchByIllustId(id string) error {

	c := createCollector()

	reqUrl := &url.URL{
		Scheme: "https",
		Host:   "www.pixiv.net",
		Path:     "/ajax/illust/" + id + "/pages",
	}

	saveImgPath := "./illusts_" + id + "/"

	isDownloaded := true

	handleGetIllustInfoSuccess := func(r *colly.Response) {
		body := string(r.Body)
		imgObjList := gjson.Get(body, "body.#.urls.original")

		downloadAndCheck := func(link string) {
			downlandRes := downloadImage(processLink(link), saveImgPath)
			if downlandRes != nil {
				isDownloaded = false
			}
		}

		imgObjList.ForEach(func(key, value gjson.Result) bool {
			downloadAndCheck(value.Str)
			return true
		})
	}

	handleGetIllustInfoError := func(r *colly.Response, err error) {
		fmt.Println("Get illust's detail error", err)
		isDownloaded = false
	}

	c.OnResponse(handleGetIllustInfoSuccess)

	c.OnError(handleGetIllustInfoError)

	c.Visit(reqUrl.String())

	c.Wait()

	if isDownloaded {
		return nil
	} else {
		return errors.New("download illust error")
	}
}

func getIllusts(jsonData string, jsonPathStr string) []string {
	imgList := make([]string, 0)
	illusts := gjson.Get(jsonData, jsonPathStr)
	illusts.ForEach(func(key, value gjson.Result) bool {
		imgList = append(imgList, key.Str)
		return true
	})
	return imgList
}

// 根据画师ID搜索并下载该画师的所有画品
func searchByUid(uid string) error {

	uidInfoUrl := &url.URL{
		Scheme: "https",
		Host:   "www.pixiv.net",
		Path: "/ajax/user/" + uid + "/profile/all",
	}

	c := createCollector()

	handleGetUserInfoError := func(r *colly.Response, err error) {
		fmt.Println("Get user's detail error", err)
	}

	handleGetUserInfoSuccess := func(r *colly.Response) {
		imgsRes := getIllusts(string(r.Body), "body.illusts")

		totalImgs := len(imgsRes)

		downloadedImgs := 0

		errImg := make([]string, 0)

		for _, imgId := range imgsRes {
			if searchByIllustId(imgId) != nil {
				errImg = append(errImg, imgId)
			} else {
				downloadedImgs++
			}
		}

		fmt.Printf("Total images: %d, downloaded: %d, failed: %d\n", totalImgs, downloadedImgs, len(errImg))
		fmt.Println("Failed images: ", errImg)
	}

	c.OnError(handleGetUserInfoError)

	c.OnResponse(handleGetUserInfoSuccess)

	c.Visit(uidInfoUrl.String())

	c.Wait()
	return nil
}

func menu() {
	fmt.Println("1.Dowanload image by illust's id")
	fmt.Println("2.Dowanload image by uid")
	fmt.Println("3.Doh test")
	fmt.Println("input your option")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	switch scanner.Text() {
	case "1":
		fmt.Println("Please enter the illust's id of the image you want to download: ")
		scanner.Scan()
		searchByIllustId(scanner.Text())
	case "2":
		fmt.Println("Please enter the uid of the image you want to download: ")
		scanner.Scan()
		searchByUid(scanner.Text())
	case "3":
		fmt.Println("Doh test domain: ")
		scanner.Scan()
		getDnsResponse(scanner.Text())
	default:
		fmt.Println("Goodbye!")
	}
}

func main() {
	menu()
}
