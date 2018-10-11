package rpc

// GetInfo is the response from getinfo
type GetInfo struct {
	Version           uint32  `json:"version"`
	ProtocolVersion   uint32  `json:"protocolversion"`
	WalletVersion     uint32  `json:"walletversion"`
	Balance           float64 `json:"balance"`
	Blocks            uint32  `json:"blocks"`
	TimeOffset        int64   `json:"timeoffset"`
	Connections       uint32  `json:"connections"`
	Proxy             string  `json:"proxy"`
	PoWAlgoID         uint32  `json:"pow_algo_id"`
	PoWAlgo           string  `json:"pow_algo"`
	Difficulty        float64 `json:"difficulty"`
	DifficultySHA256d float64 `json:"difficulty_sha256d"`
	DifficultyScrypt  float64 `json:"difficulty_scrypt"`
	Testnet           bool    `json:"testnet"`
	KeyPoolOldest     int64   `json:"keypoololdest"`
	KeyPoolSize       uint32  `json:"keypoolsize"`
	PayTxFee          float64 `json:"paytxfee"`
	Errors            string  `json:"errors"`
}
