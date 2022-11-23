package chainList

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
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

		_strRpcs := toStringArray(_rpc, true)
		_strFaucets := toStringArray(_faucets, false)

		chainlist = append(chainlist, Chain{_name, _chain, _icon, _strRpcs, _strFaucets, _nativeCurrency, _infoURL, _shortName, _chainId, _networkId})
	}

	cbytes, _ := JSONMarshal(chainlist)

	os.WriteFile("chainlist.json", cbytes, 0666)
}
