package healthChecker

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
