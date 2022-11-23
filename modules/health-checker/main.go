package healthChecker

import (
	"bytes"
	"chainList"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

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
	bytes, err := os.ReadFile("chainlist.json")

	if err != nil {
		fmt.Println("There is no chainlist data file. The file will be regenerated.")
		chainList.Execute()
		bytes, err = os.ReadFile("chainlist.json")
	}

	return bytes, err
}

// * Export
func Execute(chainId uint64) (string, error) {
	bytes, _ := readChainListFile()

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
	// fmt.Println(sortedNodes)
	fmt.Println(">---------------------------------------------------------------<")

	return sortedNodes[0].Url, nil
}
