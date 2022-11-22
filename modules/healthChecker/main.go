package healthChecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func fetchRPC(node string, c chan<- Node) {
	client := http.Client{Timeout: time.Millisecond * 500} // 500ms
	parameters := RpcRequest{1, "2.0", "eth_blockNumber", []interface{}{}}
	pbytes, _ := json.Marshal(parameters)
	pbuff := bytes.NewBuffer(pbytes)

	start := time.Now()
	resp, err := client.Post(node, "application/json", pbuff)
	latency := time.Since(start)

	if err != nil {
		c <- Node{"", 0, 0}

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c <- Node{"", 0, 0}

		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyToString := string(body)
	hexHeight := gjson.Get(bodyToString, "result").String()
	cleanedHeight := strings.Replace(hexHeight, "0x", "", -1)
	intHeight, _ := strconv.ParseInt(cleanedHeight, 16, 64)

	c <- Node{node, latency.Milliseconds(), intHeight}
}

func sortNodes(nodes []Node) []Node {
	sort.Sort(SortByLatencyWithHeight(nodes))

	return nodes
}

func fetchURL(url string, chainId uint64, c chan<- []gjson.Result) {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)

	CheckErr(err)
	defer resp.Body.Close()
	CheckCode(resp)

	body, err := ioutil.ReadAll(resp.Body)
	chainJson := string(body)

	chainInfoJSON := gjson.Parse(chainJson)

	rpcs := make([]gjson.Result, 0)
	if url == BaseUrl {
		rpcs = gjson.Get(chainInfoJSON.String(), "#(chainId=="+strconv.Itoa(int(chainId))+").rpc").Array()
	} else {
		rpcs = gjson.Get(chainInfoJSON.String(), strconv.Itoa(int(chainId))+".rpcs").Array()
	}

	c <- rpcs
}

// * Export
func Execute(chainId uint64) (string, error) {
	var resourceUrls [2]string = [2]string{BaseUrl, ExtraRPCUrl}

	c1 := make(chan []gjson.Result)
	chainUrls := []gjson.Result{}

	for _, url := range resourceUrls {
		go fetchURL(url, chainId, c1)
	}

	fetchURLStart := time.Now()
	for i := 0; i < len(resourceUrls); i++ {
		urls := <-c1
		chainUrls = append(chainUrls, urls...)
	}
	fetchURLResponseTime := time.Since(fetchURLStart)
	fmt.Println("\nParallel URL Response Time :", fetchURLResponseTime)

	// **********************************
	// * get rpc urls ↑
	// * rpc status check ↓
	// **********************************

	chainUrlLength := len(chainUrls)

	if chainUrlLength == 0 {
		return "", InvalidChainError
	}

	c2 := make(chan Node)
	nodes := []Node{}

	for _, chainUrl := range chainUrls {
		go fetchRPC(chainUrl.String(), c2)
	}

	fetchRPCStart := time.Now()
	for i := 0; i < chainUrlLength; i++ {
		fetchedRPC := <-c2
		if fetchedRPC.Url != "" && fetchedRPC.Height > 0 {
			nodes = append(nodes, fetchedRPC)
		}
	}
	fetchRPCResponseTime := time.Since(fetchRPCStart)
	fmt.Println("Parallel RPC Response Time :", fetchRPCResponseTime)

	if len(nodes) == 0 {
		return "", RPCDeadError
	}

	sortedNodes := sortNodes(nodes)
	fmt.Println(">-------- Sorted RPC Node --------------------------------------<")
	fmt.Println(sortedNodes)
	fmt.Println(">---------------------------------------------------------------<")

	return sortedNodes[0].Url, nil

}
