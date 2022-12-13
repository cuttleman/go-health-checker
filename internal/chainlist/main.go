package chainlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type NativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int64  `json:"decimals"`
}

type ChainInfo struct {
	Name           string         `json:"name"`
	Chain          string         `json:"chain"`
	Icon           string         `json:"icon"`
	Rpc            []string       `json:"rpc"`
	Faucets        []string       `json:"faucets"`
	NativeCurrency NativeCurrency `json:"nativeCurrency"`
	InfoURL        string         `json:"infoURL"`
	ShortName      string         `json:"shortName"`
	ChainId        int64          `json:"chainId"`
	NetworkId      int64          `json:"networkId"`
}

const (
	BaseUrl       = "https://chainid.network/chains.json"
	ExtraRPCUrl   = "https://raw.githubusercontent.com/DefiLlama/chainlist/main/constants/extraRpcs.js"
	AssetsDir     = "assets"
	ChainListPath = AssetsDir + "/chainlist.json"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func CheckCode(resp *http.Response) {
	if resp.StatusCode != 200 {
		fmt.Println("Request failed with status:", resp.StatusCode)
	}
}

func fetchURL(url string) string {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)

	CheckErr(err)
	defer resp.Body.Close()
	CheckCode(resp)

	body, err := ioutil.ReadAll(resp.Body)
	chainJson := string(body)
	if strings.HasPrefix(chainJson, "export default") {
		chainJson = chainJson[len("export default"):]
	}

	return chainJson
}

func toStringArray(g []gjson.Result, rTrailingSlash bool) []string {
	allKeys := make(map[string]bool)
	result := []string{}
	for _, item := range g {
		strItem := item.String()
		if strings.HasPrefix(strItem, "wss://") {
			continue
		}
		if rTrailingSlash && strings.HasSuffix(strItem, "/") {
			strItem = strItem[:len(strItem)-1]
		}
		if _, value := allKeys[strItem]; !value {
			allKeys[strItem] = true
			result = append(result, strItem)
		}
	}

	return result
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(t)

	return buf.Bytes(), err
}

func Execute() error {
	baseList := fetchURL(BaseUrl)
	extraList := fetchURL(ExtraRPCUrl)

	chainlist := []ChainInfo{}
	for _, chain := range gjson.Parse(baseList).Array() {
		_name := gjson.Get(chain.String(), "name").String()
		_chain := gjson.Get(chain.String(), "chain").String()
		_icon := gjson.Get(chain.String(), "icon").String()
		_faucets := gjson.Get(chain.String(), "faucets").Array()
		_nativeCurrencyName := gjson.Get(chain.String(), "nativeCurrency.name").String()
		_nativeCurrencySymbol := gjson.Get(chain.String(), "nativeCurrency.symbol").String()
		_nativeCurrencyDecimals := gjson.Get(chain.String(), "nativeCurrency.decimals").Int()
		_infoURL := gjson.Get(chain.String(), "infoURL").String()
		_shortName := gjson.Get(chain.String(), "shortName").String()
		_chainId := gjson.Get(chain.String(), "chainId").Int()
		_networkId := gjson.Get(chain.String(), "networkId").Int()
		_rpc := gjson.Get(chain.String(), "rpc").Array()
		_extraRpcArr := gjson.Get(extraList, strconv.Itoa(int(_chainId))+".rpcs").Array()
		_extraRpc := []gjson.Result{}
		for _, extra := range _extraRpcArr {
			if extra.Type == gjson.String && extra.String() != "rpcWorking:false" {
				_extraRpc = append(_extraRpc, extra)
			}
		}

		_nativeCurrency := NativeCurrency{_nativeCurrencyName, _nativeCurrencySymbol, _nativeCurrencyDecimals}
		_rpc = append(_rpc, _extraRpc...)

		_strRpcs := toStringArray(_rpc, true)
		_strFaucets := toStringArray(_faucets, false)

		chainlist = append(chainlist, ChainInfo{_name, _chain, _icon, _strRpcs, _strFaucets, _nativeCurrency, _infoURL, _shortName, _chainId, _networkId})
	}

	cbytes, _ := JSONMarshal(chainlist)

	if _, readDirErr := os.ReadDir(AssetsDir); readDirErr != nil {
		os.Mkdir(AssetsDir, os.ModePerm)
		fmt.Println("The assets directory did not exist, so it was created.")
	}

	err := os.WriteFile(ChainListPath, cbytes, 0666)

	return err
}
