// Package scoped store for information about the version of the duo golang parallelcoin blockchain network client
package version
import (
	"fmt"
)
const (
	Name          = "duod"
	Major         = 0
	Minor         = 1
	Revision      = 0
	Build         = 0
	Client        = 1000000*Major + 10000*Minor + 100*Revision + 1*Build
	IsRelease     = true
	CopyrightYear = 2018
	Protocol      = 80000
	MinProto      = 209
	AddrTime      = 31042
	NoBlocksStart = 32000
	NoBlocksEnd   = 32400
	BIP0031       = 60000
	MempoolGD     = 60002
	ClientName    = "Parallelcoin"
	Suffix        = "regenerator"
	Description   = "Client for Parallelcoin cryptocurrency network"
)
var (
	Git = map[string]string{
		"Commit": "main.GitCommit",
		"Branch": "main.GitBranch",
		"State":  "main.GitState",
	}
	BuildDate = "main.BuildDate"
)
func Print() {
	fmt.Printf("%s %d.%d.%d.%d - %s\n\n", Name, Major, Minor, Revision, Build, Description)
}
