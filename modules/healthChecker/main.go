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

func getExtraChainInfo(chainId uint64) []gjson.Result {
	resp, err := http.Get(ExtraRPCUrl)

	CheckErr(err)
	defer resp.Body.Close()
	CheckCode(resp)

	body, err := ioutil.ReadAll(resp.Body)
	chainJson := string(body)

	chainInfoJSON := gjson.Parse(chainJson)

	nodes := gjson.Get(chainInfoJSON.String(), strconv.Itoa(int(chainId))+".rpcs").Array()

	return nodes
}

func getChainInfo() map[uint64][]gjson.Result {
	resp, err := http.Get(BaseUrl)

	CheckErr(err)
	defer resp.Body.Close()
	CheckCode(resp)

	body, err := ioutil.ReadAll(resp.Body)
	chainJson := string(body)

	bscChainInfoArr := gjson.Parse(chainJson).Array()

	var chainInfo = make(map[uint64][]gjson.Result)

	for _, data := range bscChainInfoArr {
		chainId := gjson.Get(data.String(), "chainId").Uint()
		rpc := gjson.Get(data.String(), "rpc").Array()
		chainInfo[chainId] = rpc
	}

	return chainInfo
}

func fetchNode(node string, c chan<- Node) {
	client := http.Client{
		Timeout: time.Second / 2,
	}
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

	body, _ := ioutil.ReadAll(resp.Body)
	bodyToString := string(body)
	hexHeight := gjson.Get(bodyToString, "result").String()
	cleanedHeight := strings.Replace(hexHeight, "0x", "", -1)
	intHeight, _ := strconv.ParseInt(cleanedHeight, 16, 64)

	c <- Node{node, latency.Milliseconds(), intHeight}
}

// * Export
func Execute(chainId uint64) (string, error) {
	chainInfos := getChainInfo()
	extraNodes := getExtraChainInfo(chainId)

	c := make(chan Node)
	nodes := []Node{}
	chainNodes := chainInfos[chainId]
	extraNodesLength := len(extraNodes)

	if extraNodesLength > 0 {
		chainNodes = append(chainNodes, extraNodes...)
	}

	nodesLength := len(chainNodes)

	if nodesLength > 0 {
		for _, node := range chainNodes {
			go fetchNode(node.String(), c)
		}

		for i := 0; i < nodesLength; i++ {
			fetchedNode := <-c
			if fetchedNode.Url != "" && fetchedNode.Height > 0 && fetchedNode.Latency > 0 {
				nodes = append(nodes, fetchedNode)
			}
		}

		if len(nodes) > 0 {
			sort.Sort(SortByLatencyWithHeight(nodes))
			fmt.Println(">-------- Sorted RPC Node --------------------------------------<")
			fmt.Println(nodes)
			fmt.Println(">---------------------------------------------------------------<")

			return nodes[0].Url, nil
		}

		return "", RPCDeadError
	}

	return "", InvalidChainError
}
