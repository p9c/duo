package rpc

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
)

var (
	TallyItems Tally
)

type Tally struct {
	Amount int64
	Conf   int
	TxIDs  []Uint.U256
}

func getnewaddress(help bool, params ...string) interface{} {
	return nil
}
func getaccountaddress(help bool, params ...string) interface{} {
	return nil
}
func setaccount(help bool, params ...string) interface{} {
	return nil
}
func getaccount(help bool, params ...string) interface{} {
	return nil
}
func getaddressesbyaccount(help bool, params ...string) interface{} {
	return nil
}
func sendtoaddress(help bool, params ...string) interface{} {
	return nil
}
func signmessage(help bool, params ...string) interface{} {
	return nil
}
func verifymessage(help bool, params ...string) interface{} {
	return nil
}
func getreceivedbyaddress(help bool, params ...string) interface{} {
	return nil
}
func getreceivedbyaccount(help bool, params ...string) interface{} {
	return nil
}
func getbalance(help bool, params ...string) interface{} {
	return nil
}
func movecmd(help bool, params ...string) interface{} {
	return nil
}
func sendfrom(help bool, params ...string) interface{} {
	return nil
}
func sendmany(help bool, params ...string) interface{} {
	return nil
}
func addmultisigaddress(help bool, params ...string) interface{} {
	return nil
}
func createmultisig(help bool, params ...string) interface{} {
	return nil
}
func listreceivedbyaddress(help bool, params ...string) interface{} {
	return nil
}
func listreceivedbyaccount(help bool, params ...string) interface{} {
	return nil
}
func listtransactions(help bool, params ...string) interface{} {
	return nil
}
func listaddressgroupings(help bool, params ...string) interface{} {
	return nil
}
func listaccounts(help bool, params ...string) interface{} {
	return nil
}
func listsinceblock(help bool, params ...string) interface{} {
	return nil
}
func gettransaction(help bool, params ...string) interface{} {
	return nil
}
func backupwallet(help bool, params ...string) interface{} {
	return nil
}
func keypoolrefill(help bool, params ...string) interface{} {
	return nil
}
func walletpassphrase(help bool, params ...string) interface{} {
	return nil
}
func walletpassphrasechange(help bool, params ...string) interface{} {
	return nil
}
func walletlock(help bool, params ...string) interface{} {
	return nil
}
func encryptwallet(help bool, params ...string) interface{} {
	return nil
}
func validateaddress(help bool, params ...string) interface{} {
	return nil
}
func getinfo(help bool, params ...string) interface{} {
	return nil
}

func makekeypair(help bool, params ...string) interface{} {
	return nil
}
