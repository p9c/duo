package main

import (
	"github.com/parallelcointeam/duo/pkg/sync"
)

func main() {
	node := sync.NewNode()
	node.Sync()

}
