package cmds

import (
	"context"
	"flag"
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/version"

	"gitlab.com/parallelcoin/duo/pkg/subcmd"
)

func result(cmd string, f *flag.FlagSet, q bool, usage string) subcmd.ExitStatus {
	r, err := Cmd[cmd]("cli", f.Args(), nil)
	fmt.Print(r)
	if err == nil {
		return subcmd.ExitSuccess
	}
	fmt.Println(usage)
	if q {
		return subcmd.ExitQuit
	}
	return subcmd.ExitFailure
}

func init() {
	subcmd.Register(subcmd.HelpCommand(), "")
	subcmd.Register(subcmd.FlagsCommand(), "")
	subcmd.Register(subcmd.CommandsCommand(), "")
	subcmd.Register(&addmultisigaddress{}, "wallet")
	subcmd.Register(&addnode{}, "network")
	subcmd.Register(&backupwallet{}, "wallet")
	subcmd.Register(&createmultisig{}, "wallet")
	subcmd.Register(&createrawtransaction{}, "wallet")
	subcmd.Register(&decoderawtransaction{}, "wallet")
	subcmd.Register(&dumpprivkey{}, "wallet")
	subcmd.Register(&dumpwallet{}, "wallet")
	subcmd.Register(&encryptwallet{}, "wallet")
	subcmd.Register(&getaccount{}, "wallet")
	subcmd.Register(&getaccountaddress{}, "wallet")
	subcmd.Register(&getaddednodeinfo{}, "network")
	subcmd.Register(&getaddressesbyaccount{}, "wallet")
	subcmd.Register(&getbalance{}, "wallet")
	subcmd.Register(&getbestblockhash{}, "chain")
	subcmd.Register(&getblock{}, "chain")
	subcmd.Register(&getblockcount{}, "chain")
	subcmd.Register(&getblockhash{}, "chain")
	subcmd.Register(&getblocktemplate{}, "chain")
	subcmd.Register(&getconnectioncount{}, "network")
	subcmd.Register(&getdifficulty{}, "chain")
	subcmd.Register(&getgenerate{}, "client")
	subcmd.Register(&gethashespersec{}, "client")
	subcmd.Register(&getinfo{}, "client")
	subcmd.Register(&getmininginfo{}, "client")
	subcmd.Register(&getnetworkhashps{}, "chain")
	subcmd.Register(&getnewaddress{}, "wallet")
	subcmd.Register(&getpeerinfo{}, "network")
	subcmd.Register(&getrawmempool{}, "client")
	subcmd.Register(&getrawtransaction{}, "chain")
	subcmd.Register(&getreceivedbyaccount{}, "chain")
	subcmd.Register(&getreceivedbyaddress{}, "chain")
	subcmd.Register(&gettransaction{}, "chain")
	subcmd.Register(&gettxout{}, "chain")
	subcmd.Register(&gettxoutsetinfo{}, "client")
	subcmd.Register(&getwork{}, "chain")
	subcmd.Register(&help{}, "client")
	subcmd.Register(&importprivkey{}, "wallet")
	subcmd.Register(&importwallet{}, "wallet")
	subcmd.Register(&keypoolrefill{}, "wallet")
	subcmd.Register(&listaccounts{}, "wallet")
	subcmd.Register(&listaddressgroupings{}, "wallet")
	subcmd.Register(&listlockunspent{}, "wallet")
	subcmd.Register(&listreceivedbyaccount{}, "wallet")
	subcmd.Register(&listreceivedbyaddress{}, "wallet")
	subcmd.Register(&listsinceblock{}, "chain")
	subcmd.Register(&listtransactions{}, "wallet")
	subcmd.Register(&listunspent{}, "wallet")
	subcmd.Register(&lockunspent{}, "wallet")
	subcmd.Register(&makekeypair{}, "wallet")
	subcmd.Register(&move{}, "wallet")
	subcmd.Register(&sendalert{}, "network")
	subcmd.Register(&sendfrom{}, "wallet")
	subcmd.Register(&sendmany{}, "wallet")
	subcmd.Register(&sendrawtransaction{}, "wallet")
	subcmd.Register(&sendtoaddress{}, "wallet")
	subcmd.Register(&setaccount{}, "wallet")
	subcmd.Register(&setgenerate{}, "client")
	subcmd.Register(&settxfee{}, "wallet")
	subcmd.Register(&signmessage{}, "wallet")
	subcmd.Register(&signrawtransaction{}, "wallet")
	subcmd.Register(&stop{}, "client")
	subcmd.Register(&submitblock{}, "chain")
	subcmd.Register(&validateaddress{}, "chain")
	subcmd.Register(&verifychain{}, "client")
	subcmd.Register(&verifymessage{}, "client")
}

