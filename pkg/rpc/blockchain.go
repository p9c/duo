package rpc

import (
	"gitlab.com/parallelcoin/duo/pkg/algos"
	"gitlab.com/parallelcoin/duo/pkg/block"
)

func getblockcount(help bool, params ...string) interface{} {
	return nil
}
func getbestblockhash(help bool, params ...string) interface{} {
	return nil
}
func getdifficulty(help bool, params ...string) interface{} {
	return nil
}
func GetDifficulty(index block.Index, algo int) float64 {
	switch algo {
	case algos.SHA256D:
	case algos.SCRYPT:
	}
	return 0.0
}
func settxfee(help bool, params ...string) interface{} {
	return nil
}
func getrawmempool(help bool, params ...string) interface{} {
	return nil
}
func getblockhash(help bool, params ...string) interface{} {
	return nil
}
func getblock(help bool, params ...string) interface{} {
	return nil
}
func gettxoutsetinfo(help bool, params ...string) interface{} {
	return nil
}
func gettxout(help bool, params ...string) interface{} {
	return nil
}
func verifychain(help bool, params ...string) interface{} {
	return nil
}
