package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Location struct {
	Loc string `xml:"loc"`
}

type SiteMapIndex struct {
	Locations []Location `xml:"sitemap"`
}

func (location Location) String() string {
	return fmt.Sprintf(location.Loc)
}

func getBodyData(address string) {
	response, _ := http.Get(address)
	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	var siteMapIndex SiteMapIndex
	xml.Unmarshal(bytes, &siteMapIndex)

	fmt.Println(siteMapIndex.Locations)
}

func main() {
	getBodyData("https://www.washingtonpost.com/news-sitemap-index.xml")
}