type addmultisigaddress struct {
	Info string
}

func (a *addmultisigaddress) Name() string {
	return "addmultisigaddress"
}
func (a *addmultisigaddress) Synopsis() string {
	return `<nrequired> <'["key","key"]'> [account]`
}
func (a *addmultisigaddress) Usage() string {
	return `
  Add a nrequired-to-sign multisignature address to the wallet each key is a Parallelcoin address or hex-encoded public key.
  If [account] is specified, assign address to [account].
  * Single quotes are required around the list, it is JSON array syntax.
`
}
func (a *addmultisigaddress) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *addmultisigaddress) Execute(c context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("multisigaddress", f, true, a.Info)
}

type addnode struct {
	Info string
}

func (a *addnode) Name() string {
	return "addnode"
}
func (a *addnode) Synopsis() string {
	return `<node> <add|remove|onetry>`
}
func (a *addnode) Usage() string {
	return `
  Attempts add or remove <node> from the addnode list or try a connection to <node> once.
`
}
func (a *addnode) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *addnode) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("addnode", f, true, a.Info)
}

type backupwallet struct {
	Info string
}

func (a *backupwallet) Name() string {
	return "backupwallet"
}
func (a *backupwallet) Synopsis() string {
	return `<destination>`
}
func (a *backupwallet) Usage() string {
	return `
  Safely copies wallet.dat to destination, which can be a directory or a path with filename.
`
}
func (a *backupwallet) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *backupwallet) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("backupwallet", f, true, a.Info)
}

type createmultisig struct {
	Info string
}

func (a *createmultisig) Name() string {
	return "createmultisig"
}
func (a *createmultisig) Synopsis() string {
	return `<nrequired> <'["key","key"]'>`
}
func (a *createmultisig) Usage() string {
	return `
  Creates a multi-signature address and returns a json object with keys:
         address : parallelcoin address
    redeemScript : hex-encoded redemption script
`
}
func (a *createmultisig) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *createmultisig) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("createmultisig", f, true, a.Info)
}

type createrawtransaction struct {
	Info string
}

func (a *createrawtransaction) Name() string {
	return "createrawtransaction"
}
func (a *createrawtransaction) Synopsis() string {
	return `<'[{"txid":txid,"vout":n},...]'> <'{address:amount,...}'>`
}
func (a *createrawtransaction) Usage() string {
	return `
  Create a transaction spending given inputs (array of objects containing transaction id and output number), sending to given address(es).
  Returns hex-encoded raw transaction.
  Note that the transaction's inputs are not signed, and it is not stored in the wallet or transmitted to the network.
`
}
func (a *createrawtransaction) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *createrawtransaction) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("createrawtransaction", f, true, a.Info)
}

type decoderawtransaction struct {
	Info string
}

func (a *decoderawtransaction) Name() string {
	return "decoderawtransaction"
}
func (a *decoderawtransaction) Synopsis() string {
	return `<hex string>`
}
func (a *decoderawtransaction) Usage() string {
	return `
  Return a JSON object representing the serialized, hex-encoded transaction.
`
}
func (a *decoderawtransaction) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *decoderawtransaction) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("decoderawtransaction", f, true, a.Info)
}

type dumpprivkey struct {
	Info string
}

func (a *dumpprivkey) Name() string {
	return "dumpprivkey"
}
func (a *dumpprivkey) Synopsis() string {
	return `<parallelcoinaddress>`
}
func (a *dumpprivkey) Usage() string {
	return `
  Reveals the private key corresponding to <parallelcoinaddress>.
`
}
func (a *dumpprivkey) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *dumpprivkey) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("dumpprivkey", f, true, a.Info)
}

type dumpwallet struct {
	Info string
}

func (a *dumpwallet) Name() string {
	return "dumpwallet"
}
func (a *dumpwallet) Synopsis() string {
	return `<filename>`
}
func (a *dumpwallet) Usage() string {
	return `
  Dumps all wallet keys in a human-readable format.
`
}
func (a *dumpwallet) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *dumpwallet) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("dumpwallet", f, true, a.Info)
}

type encryptwallet struct {
	Info string
}

func (a *encryptwallet) Name() string {
	return "encryptwallet"
}
func (a *encryptwallet) Synopsis() string {
	return `<passphrase>`
}
func (a *encryptwallet) Usage() string {
	return `
  Encrypts the wallet with <passphrase>. 
  (use single quotes if it contains spaces)
`
}
func (a *encryptwallet) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *encryptwallet) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("encryptwallet", f, true, a.Info)
}

type getaccount struct {
	Info string
}

func (a *getaccount) Name() string {
	return "getaccount"
}
func (a *getaccount) Synopsis() string {
	return `<parallelcoinaddress>`
}
func (a *getaccount) Usage() string {
	return `
  Returns the account associated with the given address.
`
}
func (a *getaccount) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getaccount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getaccount", f, true, a.Info)
}

