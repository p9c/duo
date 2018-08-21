package net
import (
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	"sync"
)
const (
	// NetUnrouteable means no network connection
	NetUnrouteable = iota
	// NetIPv4 means connected over IPv4
	NetIPv4
	// NetIPv6 means connected over IPv6
	NetIPv6
	// NetTor means connected over Tor
	NetTor
	// NetZeroTier means connected to the Parallelcoin virtual public Zerotier network
	NetZeroTier
	// NetMax is the limit of the type of connections list
	NetMax
)
var (
	// ProxyType is a map of codes representing different types of network connections
	ProxyType = map[string]int{
		"Unrouteable": 0,
		"IPv4":        1,
		"IPv6":        2,
		"Tor":         3,
		"Unknown":     4,
		"Teredo":      5,
		"Zerotier":    6,
	}
	// ProxyInfosMutex must be unlocked to change information about proxy connections
	ProxyInfosMutex sync.RWMutex
	// ConnectTimeout is how long a connection can remain unresponsive before we disconnect
	ConnectTimeout = 5000
	// NameLookup sets whether we will use DNS to resolve addresses
	NameLookup = false
	// IPv4Mask is a bitmask over 128 bits
	IPv4Mask = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}
	// OnionMask is the bitmask for a Tor hidden address
	OnionMask = []byte{0xFD, 0x87, 0xD8, 0x7E, 0xEB, 0x43}
)
// Addr stores a peer network address
type Addr struct {
	IP [16]byte
}
// Service stores the port of a service
type Service struct {
	Port uint16
}
func (a *Addr) IsValid() bool {
	// placeholder
	return true
}
// ToStringIPPort returns the string for conneecting to a service
func (s *Service) ToStringIPPort() string {
	// placeholder
	return "127.0.0.1:" + fmt.Sprint(*args.Proxy)
}
