package rpc

// GetBlock is the result of the getblock API call with true (or no) second parameter
type GetBlock struct {
	Hash              string   `json:"hash"`
	Confirmations     uint32   `json:"confirmations"`
	Size              uint32   `json:"size"`
	Height            uint32   `json:"height"`
	Version           uint32   `json:"version"`
	PowAlgoID         uint32   `json:"pow_algo_id"`
	PowAlgo           string   `json:"pow_algo"`
	PowHash           string   `json:"pow_hash"`
	MerkleRoot        string   `json:"merkleroot"`
	Tx                []string `json:"tx"`
	Time              int64    `json:"time"`
	Nonce             uint64   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	PreviousBlockHash string   `json:"previousblockhash"`
	NextBlockHash     string   `json:"nextblockhash"`
}
