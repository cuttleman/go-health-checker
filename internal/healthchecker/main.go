package healthchecker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"healthchecker-server/internal/chainlist"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type RpcRequest struct {
	Id      int8          `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Node struct {
	Url     string
	Latency int64
	Height  int64
}

type SortByLatencyWithHeight []Node

func (a SortByLatencyWithHeight) Len() int      { return len(a) }
func (a SortByLatencyWithHeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByLatencyWithHeight) Less(i, j int) bool {
	h1 := a[i].Height
	h2 := a[j].Height
	l1 := a[i].Latency
	l2 := a[j].Latency

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

var RPCDeadError error = errors.New("All RPC Nodes are Dead.")
var InvalidChainError error = errors.New("Invalid Chain ID.")

func fetchNode(node string, c chan<- Node) {
	client := http.Client{Timeout: time.Second / 2} // 500ms
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

func readChainListFile() ([]byte, error) {
	bytes, err := os.ReadFile(chainlist.ChainListPath)

	if err != nil {
		fmt.Println("There is no chainlist data file. The file will be regenerated.")
		chainlist.Execute()
		bytes, err = os.ReadFile(chainlist.ChainListPath)
	}

	return bytes, err
}

// * Export
func Execute(chainId uint64) (string, error) {
	bytes, err := readChainListFile()

	if err != nil {
		return "", err
	}

	bytesToString := gjson.Parse(string(bytes)).String()
	chainRpcs := gjson.Get(bytesToString, "#(chainId=="+strconv.Itoa(int(chainId))+").rpc").Array()
	fmt.Println(chainRpcs)

	chainUrlLength := len(chainRpcs)

	if chainUrlLength == 0 {
		return "", InvalidChainError
	}

	c := make(chan Node)
	nodes := []Node{}

	for _, chainUrl := range chainRpcs {
		go fetchNode(chainUrl.String(), c)
	}

	fetchNodeStart := time.Now()
	for i := 0; i < chainUrlLength; i++ {
		fetchedNode := <-c
		if fetchedNode.Url != "" && fetchedNode.Height > 0 {
			nodes = append(nodes, fetchedNode)
		}
	}
	fetchNodeResponseTime := time.Since(fetchNodeStart)
	fmt.Println("Parallel Node Response Time :", fetchNodeResponseTime)

	if len(nodes) == 0 {
		return "", RPCDeadError
	}

	sortedNodes := sortNodes(nodes)
	fmt.Println(">-------- Sorted RPC Node --------------------------------------<")
	fmt.Println(sortedNodes)
	fmt.Println(">---------------------------------------------------------------<")

	return sortedNodes[0].Url, nil
}
