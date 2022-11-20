package healthChecker

import "strconv"

type RpcRequest struct {
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
