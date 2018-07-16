package rpc

import (
	"sync"
)

var (
	HTTPStatus = map[string]int{
		"OK":                  200,
		"BadRequest":          400,
		"Unauthorized":        401,
		"Forbidden":           403,
		"NotFound":            404,
		"InternalServerError": 500,
	}
	RPCError = map[string]int{
		// Standard JSON-RPC 2.0 errors
		"InvalidRequest": -32600,
		"MethodNotFound": -32601,
		"InvalidParams":  -32602,
		"InternalError":  -32603,
		"ParseError":     -32700,
		// General application defined errors
		"MiscError":            -1,
		"ForbiddenBySafeMode":  -2,
		"TypeError":            -3,
		"InvalidAddressOrKey":  -5,
		"OutOfMemory":          -7,
		"InvalidParameter":     -8,
		"DatabaseError":        -20,
		"DeserializationError": -22,
		// P2P client errors
		"ClientNotConnected":      -9,
		"ClientInInitialDownload": -10,
		// Wallet errors
		"WalletError":               -4,
		"WalletInsufficientFunds":   -6,
		"WalletInvalidAccounntName": -11,
		"WalletKeypoolRanOut":       -12,
		"WalletUnlockNeeded":        -13,
		"WalletPassphraseIncorrect": -14,
		"WalletWrongEncState":       -15,
		"WalletEncryptionFailed":    -16,
		"WalletAlreadyUnlocked":     -17,
	}
)

type RPCCommand struct {
	Name                   string
	Actor                  func(help bool, params ...string) interface{}
	OKSafeMode, ThreadSafe bool
}

var (
	RPCUserColonPass      string
	WalletUnlockTime      int64
	WalletUnlockTimeMutex sync.RWMutex
	RPCCommands           = []RPCCommand{
		{"help", help, true, true},
		{"stop", stop, true, true},
		{"getblockcount", getblockcount, true, false},
		{"getbestblockhash", getbestblockhash, true, false},
		{"getconnectioncount", getconnectioncount, true, false},
		{"getpeerinfo", getpeerinfo, true, false},
		{"addnode", addnode, true, true},
		{"getaddednodeinfo", getaddednodeinfo, true, true},
		{"getdifficulty", getdifficulty, true, false},
		{"getnetworkhashps", getnetworkhashps, true, false},
		{"getgenerate", getgenerate, true, false},
		{"setgenerate", setgenerate, true, false},
		{"gethashespersec", gethashespersec, true, false},
		{"getinfo", getinfo, true, false},
		{"getmininginfo", getmininginfo, true, false},
		{"getnewaddress", getnewaddress, true, false},
		{"getaccountaddress", getaccountaddress, true, false},
		{"setaccount", setaccount, true, false},
		{"getaccount", getaccount, false, false},
		{"getaddressesbyaccount", getaddressesbyaccount, true, false},
		{"sendtoaddress", sendtoaddress, false, false},
		{"getreceivedbyaddress", getreceivedbyaddress, false, false},
		{"getreceivedbyaccount", getreceivedbyaccount, false, false},
		{"listreceivedbyaddress", listreceivedbyaddress, false, false},
		{"listreceivedbyaccount", listreceivedbyaccount, false, false},
		{"backupwallet", backupwallet, true, false},
		{"keypoolrefill", keypoolrefill, true, false},
		{"walletpassphrase", walletpassphrase, true, false},
		{"walletpassphrasechange", walletpassphrasechange, false, false},
		{"walletlock", walletlock, true, false},
		{"encryptwallet", encryptwallet, false, false},
		{"validateaddress", validateaddress, true, false},
		{"getbalance", getbalance, false, false},
		{"move", movecmd, false, false},
		{"sendfrom", sendfrom, false, false},
		{"sendmany", sendmany, false, false},
		{"addmultisigaddress", addmultisigaddress, false, false},
		{"createmultisig", createmultisig, true, true},
		{"getrawmempool", getrawmempool, true, false},
		{"getblock", getblock, false, false},
		{"getblockhash", getblockhash, false, false},
		{"gettransaction", gettransaction, false, false},
		{"listtransactions", listtransactions, false, false},
		{"listaddressgroupings", listaddressgroupings, false, false},
		{"signmessage", signmessage, false, false},
		{"verifymessage", verifymessage, false, false},
		{"getwork", getwork, true, false},
		{"listaccounts", listaccounts, false, false},
		{"settxfee", settxfee, false, false},
		{"getblocktemplate", getblocktemplate, true, false},
		{"submitblock", submitblock, false, false},
		{"listsinceblock", listsinceblock, false, false},
		{"dumpprivkey", dumpprivkey, true, false},
		{"dumpwallet", dumpwallet, true, false},
		{"importprivkey", importprivkey, false, false},
		{"importwallet", importwallet, false, false},
		{"listunspent", listunspent, false, false},
		{"getrawtransaction", getrawtransaction, false, false},
		{"createrawtransaction", createrawtransaction, false, false},
		{"decoderawtransaction", decoderawtransaction, false, false},
		{"signrawtransaction", signrawtransaction, false, false},
		{"sendrawtransaction", sendrawtransaction, false, false},
		{"gettxoutsetinfo", gettxoutsetinfo, true, false},
		{"gettxout", gettxout, true, false},
		{"lockunspent", lockunspent, false, false},
		{"listlockunspent", listlockunspent, false, false},
		{"verifychain", verifychain, true, false},
		{"makekeypair", makekeypair, false, false},
		{"sendalert", sendalert, true, false},
	}
)

func help(help bool, params ...string) interface{} {
	return nil
}
func stop(help bool, params ...string) interface{} {
	return nil
}

type SSLIOStreamDevice struct {
	needHandshake, useSSL bool
	// SSL socket here
}
type AcceptedConnection struct {
	// Device Protocol here
	// SSL socket here
}
type JSONRequest struct {
	ID     interface{}
	Method string
	params []string
}
type RPCTable struct {
	Commands []RPCCommand
}

func sendalert(help bool, params ...string) interface{} {
	return nil
}
