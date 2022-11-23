package healthChecker

import "errors"

var RPCDeadError error = errors.New("All RPC Nodes are Dead.")

var InvalidChainError error = errors.New("Invalid Chain ID.")