type getaccountaddress struct {
	Info string
}

func (a *getaccountaddress) Name() string {
	return "getaccountaddress"
}
func (a *getaccountaddress) Synopsis() string {
	return `<account>`
}
func (a *getaccountaddress) Usage() string {
	return `
  Returns the current Parallelcoin address for receiving payments to this account.
`
}
func (a *getaccountaddress) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getaccountaddress) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getaccountaddress", f, true, a.Info)
}

type getaddednodeinfo struct {
	Info string
}

func (a *getaddednodeinfo) Name() string {
	return "getaddednodeinfo"
}
func (a *getaddednodeinfo) Synopsis() string {
	return `<dns> [node]`
}
func (a *getaddednodeinfo) Usage() string {
	return `
  Returns information about the given added node, or all added nodes (note that onetry addnodes are not listed here).
  If <dns> is false, only a list of added nodes will be provided, otherwise connected information will also be available.
`
}
func (a *getaddednodeinfo) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getaddednodeinfo) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getaddednodeinfo", f, true, a.Info)
}

type getaddressesbyaccount struct {
	Info string
}

func (a *getaddressesbyaccount) Name() string {
	return "getaddressesbyaccount"
}
func (a *getaddressesbyaccount) Synopsis() string {
	return `<account>`
}
func (a *getaddressesbyaccount) Usage() string {
	return `
  Returns the list of addresses for the given account.
`
}
func (a *getaddressesbyaccount) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getaddressesbyaccount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getaddressesbyaccount", f, true, a.Info)
}

type getbalance struct {
	Info string
}

func (a *getbalance) Name() string {
	return "getbalance"
}
func (a *getbalance) Synopsis() string {
	return `[account] [minconf=1]`
}
func (a *getbalance) Usage() string {
	return `
  If [account] is not specified, returns the server's total available balance.
  If [account] is specified, returns the balance in the account.
  [minconf=x] sets the number of confirmations required for an transaction to be counted in the total.
`
}
func (a *getbalance) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getbalance) Execute(_ context.Context, f *flag.FlagSet, i ...interface{}) subcmd.ExitStatus {
	return result("getbalance", f, true, a.Info)
}

type getbestblockhash struct {
	Info string
}

func (a *getbestblockhash) Name() string {
	return "getbestblockhash"
}
func (a *getbestblockhash) Synopsis() string {
	return ``
}
func (a *getbestblockhash) Usage() string {
	return `
  Returns the hash of the best (tip) block in the longest block chain (the head block).
`
}
func (a *getbestblockhash) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getbestblockhash) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getbestblockhash", f, true, a.Info)
}

type getblock struct {
	Info string
}

func (a *getblock) Name() string {
	return "getblock"
}
func (a *getblock) Synopsis() string {
	return `<hash> [verbose]`
}
func (a *getblock) Usage() string {
	return `
  If verbose is not present, returns a string that is serialized, hex-encoded data for block <hash>.
  If verbose is present, returns an Object with information about block <hash>.
`
}
func (a *getblock) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getblock) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getblock", f, true, a.Info)
}

type getblockcount struct {
	Info string
}

func (a *getblockcount) Name() string {
	return "getblockcount"
}
func (a *getblockcount) Synopsis() string {
	return ``
}
func (a *getblockcount) Usage() string {
	return `
  Returns the number of blocks in the longest block chain.
`
}
func (a *getblockcount) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getblockcount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getblockcount", f, true, a.Info)
}

type getblockhash struct {
	Info string
}

func (a *getblockhash) Name() string {
	return "getblockhash"
}
func (a *getblockhash) Synopsis() string {
	return `<index>`
}
func (a *getblockhash) Usage() string {
	return `
  Returns hash of block in best-block-chain at <index>.
`
}
func (a *getblockhash) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getblockhash) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getblockhash", f, true, a.Info)
}

type getblocktemplate struct {
	Info string
}

func (a *getblocktemplate) Name() string {
	return "getblocktemplate"
}
func (a *getblocktemplate) Synopsis() string {
	return `[params]`
}
func (a *getblocktemplate) Usage() string {
	return `
  Returns data needed to construct a block to work on:
              "version" : block version
    "previousblockhash" : hash of current highest block
         "transactions" : contents of non-coinbase transactions that should be included in the next block
          "coinbaseaux" : data that should be included in coinbase
        "coinbasevalue" : maximum allowable input to coinbase transaction, including the generation award and transaction fees
               "target" : hash target
              "mintime" : minimum timestamp appropriate for next block
              "curtime" : current timestamp
              "mutable" : list of ways the block template may be changed
           "noncerange" : range of valid nonces
           "sigoplimit" : limit of sigops in blocks
            "sizelimit" : limit of block size
                 "bits" : compressed target of next block
               "height" : height of the next block
  See https://en.bitcoin.it/wiki/BIP_0022 for full specification.
`
}
func (a *getblocktemplate) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getblocktemplate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getblocktemplate", f, true, a.Info)
}

