package chainList

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type NativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int64  `json:"decimals"`
}

type Chain struct {
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

var BaseUrl string = "https://chainid.network/chains.json"
var ExtraRPCUrl string = "https://raw.githubusercontent.com/DefiLlama/chainlist/main/constants/extraRpcs.json"

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

func fetchURL(url string) string {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)

	CheckErr(err)
	defer resp.Body.Close()
	CheckCode(resp)

	body, err := ioutil.ReadAll(resp.Body)
	chainJson := string(body)

	chainInfoJSON := gjson.Parse(chainJson).String()

	return chainInfoJSON
}

func toStringArray(g []gjson.Result) []string {
	result := []string{}
	for _, value := range g {
		result = append(result, value.String())
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

func Execute() {
	baseList := fetchURL(BaseUrl)
	extraList := fetchURL(ExtraRPCUrl)

	chainlist := []Chain{}
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
		_extraRpc := gjson.Get(extraList, strconv.Itoa(int(_chainId))+".rpcs").Array()

		_nativeCurrency := NativeCurrency{_nativeCurrencyName, _nativeCurrencySymbol, _nativeCurrencyDecimals}
		_rpc = append(_rpc, _extraRpc...)

		_strRpcs := toStringArray(_rpc)
		_strFaucets := toStringArray(_faucets)

		chainlist = append(chainlist, Chain{_name, _chain, _icon, _strRpcs, _strFaucets, _nativeCurrency, _infoURL, _shortName, _chainId, _networkId})
	}

	cbytes, _ := JSONMarshal(chainlist)

	os.WriteFile("chainlist.json", cbytes, 0666)
}
