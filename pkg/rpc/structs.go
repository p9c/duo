package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/parallelcointeam/duo/pkg/core"
)

var er core.CommonErrors

const (
	// VERSION represents bicoind package version
	VERSION = 0.1
	// RPCClientTimeout represent http timeout for rcp client
	RPCClientTimeout = 30
)

// A Client is a connection to a websocket JSON RPC server
type Client struct {
	URL        string
	Username   string
	Password   string
	httpClient *http.Client
	core.State
}

// Request is a JSON RPC request
type Request struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

// Response is the response to a Request
type Response struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}
