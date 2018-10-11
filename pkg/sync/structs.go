package sync

import (
	"os"

	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/rpc"

	"github.com/dgraph-io/badger"
)

// Node is a sync client that updates by polling from full node, but will eventually implement as a standalone p2p client
type Node struct {
	RPC        *rpc.Client
	DB         *badger.DB
	Chain      *os.File
	Latest     uint32
	LatestHash []byte
	Best       uint32
	BestTime   int64
	core.State
}
