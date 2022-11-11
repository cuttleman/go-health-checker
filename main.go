package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type result struct {
	url  string
	code int
}

func main() {
	c := make(chan *result)
	results := make(map[string]int)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://academy.nomadcoders.co/",
	}

	for _, url := range urls {
		go urlChecker(url, c)
	}

	for i := 0; i < len(urls); i++ {
		resulted := <-c
		results[resulted.url] = resulted.code
	}

	for url, code := range results {
		fmt.Println(url, code)
	}
}

func urlChecker(url string, c chan<- *result) {
	resp, err := http.Get(url)

	if err != nil {
		c <- &result{url: url, code: 0}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	fmt.Println(doc)
	c <- &result{url: url, code: resp.StatusCode}
}
