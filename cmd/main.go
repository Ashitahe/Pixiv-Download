package main

import (
	"bufio"
	"fmt"
	"os"

	"pixivDownload/pkg/config"
	"pixivDownload/pkg/dns"
	"pixivDownload/pkg/jsonutil"
	"pixivDownload/pkg/search"
)

func menu() {
	fmt.Println("1. Download image by illust's id")
	fmt.Println("2. Download image by uid")
	fmt.Println("3. Download image by json file")
	fmt.Println("4. Doh test")
	fmt.Println("input your option")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	option := scanner.Text()

	config.LoadConfig("")

	switch option {
	case "1":
		fmt.Println("Please enter the illust's id of the image you want to download: ")
		scanner.Scan()
		illustID := scanner.Text()
		err := search.SearchByIllustId(illustID, "./downloads/")
		if err != nil {
			fmt.Printf("Error downloading images: %v\n", err)
		}
	case "2":
		fmt.Println("Please enter the uid of the artist you want to download images from: ")
		scanner.Scan()
		uid := scanner.Text()
		err := search.SearchByUid(uid)
		if err != nil {
			fmt.Printf("Error downloading images: %v\n", err)
		}
	case "3":
		fmt.Println("Please enter the filename of the JSON file to download images: ")
		scanner.Scan()
		filename := scanner.Text()
		err := jsonutil.ReadJSONFileToDownload(filename)
		if err != nil {
			fmt.Printf("Error processing JSON file: %v\n", err)
		}
	case "4":
		fmt.Println("Enter a domain to test DNS resolution: ")
		scanner.Scan()
		domain := scanner.Text()
		response, err := dns.GetDnsResponse(domain)
		if err != nil {
			fmt.Printf("Error getting DNS response: %v\n", err)
		} else {
			fmt.Printf("DNS response: %s\n", response)
		}
	default:
		fmt.Println("Invalid option. Exiting.")
	}
}

func main() {
	menu()
}
