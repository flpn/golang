package main

import (
	"encoding/xml"
	"html/template"
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

type NewsAggregatorPage struct {
	Title string
	News  map[string]NewsMap
}

func NewsAggregatorHandler(w http.ResponseWriter, r *http.Request) {
	var siteMapIndex SiteMapIndex
	var news News
	newsMap := make(map[string]NewsMap)

	response, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(response.Body)
	xml.Unmarshal(bytes, &siteMapIndex)

	for _, Location := range siteMapIndex.Locations {
		response, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(response.Body)

		xml.Unmarshal(bytes, &news)

		for index := range news.Keywords {
			newsMap[news.Titles[index]] = NewsMap{news.Keywords[index], news.Locations[index]}
		}
	}

	data := NewsAggregatorPage{"Amazing News Aggregator", newsMap}
	temp, _ := template.ParseFiles("index.html")
	temp.Execute(w, data)
}

func main() {
	http.HandleFunc("/", NewsAggregatorHandler)
	http.ListenAndServe(":8000", nil)
}