type getconnectioncount struct {
	Info string
}

func (a *getconnectioncount) Name() string {
	return "getconnectioncount"
}
func (a *getconnectioncount) Synopsis() string {
	return ``
}
func (a *getconnectioncount) Usage() string {
	return `
  Returns the number of connections to other nodes.
`
}
func (a *getconnectioncount) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getconnectioncount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getconnectioncount", f, true, a.Info)
}

type getdifficulty struct {
	Info string
}

func (a *getdifficulty) Name() string {
	return "getdifficulty"
}
func (a *getdifficulty) Synopsis() string {
	return ``
}
func (a *getdifficulty) Usage() string {
	return `
  Returns the proof-of-work difficulty as a multiple of the minimum difficulty.
`
}
func (a *getdifficulty) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getdifficulty) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getdifficulty", f, true, a.Info)
}

type getgenerate struct {
	Info string
}

func (a *getgenerate) Name() string {
	return "getgenerate"
}
func (a *getgenerate) Synopsis() string {
	return ``
}
func (a *getgenerate) Usage() string {
	return `
  Returns true or false whether this node is mining.
`
}
func (a *getgenerate) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getgenerate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getgenerate", f, true, a.Info)
}

type gethashespersec struct {
	Info string
}

func (a *gethashespersec) Name() string {
	return "gethashespersec"
}
func (a *gethashespersec) Synopsis() string {
	return ``
}
func (a *gethashespersec) Usage() string {
	return `
  Returns a recent hashes per second performance measurement while generating.
`
}
func (a *gethashespersec) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *gethashespersec) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("gethashespersec", f, true, a.Info)
}

type getinfo struct {
	Info string
}

func (a *getinfo) Name() string {
	return "getinfo"
}
func (a *getinfo) Synopsis() string {
	return ``
}
func (a *getinfo) Usage() string {
	return `
  Returns an object containing various state info.
`
}
func (a *getinfo) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getinfo) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getinfo", f, true, a.Info)
}

type getmininginfo struct {
	Info string
}

func (a *getmininginfo) Name() string {
	return "getmininginfo"
}
func (a *getmininginfo) Synopsis() string {
	return ``
}
func (a *getmininginfo) Usage() string {
	return `
  Returns an object containing mining-related information.
`
}
func (a *getmininginfo) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getmininginfo) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getmininginfo", f, true, a.Info)
}

type getnetworkhashps struct {
	Info string
}

func (a *getnetworkhashps) Name() string {
	return "getnetworkhashps"
}
func (a *getnetworkhashps) Synopsis() string {
	return `[blocks=120] [height=-1]`
}
func (a *getnetworkhashps) Usage() string {
	return `
  Returns the estimated network hashes per second based on the last 120 blocks.
    Pass in [blocks] to override # of blocks.
    Pass in [height] to estimate the network speed at the time when a certain block was found, the default is -1 meaning best new block.
`
}
func (a *getnetworkhashps) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getnetworkhashps) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getnetworkhashps", f, true, a.Info)
}

type getnewaddress struct {
	Info string
}

func (a *getnewaddress) Name() string {
	return "getnewaddress"
}
func (a *getnewaddress) Synopsis() string {
	return `[account]`
}
func (a *getnewaddress) Usage() string {
	return `
  Returns a new Parallelcoin address for receiving payments.
  If [account] is specified (recommended), it is added to the address book so payments received with the address will be credited to [account].
`
}
func (a *getnewaddress) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getnewaddress) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getnewaddress", f, true, a.Info)
}

type getpeerinfo struct {
	Info string
}

func (a *getpeerinfo) Name() string {
	return "getpeerinfo"
}
func (a *getpeerinfo) Synopsis() string {
	return ``
}
func (a *getpeerinfo) Usage() string {
	return `
  Returns data about each connected network node.
`
}
func (a *getpeerinfo) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getpeerinfo) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getpeerinfo", f, true, a.Info)
}

type getrawmempool struct {
	Info string
}

func (a *getrawmempool) Name() string {
	return "getrawmempool"
}
func (a *getrawmempool) Synopsis() string {
	return ``
}
func (a *getrawmempool) Usage() string {
	return `
  Returns all transaction ids in memory pool.
`
}
func (a *getrawmempool) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getrawmempool) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getrawmempool", f, true, a.Info)
}

type getrawtransaction struct {
	Info string
}

