package healthChecker

import "errors"

var BaseUrl string = "https://chainid.network/chains.json"

var ExtraRPCUrl string = "https://raw.githubusercontent.com/DefiLlama/chainlist/main/constants/extraRpcs.json"

var RPCDeadError error = errors.New("All RPC Nodes are Dead.")

var InvalidChainError error = errors.New("Invalid Chain ID.")
