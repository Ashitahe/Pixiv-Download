package dns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DnsResponse struct {
	Answer []struct {
		Data string `json:"data"`
	}
}

func GetDnsResponse(domain string) (string, error) {
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

	for _, url := range urls {
		dnsResp, err := queryDNS(url, domain, client)
		if err != nil {
			continue
		}
		if len(dnsResp.Answer) > 0 {
			return dnsResp.Answer[0].Data, nil
		}
	}

	return "", fmt.Errorf("no valid DNS response")
}

func queryDNS(url, domain string, client *http.Client) (DnsResponse, error) {
	var dnsResp DnsResponse
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("name", domain)
	q.Add("type", "A")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/dns-json")

	resp, err := client.Do(req)
	if err != nil {
		return dnsResp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&dnsResp)
	if err != nil {
		return dnsResp, err
	}

	return dnsResp, nil
}