func (a *getrawtransaction) Name() string {
	return "getrawtransaction"
}
func (a *getrawtransaction) Synopsis() string {
	return `<txid> [verbose]`
}
func (a *getrawtransaction) Usage() string {
	return `
  If verbose is not set, returns a string that is serialized, hex-encoded data for <txid>.
  If verbose is set, returns an Object with information about <txid>.
`
}
func (a *getrawtransaction) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getrawtransaction) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getrawtransaction", f, true, a.Info)
}

type getreceivedbyaccount struct {
	Info string
}

func (a *getreceivedbyaccount) Name() string {
	return "getreceivedbyaccount"
}
func (a *getreceivedbyaccount) Synopsis() string {
	return `<account> [minconf=1]`
}
func (a *getreceivedbyaccount) Usage() string {
	return `
  Returns the total amount received by addresses with <account> in transactions with at least [minconf] confirmations.
`
}
func (a *getreceivedbyaccount) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getreceivedbyaccount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("", f, true, a.Info)
}

type getreceivedbyaddress struct {
	Info string
}

func (a *getreceivedbyaddress) Name() string {
	return "getreceivedbyaddress"
}
func (a *getreceivedbyaddress) Synopsis() string {
	return `<parallelcoinaddress> [minconf=1]`
}
func (a *getreceivedbyaddress) Usage() string {
	return `
  Returns the total amount received by <parallelcoinaddress> in transactions with at least [minconf] confirmations.
`
}
func (a *getreceivedbyaddress) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getreceivedbyaddress) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getreceivedbyaddress", f, true, a.Info)
}

type gettransaction struct {
	Info string
}

func (a *gettransaction) Name() string {
	return "gettransaction"
}
func (a *gettransaction) Synopsis() string {
	return `<txid>`
}
func (a *gettransaction) Usage() string {
	return `
  Get detailed information about in-wallet transaction <txid>.
`
}
func (a *gettransaction) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *gettransaction) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("gettransaction", f, true, a.Info)
}

type gettxout struct {
	Info string
}

func (a *gettxout) Name() string {
	return "gettxout"
}
func (a *gettxout) Synopsis() string {
	return `<txid> <n> [includemempool]`
}
func (a *gettxout) Usage() string {
	return `
  Returns details about an unspent transaction output.
`
}
func (a *gettxout) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *gettxout) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("gettxout", f, true, a.Info)
}

type gettxoutsetinfo struct {
	Info string
}

func (a *gettxoutsetinfo) Name() string {
	return "gettxoutsetinfo"
}
func (a *gettxoutsetinfo) Synopsis() string {
	return ``
}
func (a *gettxoutsetinfo) Usage() string {
	return `
  Returns statistics about the unspent transaction output set.
`
}
func (a *gettxoutsetinfo) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *gettxoutsetinfo) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("gettxoutsetinfo", f, true, a.Info)
}

type getwork struct {
	Info string
}

func (a *getwork) Name() string {
	return "getwork"
}
func (a *getwork) Synopsis() string {
	return `[data]`
}
func (a *getwork) Usage() string {
	return `
  If [data] is not specified, returns formatted hash data to work on:
    "midstate" : precomputed hash state after hashing the first half of the data (DEPRECATED)
        "data" : block data
       "hash1" : formatted hash buffer for second hash (DEPRECATED)
      "target" : little endian hash target
  If [data] is specified, tries to solve the block and returns true if it was successful.
`
}
func (a *getwork) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *getwork) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("getwork", f, true, a.Info)
}

type help struct {
	Info string
}

func (a *help) Name() string {
	return "help"
}
func (a *help) Synopsis() string {
	return `[command]`
}
func (a *help) Usage() string {
	return `
  List commands, or get help for a command.
`
}
func (a *help) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *help) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	version.Print()
	return result("help", f, true, a.Info)
}

type importprivkey struct {
	Info string
}

func (a *importprivkey) Name() string {
	return "importprivkey"
}
func (a *importprivkey) Synopsis() string {
	return `<parallelcoinprivkey> [label] [rescan]`
}
func (a *importprivkey) Usage() string {
	return `
  Adds a private key (as returned by dumpprivkey) to your wallet.
`
}
func (a *importprivkey) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *importprivkey) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("importprivkey", f, true, a.Info)
}

type importwallet struct {
	Info string
}

func (a *importwallet) Name() string {
	return "importwallet"
}
func (a *importwallet) Synopsis() string {
	return `<filename>`
}
func (a *importwallet) Usage() string {
	return `
  Imports keys from a wallet dump file (see dumpwallet).
`
}
func (a *importwallet) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *importwallet) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("importwallet", f, true, a.Info)
}

type keypoolrefill struct {
	Info string
}

