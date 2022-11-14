package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.saramin.co.kr/zf_user/jobs/list/job-category?page=1&cat_kewd=87&search_optional_item=n&search_done=y&panel_count=y&preview=y&isAjaxRequest=0&page_count=50&sort=RL&type=job-category&is_param=1&isSearchResultEmpty=1&isSectionHome=0&searchParamCount=1#searchTitle"

func main() {
	pages := getPages()
	fmt.Println(pages)
}

func getPages() int {
	resp, err := http.Get(baseURL)

	checkErr(err)
	defer resp.Body.Close()
	checkCode(resp)

	doc, queryErr := goquery.NewDocumentFromReader(resp.Body)
	checkErr(queryErr)
	fmt.Println(doc)

	return 0
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func checkCode(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalln("Request failed with status:", resp.StatusCode)
	}
}
