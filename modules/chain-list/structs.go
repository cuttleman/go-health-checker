package chainList

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
