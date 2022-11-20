package healthChecker

import (
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func CheckCode(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalln("Request failed with status:", resp.StatusCode)
	}
}
