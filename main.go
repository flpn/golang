package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SiteMapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword  string
	Location string
}

func main() {
	var siteMapIndex SiteMapIndex
	var news News
	newsMap := make(map[string]NewsMap)

	response, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(response.Body)
	// response.Body.Close()

	xml.Unmarshal(bytes, &siteMapIndex)

	for _, Location := range siteMapIndex.Locations {
		response, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(response.Body)
		// response.Body.Close()

		xml.Unmarshal(bytes, &news)

		for index, _ := range news.Keywords {
			newsMap[news.Titles[index]] = NewsMap{news.Keywords[index], news.Locations[index]}
		}
	}

	for index, data := range newsMap {
		fmt.Println(index, data)
	}

}
