package sync

import (
	"github.com/parallelcointeam/duo/pkg/rpc"

	"github.com/dgraph-io/badger"
)

// Node is a sync client that updates by polling from full node, but will eventually implement as a standalone p2p client
type Node struct {
	RPC *rpc.Client
	DB  *badger.DB
}
