package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getBodyData(address string) string {
	response, _ := http.Get(address)
	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	return string(bytes)
}

func main() {
	siteData := getBodyData("https://www.washingtonpost.com/news-sitemap-index.xml")
	fmt.Println(siteData)
}
