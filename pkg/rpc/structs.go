package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/parallelcointeam/duo/pkg/core"
)

var er core.Errors

type Client struct {
	URL        string
	Username   string
	Password   string
	httpClient *http.Client
	core.State
}

type Request struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

type Response struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}
