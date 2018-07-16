package net

import (
	"sync"

	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/bloom"
	"gitlab.com/parallelcoin/duo/pkg/semaphore"
	"gitlab.com/parallelcoin/duo/pkg/serialize"
)

const (
	// LocalNone is
	LocalNone = iota
	// LocalIf is
	LocalIf
	// LocalBind is
	LocalBind
	// LocalUPNP is
	LocalUPNP
	// LocalHTTP is
	LocalHTTP
	// LocalManual is
	LocalManual
	// LocalMax is the end of the Local node types
	LocalMax
)

// Signals binds signals
type Signals struct {
}

// NodeStats is stores information about the state of a node
type NodeStats struct {
	Services                          uint64
	LastSend, LastRecv, TimeConnected int64
	Name                              string
	Version                           int
	Inbound                           bool
	StartingHeight, Misbehavior       int
	SendBytes, RecvBytes              uint64
	SyncNode                          bool
}

// Message is a message to be sent to another node
type Message struct {
	InData     bool
	DataStream ser.DataStream
	Header     MessageHeader
	Recv       ser.DataStream
	DataPos    uint
}

// Node is the complete data structure for a record for a peer in the network
type Node struct {
	Services                                                                           uint64
	Socket                                                                             uint
	DataStream                                                                         ser.DataStream
	SendSize, SendOffset                                                               uint
	SendBytes                                                                          uint64
	SendMsg                                                                            ser.Data
	SendMutex                                                                          sync.RWMutex
	RecvGetData                                                                        []Inv
	RecvMsg                                                                            []Message
	RecvMutex                                                                          sync.RWMutex
	RecvBytes                                                                          uint64
	RecvVersion                                                                        int
	LastSend, LastRecv, LastSendEmpty, TimeConnected                                   int64
	Address                                                                            Addr
	Name                                                                               string
	LocalAddr                                                                          Service
	Version                                                                            int
	SubVersion                                                                         string
	OneShot, Client, Inbound, NetworkNode, SuccessfullyConnected, Disconnect, RelayTXs bool
	GrantOutbound                                                                      semaphore.Grant
	Filter                                                                             bloom.Filter
	RefCount                                                                           int
	BannedSet                                                                          map[Addr]int64
	BannedMutex                                                                        sync.RWMutex
	Misbehavior                                                                        int
	HashContinue                                                                       Uint.U256
	IndexLastBlocksBegin                                                               *block.Index
	LastGetBlocksEnd                                                                   Uint.U256
	StartingHeight                                                                     int
	StartSync                                                                          bool
	AddrToSend                                                                         []Addr
	AddrKnown                                                                          []Addr
	GetAddr                                                                            bool
	Known                                                                              []Uint.U256
	InventoryKnown                                                                     []Inv
	InventoryToSend                                                                    []Inv
	InventoryMutex                                                                     sync.RWMutex
	AskForMap                                                                          map[int64]Inv
}