func (a *keypoolrefill) Name() string {
	return "keypoolrefill"
}
func (a *keypoolrefill) Synopsis() string {
	return ``
}
func (a *keypoolrefill) Usage() string {
	return `
  Fills the keypool.
`
}
func (a *keypoolrefill) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *keypoolrefill) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("keypoolrefill", f, true, a.Info)
}

type listaccounts struct {
	Info string
}

func (a *listaccounts) Name() string {
	return "listaccounts"
}
func (a *listaccounts) Synopsis() string {
	return `[minconf=1]`
}
func (a *listaccounts) Usage() string {
	return `
  Returns Object that has account names as keys, account balances as values.
`
}
func (a *listaccounts) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listaccounts) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listaccounts", f, true, a.Info)
}

type listaddressgroupings struct {
	Info string
}

func (a *listaddressgroupings) Name() string {
	return "listaddressgroupings"
}
func (a *listaddressgroupings) Synopsis() string {
	return ``
}
func (a *listaddressgroupings) Usage() string {
	return `
  Lists groups of addresses which have had their common ownership made public by common use as inputs or as the resulting change in past transactions.
`
}
func (a *listaddressgroupings) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listaddressgroupings) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listaddressgroupings", f, true, a.Info)
}

type listlockunspent struct {
	Info string
}

func (a *listlockunspent) Name() string {
	return "listlockunspent"
}
func (a *listlockunspent) Synopsis() string {
	return ``
}
func (a *listlockunspent) Usage() string {
	return `
  Returns list of temporarily unspendable outputs.
`
}
func (a *listlockunspent) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listlockunspent) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listlockunspent", f, true, a.Info)
}

type listreceivedbyaccount struct {
	Info string
}

func (a *listreceivedbyaccount) Name() string {
	return "listreceivedbyaccount"
}
func (a *listreceivedbyaccount) Synopsis() string {
	return `[minconf=1] [includeempty]`
}
func (a *listreceivedbyaccount) Usage() string {
	return `
  [minconf] is the minimum number of confirmations before payments are included.
  [includeempty] whether to include accounts that haven't received any payments.
  Returns an array of objects containing:
          "account" : the account of the receiving addresses
           "amount" : total amount received by addresses with this account
  "confirmations" : number of confirmations of the most recent transaction included
`
}
func (a *listreceivedbyaccount) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listreceivedbyaccount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listreceivedbyaccount", f, true, a.Info)
}

type listreceivedbyaddress struct {
	Info string
}

func (a *listreceivedbyaddress) Name() string {
	return "listreceivedbyaddress"
}
func (a *listreceivedbyaddress) Synopsis() string {
	return `[minconf=1] [includeempty]`
}
func (a *listreceivedbyaddress) Usage() string {
	return `
  [minconf] is the minimum number of confirmations before payments are included.
  [includeempty] whether to include addresses that haven't received any payments.
  Returns an array of objects containing:
          "address" : receiving address
          "account" : the account of the receiving address
           "amount" : total amount received by the address
    "confirmations" : number of confirmations of the most recent transaction included
            "txids" : list of transactions with outputs to the address
`
}
func (a *listreceivedbyaddress) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listreceivedbyaddress) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listreceivedbyaddress", f, true, a.Info)
}

type listsinceblock struct {
	Info string
}

func (a *listsinceblock) Name() string {
	return "listsinceblock"
}
func (a *listsinceblock) Synopsis() string {
	return `[blockhash] [target-confirmations]`
}
func (a *listsinceblock) Usage() string {
	return `
  Get all transactions in blocks since block [blockhash], or all transactions if omitted.
`
}
func (a *listsinceblock) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listsinceblock) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listsinceblock", f, true, a.Info)
}

type listtransactions struct {
	Info string
}

func (a *listtransactions) Name() string {
	return "listtransactions"
}
func (a *listtransactions) Synopsis() string {
	return `[account] [count=10] [from=0]`
}
func (a *listtransactions) Usage() string {
	return `
  Returns up to [count] most recent transactions skipping the first [from] transactions for account [account].
`
}
func (a *listtransactions) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listtransactions) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listtransactions", f, true, a.Info)
}

type listunspent struct {
	Info string
}

func (a *listunspent) Name() string {
	return "listunspent"
}
func (a *listunspent) Synopsis() string {
	return `[minconf=1] [maxconf=9999999] ['["address",...]']`
}
func (a *listunspent) Usage() string {
	return `
  Returns array of unspent transaction outputs with between minconf and maxconf (inclusive) confirmations.
  Optionally filtered to only include txouts paid to specified addresses.
  Results are an array of Objects, each of which has: {txid, vout, scriptPubKey, amount, confirmations}
`
}
func (a *listunspent) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *listunspent) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("listunspent", f, true, a.Info)
}

