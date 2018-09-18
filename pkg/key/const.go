package key
import (
	"encoding/hex"
)
var (
	B58prefixes = map[string]map[string]string{
		"mainnet": {
			"pubkey":  hex.EncodeToString([]byte{byte(83)}),
			"script":  hex.EncodeToString([]byte{byte(9)}),
			"privkey": hex.EncodeToString([]byte{byte(178)}),
		},
		"testnet": {
			"pubkey":  hex.EncodeToString([]byte{byte(18)}),
			"script":  hex.EncodeToString([]byte{byte(188)}),
			"privkey": hex.EncodeToString([]byte{byte(239)}),
		},
		"regtestnet": {
			"pubkey":  hex.EncodeToString([]byte{byte(0)}),
			"script":  hex.EncodeToString([]byte{byte(5)}),
			"privkey": hex.EncodeToString([]byte{byte(128)}),
		},
	}
)
