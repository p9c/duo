package net

import (
	"unsafe"

	"gitlab.com/parallelcoin/duo/pkg/Uint"
)

const (
	// NodeNetwork is
	NodeNetwork = 1
	// CommandSize is
	CommandSize = 12
	// Int is
	Int int = 0
	// MessageSizeSize is
	MessageSizeSize = unsafe.Sizeof(Int)
	// MessageTx is
	MessageTx = 1
	// MessageBlock is
	MessageBlock = 2
	// MessageFilteredBlock is
	MessageFilteredBlock = 3
)

var (
	// TypeName is
	TypeName = []string{
		"Error",
		"tx",
		"block",
		"filtered block",
	}
)

// MessageHeader is the header for a p2p message
type MessageHeader struct {
	MessageStart          [MessageStartSize]byte
	Command               [CommandSize]byte
	MessageSize, Checksum uint
}

// Address is the address record of a node
type Address struct {
	Service
	Services uint64
	Time     uint
	LastTry  int64
}

// Inv is the inventory of data in a message
type Inv struct {
	Type int
	Hash Uint.U256
}
