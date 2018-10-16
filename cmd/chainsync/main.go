package main

import (
	"fmt"

	"github.com/parallelcointeam/duo/pkg/sync"
)

func main() {
	node := sync.NewNode()
	// node.RemoveOldVersions()
	node.Sync()
	if !node.OK() {
		fmt.Println(node.Error())
	}
	node.RemoveOldVersions()
	// fmt.Println(node.GetLatestSynced())
	// node.UpdateAddresses()
	node.Close()
}
