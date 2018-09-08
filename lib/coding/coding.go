package coding

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/base58check"
	"math/big"
	"strconv"
	"strings"
)

// Coding is an interface for encoding raw bytes in various base number formats
type Coding interface {
	GetCoding() string
	SetCoding(string) Coding
	ListCodings() []string
}

// Codings are the types of encoding schemes available
var Codings = []string{"string", "golang", "octal", "decimal", "hex", "base32", "base58check", "base64"}

// Encode takes a byte slice and encodes it to a string with the prescribed format or if empty or not in the list, defaults to string
func Encode(b []byte, code string, extra ...int) string {
	switch code {
	case "golang":
		return fmt.Sprint(b)
	case "octal":
		bi := big.NewInt(0)
		bi.SetBytes(b)
		return bi.Text(8)
	case "decimal":
		bi := big.NewInt(0)
		bi.SetBytes(b)
		return bi.Text(10)
	case "hex":
		return hex.EncodeToString(b)
	case "base32":
		return strings.ToLower(base32.StdEncoding.EncodeToString(b))
	case "base58check":
		if len(extra) == 0 {
			extra = []int{0}
		}
		pre := hex.EncodeToString([]byte{byte(extra[0])})
		msg := hex.EncodeToString(b)
		r, _ := base58check.Encode(pre, msg)
		return r
	case "base64":
		return base64.StdEncoding.EncodeToString(b)
	default: // string
		return string(b)
	}
}

// Decode takes a byte slice and encodes it to a string with the prescribed format or if empty or not in the list, defaults to string
func Decode(s string, code string) (b []byte) {
	if len(s) == 0 {
		return []byte{}
	}
	switch code {
	case "golang":
		split1 := strings.Split(s, "[")
		split2 := strings.Split(split1[1], "]")
		splitted := strings.Split(split2[0], " ")
		b = make([]byte, len(splitted))
		for i := range splitted {
			r, _ := strconv.Atoi(splitted[i])
			b[i] = byte(r)
		}
		return b
	case "octal":
		bi := big.NewInt(0)
		bi.SetString(s, 8)
		return bi.Bytes()
	case "decimal":
		bi := big.NewInt(0)
		bi.SetString(s, 10)
		return bi.Bytes()
	case "hex":
		bi := big.NewInt(0)
		bi.SetString(s, 16)
		return bi.Bytes()
	case "base32":
		var r []byte
		r, _ = base32.StdEncoding.DecodeString(strings.ToUpper(s))
		return r
	case "base58check":
		msg, err := base58check.Decode(s)
		if err != nil {
			return []byte("Base58check decoding error")
		}
		r, _ := hex.DecodeString(msg[2:])
		return []byte(r)
	case "base64":
		b, _ = base64.StdEncoding.DecodeString(s)
		return b
	default:
		return []byte(s)
	}
}
