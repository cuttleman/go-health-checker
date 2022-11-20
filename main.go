package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

var baseUrl string = "https://chainid.network/chains.json"
var extraRPCUrl string = "https://raw.githubusercontent.com/DefiLlama/chainlist/main/constants/extraRpcs.json"

type rpcRequest struct {
	Id      int8
	Jsonrpc string
	Method  string
	Params  []interface{}
}

type Node struct {
	Url, Latency, Height string
}

type SortByLatencyWithHeight []Node

func (a SortByLatencyWithHeight) Len() int      { return len(a) }
func (a SortByLatencyWithHeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByLatencyWithHeight) Less(i, j int) bool {
	h1, _ := strconv.Atoi(a[i].Height)
	h2, _ := strconv.Atoi(a[j].Height)
	l1, _ := strconv.Atoi(a[i].Latency)
	l2, _ := strconv.Atoi(a[j].Latency)

	if h2-h1 > 0 {
		return false
	}
	if h2-h1 < 0 {
		return true
	}
	if h1 == h2 {
		if l1-l2 < 0 {
			return true
		} else {
			return false
		}
	}
	return true
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

func makeJson(data gjson.Result) (uint64, []gjson.Result) {
	chainId := gjson.Get(data.String(), "chainId").Uint()
	rpc := gjson.Get(data.String(), "rpc").Array()

	return chainId, rpc
}

func getExtraChainInfo(chainId uint64) []gjson.Result {
	resp, err := http.Get(extraRPCUrl)

	checkErr(err)
	defer resp.Body.Close()
	checkCode(resp)

	body, err := ioutil.ReadAll(resp.Body)
	chainJson := string(body)

	chainInfoJSON := gjson.Parse(chainJson)

	nodes := gjson.Get(chainInfoJSON.String(), strconv.Itoa(int(chainId))+".rpcs").Array()

	return nodes
}

func getChainInfo() map[uint64][]gjson.Result {
	resp, err := http.Get(baseUrl)

	checkErr(err)
	defer resp.Body.Close()
	checkCode(resp)

	body, err := ioutil.ReadAll(resp.Body)
	chainJson := string(body)

	bscChainInfoArr := gjson.Parse(chainJson).Array()

	var chainInfo = make(map[uint64][]gjson.Result)

	for _, data := range bscChainInfoArr {
		chainId, rpc := makeJson(data)
		chainInfo[chainId] = rpc
	}

	return chainInfo
}

func fetchNode(node string) (Node, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	parameters := rpcRequest{1, "2.0", "eth_blockNumber", []interface{}{}}
	pbytes, _ := json.Marshal(parameters)
	pbuff := bytes.NewBuffer(pbytes)

	start := time.Now()
	resp, err := client.Post(node, "application/json", pbuff)
	latency := time.Since(start)

	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyToString := string(body)
		height := gjson.Get(bodyToString, "result").String()
		return Node{node, strconv.Itoa(int(latency.Milliseconds())), height}, nil
	}

	return Node{"", "", ""}, err
}

func healthChecker(chainId uint64) (string, error) {
	chainInfos := getChainInfo()
	extraNodes := getExtraChainInfo(chainId)

	nodes := []Node{}
	chainNodes := chainInfos[chainId]
	extraNodesLength := len(extraNodes)

	if extraNodesLength > 0 {
		chainNodes = append(chainNodes, extraNodes...)
	}

	nodesLength := len(chainNodes)

	if nodesLength > 0 {
		for _, node := range chainNodes {
			fetchedNode, err := fetchNode(node.String())

			if err == nil && fetchedNode.Url != "" && fetchedNode.Height != "" && fetchedNode.Latency != "" {
				nodes = append(nodes, fetchedNode)
			}
		}

		sort.Sort(SortByLatencyWithHeight(nodes))

		return nodes[0].Url, nil
	}

	return "", errors.New("Invalid Chain ID:" + strconv.Itoa(int(chainId)))
}

func main() {
	greatNode, _ := healthChecker(97)
	fmt.Println("greatNode :", greatNode)
}
