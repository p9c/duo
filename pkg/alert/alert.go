package alert

import (
	"sync"

	"gitlab.com/parallelcoin/duo/pkg/Uint"
)

var (
	// Map is
	Map map[*Uint.U256]*Alert
	// Mutex is
	Mutex sync.RWMutex
)

// Unsigned is an unsigned alert
type Unsigned struct {
	Version                      int
	RelayUntil, Expiration       int64
	ID, Cancel                   int
	SetCancel                    []int
	MinVer, MaxVer               int
	SubVer                       []string
	Priority                     int
	Comment, StatusBar, Reserved string
}

// Alert is a signed alert
type Alert struct {
	Unsigned
	Msg, Sig []byte
}
