package main

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

var waitGroup sync.WaitGroup

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

func newsRoutine(channel chan News, Location string) {
	defer waitGroup.Done()

	var news News

	response, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(response.Body)

	xml.Unmarshal(bytes, &news)
	response.Body.Close()

	channel <- news

}

func NewsAggregatorHandler(w http.ResponseWriter, r *http.Request) {
	var siteMapIndex SiteMapIndex
	newsMap := make(map[string]NewsMap)
	channel := make(chan News, 30)

	response, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(response.Body)
	xml.Unmarshal(bytes, &siteMapIndex)
	response.Body.Close()

	for _, Location := range siteMapIndex.Locations {
		waitGroup.Add(1)
		go newsRoutine(channel, Location)
	}

	waitGroup.Wait()
	close(channel)

	for elem := range channel {
		for index := range elem.Keywords {
			newsMap[elem.Titles[index]] = NewsMap{elem.Keywords[index], elem.Locations[index]}
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