type lockunspent struct {
	Info string
}

func (a *lockunspent) Name() string {
	return "lockunspent"
}
func (a *lockunspent) Synopsis() string {
	return `[unlock] <[array-of-Objects]>`
}
func (a *lockunspent) Usage() string {
	return `
  Updates list of temporarily unspendable outputs.
`
}
func (a *lockunspent) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *lockunspent) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("lockunspent", f, true, a.Info)
}

type makekeypair struct {
	Info string
}

func (a *makekeypair) Name() string {
	return "makekeypair"
}
func (a *makekeypair) Synopsis() string {
	return `[prefix]`
}
func (a *makekeypair) Usage() string {
	return `
  Make a public/private key pair.
    [prefix] is optional preferred prefix for the public key.
`
}
func (a *makekeypair) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *makekeypair) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("makekeypair", f, true, a.Info)
}

type move struct {
	Info string
}

func (a *move) Name() string {
	return "move"
}
func (a *move) Synopsis() string {
	return `<fromaccount> <toaccount> <amount> [minconf=1] ['comment']`
}
func (a *move) Usage() string {
	return `
  Move from one account in your wallet to another.
`
}
func (a *move) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *move) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("move", f, true, a.Info)
}

type sendalert struct {
	Info string
}

func (a *sendalert) Name() string {
	return "sendalert"
}
func (a *sendalert) Synopsis() string {
	return `<message> <privatekey> <minver> <maxver> <priority> <id> [cancelupto]`
}
func (a *sendalert) Usage() string {
	return `
  <message> is the alert text message
  <privatekey> is base58 hex string of alert master private key
  <minver> is the minimum applicable internal client version
  <maxver> is the maximum applicable internal client version
  <priority> is integer priority number
  <id> is the alert id
  [cancelupto] cancels all alert id's up to this number
  Returns true or false.
`
}
func (a *sendalert) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *sendalert) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("sendalert", f, true, a.Info)
}

type sendfrom struct {
	Info string
}

func (a *sendfrom) Name() string {
	return "sendfrom"
}
func (a *sendfrom) Synopsis() string {
	return `<fromaccount> <toparallelcoinaddress> <amount> [minconf=1] [comment=""] [comment-to=""]`
}
func (a *sendfrom) Usage() string {
	return `
  <amount> is a real and is rounded to the nearest 0.00000001.
`
}
func (a *sendfrom) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *sendfrom) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("sendfrom", f, true, a.Info)
}

type sendmany struct {
	Info string
}

func (a *sendmany) Name() string {
	return "sendmany"
}
func (a *sendmany) Synopsis() string {
	return `<fromaccount> <'{"address":amount,...}'> [minconf=1] [comment]`
}
func (a *sendmany) Usage() string {
	return `
  amounts are double-precision floating point numbers.
`
}
func (a *sendmany) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *sendmany) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("sendmany", f, true, a.Info)
}

type sendrawtransaction struct {
	Info string
}

func (a *sendrawtransaction) Name() string {
	return "sendrawtransaction"
}
func (a *sendrawtransaction) Synopsis() string {
	return `<hex string>`
}
func (a *sendrawtransaction) Usage() string {
	return `
  Submits raw transaction (serialized, hex-encoded) to local node and network.
`
}
func (a *sendrawtransaction) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *sendrawtransaction) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("sendrawtransaction", f, true, a.Info)
}

type sendtoaddress struct {
	Info string
}

func (a *sendtoaddress) Name() string {
	return "sendtoaddress"
}
func (a *sendtoaddress) Synopsis() string {
	return `<parallelcoinaddress> <amount> [comment=""] [comment-to=""]`
}
func (a *sendtoaddress) Usage() string {
	return `
  <amount> is a real and is rounded to the nearest 0.00000001.
`
}
func (a *sendtoaddress) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *sendtoaddress) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("sendtoaddress", f, true, a.Info)
}

type setaccount struct {
	Info string
}

func (a *setaccount) Name() string {
	return "setaccount"
}
func (a *setaccount) Synopsis() string {
	return `<parallelcoinaddress> <account>`
}
func (a *setaccount) Usage() string {
	return `
  Sets the account associated with the given address.
`
}
func (a *setaccount) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *setaccount) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("setaccount", f, true, a.Info)
}

type setgenerate struct {
	Info string
}

func (a *setgenerate) Name() string {
	return "setgenerate"
}
func (a *setgenerate) Synopsis() string {
	return `<generate> [genproclimit]`
}
func (a *setgenerate) Usage() string {
	return `
  <generate> is true or false to turn generation on or off.
  Generation is limited to [genproclimit] processors, -1 is unlimited.
`
}
func (a *setgenerate) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *setgenerate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("setgenerate", f, true, a.Info)
}

