package net

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/algos"
	"gitlab.com/parallelcoin/duo/pkg/block"
)

const (
	// MessageStartSize is the number of bytes of the distinctive message header of this network
	MessageStartSize = 4
	// NetworkMainnet indicates the server is running on the main net
	NetworkMainnet = iota
	// NetworkTestnet indicates the server is running on the Testnet
	NetworkTestnet
	// NetworkRegtest indicates the server is running on the Regression Testnet
	NetworkRegtest
	// Base58PubKeyAddress indicates that an address is a public key
	Base58PubKeyAddress = iota
	// Base58ScriptAddress indicates that an address is a script
	Base58ScriptAddress
	// Base58SecretKey indicates that an address is a secret key
	Base58SecretKey
	// Base58MaxTypes is the bound of the Base58 address types
	Base58MaxTypes
)

var (
	// MainNetMessageStart is the set of 4 bytes that indicates the beginning of a message
	MainNetMessageStart = []byte{0xcd, 0x08, 0xac, 0xff}
	// MainNetParams is is the parameters of the main net
	MainNetParams MainParams
	// TestNetParams is the parameters of the Testnet
	TestNetParams TestParams
	// RegTestNetParams is the parameters of the Regression Testnet
	RegTestNetParams RegTestParams
	// CurrentParams is set by default to mainnet
	CurrentParams = &MainNetParams
)

// DNSSeedData stores the information about a DNS seed
type DNSSeedData struct {
	Name, Host string
}

// MessageStartChars is the marker indicating messages are for this network
type MessageStartChars [MessageStartSize]byte

// Params is the parameters of a blockchain network we are connecting to
type Params struct {
	HashGenesisBlock       Uint.U256
	MessageStart           MessageStartChars
	AlertPubKey            []byte
	DefaultPort, RPCPort   int
	ProofOfWorkLimit       [algos.NumAlgos][32]byte
	SubsidyHalvingInterval int
	DataDir                string
	Seeds                  []DNSSeedData
	Base58Prefixes         [Base58MaxTypes]int
}

// MainParams is the parameters of the mainnet
type MainParams struct {
	Genesis    block.Block
	FixedSeeds Address
}

// TestParams is the parametrs of the Testnet
type TestParams struct {
	MainParams
}

// RegTestParams is the parameters of the Regression Testnet
type RegTestParams struct {
	TestParams
}
