package net
import (
	"sync"
)
// AddrInfo stores all the information about a peer
type AddrInfo struct {
	Addr
	source                        Addr
	lastSuccess                   int64
	attempts, refCount, randomPos int
	inTried                       bool
}
type AddrMan struct {
	mutex    sync.RWMutex
	key      []byte
	idCount  int
	infoMap  map[int]Addr
	addrMap  map[*Addr]int
	random   []int
	tried    int
	triedSet [][]int
	new      int
	newSet   [][]int
}