type settxfee struct {
	Info string
}

func (a *settxfee) Name() string {
	return "settxfee"
}
func (a *settxfee) Synopsis() string {
	return `<amount>`
}
func (a *settxfee) Usage() string {
	return `
  <amount> is a real and is rounded to the nearest 0.00000001.
`
}
func (a *settxfee) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *settxfee) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("settxfee", f, true, a.Info)
}

type signmessage struct {
	Info string
}

func (a *signmessage) Name() string {
	return "signmessage"
}
func (a *signmessage) Synopsis() string {
	return `<parallelcoinaddress> <'message'>`
}
func (a *signmessage) Usage() string {
	return `
  Sign a message with the private key of an address.
`
}
func (a *signmessage) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *signmessage) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("signmessage", f, true, a.Info)
}

type signrawtransaction struct {
	Info string
}

func (a *signrawtransaction) Name() string {
	return "signrawtransaction"
}
func (a *signrawtransaction) Synopsis() string {
	return `<hex string> ['{"txid":txid,"vout":n,"scriptPubKey":hex,"redeemScript":hex},...'] [<privatekey1>,...] [sighashtype="ALL"]`
}
func (a *signrawtransaction) Usage() string {
	return `
  Sign inputs for raw transaction (serialized, hex-encoded).
  Second optional argument (may be null) is an array of previous transaction outputs that this transaction depends on but may not yet be in the block chain.
  Third optional argument (may be null) is an array of base58-encoded private keys that, if given, will be the only keys used to sign the transaction.
  Fourth optional argument is a string that is one of six values; ALL, NONE, SINGLE or ALL|ANYONECANPAY, NONE|ANYONECANPAY, SINGLE|ANYONECANPAY.
  Returns json object with keys:
         hex : raw transaction with signature(s) (hex-encoded string)
    complete : 1 if transaction has a complete set of signature (0 if not).
`
}
func (a *signrawtransaction) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *signrawtransaction) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("signrawtransaction", f, true, a.Info)
}

type stop struct {
	Info string
}

func (a *stop) Name() string {
	return "stop"
}
func (a *stop) Synopsis() string {
	return ``
}
func (a *stop) Usage() string {
	return `
  Stop Parallelcoin server.
`
}
func (a *stop) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *stop) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("stop", f, true, a.Info)
}

type submitblock struct {
	Info string
}

func (a *submitblock) Name() string {
	return "submitblock"
}
func (a *submitblock) Synopsis() string {
	return `<hex data> [optional-params-obj]`
}
func (a *submitblock) Usage() string {
	return `
  [optional-params-obj] parameter is currently ignored.
  Attempts to submit new block to network.
  See https://en.bitcoin.it/wiki/BIP_0022 for full specification..
`
}
func (a *submitblock) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *submitblock) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("submitblock", f, true, a.Info)
}

type validateaddress struct {
	Info string
}

func (a *validateaddress) Name() string {
	return "validateaddress"
}
func (a *validateaddress) Synopsis() string {
	return `<parallelcoinaddress>`
}
func (a *validateaddress) Usage() string {
	return `
  Return information about <parallelcoinaddress>.
`
}
func (a *validateaddress) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *validateaddress) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("validateaddress", f, true, a.Info)
}

type verifychain struct {
	Info string
}

func (a *verifychain) Name() string {
	return "verifychain"
}
func (a *verifychain) Synopsis() string {
	return `[checklevel=3] [num blocks]`
}
func (a *verifychain) Usage() string {
	return `
  Verifies blockchain database. 
  [checklevel] is 1-4 default is 3
`
}
func (a *verifychain) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *verifychain) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("verifychain", f, true, a.Info)
}

type verifymessage struct {
	Info string
}

func (a *verifymessage) Name() string {
	return "verifymessage"
}
func (a *verifymessage) Synopsis() string {
	return `<parallelcoinaddress> <signature> <message>`
}
func (a *verifymessage) Usage() string {
	return `
  Verify a signed message.
`
}
func (a *verifymessage) SetFlags(f *flag.FlagSet) {
	a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage())
}
func (a *verifymessage) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
	return result("verifymessage", f, true, a.Info)
}

/*
  subcmd.Register(&COMMAND{}, "")
type COMMAND struct {
}
func (a *COMMAND) Name() string {
  return "COMMAND"
}
func (a *COMMAND) Synopsis() string {
  return ``
}
func (a *COMMAND) Usage() string {
  return `.`
}
func (a *COMMAND) SetFlags(f *flag.FlagSet) { a.Info = fmt.Sprint("\n\n" + a.Name() + " " + a.Synopsis() + "\n" + a.Usage()) }
func (a *COMMAND) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcmd.ExitStatus {
  return subcmd.ExitSuccess
}
*/
