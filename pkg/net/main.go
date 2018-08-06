package net
import (
	"sync"
	"gitlab.com/parallelcoin/duo/pkg/semaphore"
	"gitlab.com/parallelcoin/duo/pkg/serialize"
)
const (
	// DumpAddressesInterval is the timeout after which an inactive node is removed from the list
	DumpAddressesInterval = 900
	// MaxOutboundConnections is the maximum number of connections we will make out to other nodes
	MaxOutboundConnections = 8
)
var (
	// Discover indicates that this server will learn its own external IP address
	Discover bool
	// LocalServices is
	LocalServices = NodeNetwork
	// LocalHostNonce is the distinct random number that identifies this node
	LocalHostNonce uint64
	// LocalHostMapMutex is the global mutex for this node
	LocalHostMapMutex sync.RWMutex
	// LocalHostMap is a list of local services
	LocalHostMap map[*Addr]*LocalServiceInfo
	// Reachable is a list of all the nodes that currently are reachable
	Reachable [NetMax]bool
	// Limited is the list of nodes that cannot receive inbound connections
	Limited [NetMax]bool
	// LocalHostNode is a node on the same machine
	LocalHostNode *Node
	// SyncNode is a node that we are syncing to
	SyncNode *Node
	// ListenSocket is the ports nodes are listening on
	ListenSocket []uint
	// Addrman is our Address Manager record
	Addrman AddrMan
	// MaxConnections is the maximum number of network connections we will establish in total
	MaxConnections = 125
	// Nodes is the list of all the nodes we know about
	Nodes []Node
	// NodesMutex must be locked to change the Nodes list
	NodesMutex sync.RWMutex
	// RelayMap is a stream of messages to be relayed
	RelayMap map[*Inv]*ser.DataStream
	// RelayExpiration records when the messages expire
	RelayExpiration map[int64]Inv
	// RelayMapMutex must be unlocked in order to change its contents
	RelayMapMutex sync.RWMutex
	// Oneshots are connections we will only try once
	Oneshots []string
	// OneshotsMutex must be unlocked to change the list of Oneshots
	OneshotsMutex sync.RWMutex
	// ServerNodeAddresses is a list of addresses to server nodes
	ServerNodeAddresses []Addr
	// AddedNodes are the list of addresses of nodes being added to our list
	AddedNodes []string
	// AddedNodesMutex must be unlocked in order to change AddedNodes
	AddedNodesMutex sync.RWMutex
	// OutboundSemaphore is used to send and receive status messages to/from other nodes
	OutboundSemaphore *semaphore.Semaphore
	// NodeSignals is
	NodeSignals Signals
	// NetCleanup is
	NetCleanup Cleanup
)
type LocalServiceInfo struct {
	Score, Port int
}
type Cleanup struct {
}
