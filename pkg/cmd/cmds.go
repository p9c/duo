package cmds

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anaskhan96/base58check"
	"gitlab.com/parallelcoin/duo/pkg/algos"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/net"
	"gitlab.com/parallelcoin/duo/pkg/rpc"
	"gitlab.com/parallelcoin/duo/pkg/server"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	"gitlab.com/parallelcoin/duo/pkg/server/state"
	"gitlab.com/parallelcoin/duo/pkg/subcmd"
	"gitlab.com/parallelcoin/duo/pkg/util"
	"gitlab.com/parallelcoin/duo/pkg/version"
	"gitlab.com/parallelcoin/duo/pkg/wallet"
	"os"
	"strconv"
	"strings"
)

// GetInfo is the list of fields returned by the getinfo function
type GetInfo struct {
	Version           int     `json:"version"`
	Protocolversion   int     `json:"protocolversion"`
	Walletversion     int     `json:"walletversion"`
	Balance           float64 `json:"balance"`
	Blocks            int     `json:"blocks"`
	Timeoffset        int64   `json:"timeoffset"`
	Connections       int     `json:"connections"`
	Proxy             string  `json:"proxy"`
	PowAlgoID         int     `json:"pow_algo_id"`
	PowAlgo           string  `json:"pow_algo"`
	Difficulty        float64 `json:"difficulty"`
	DifficultySha256D float64 `json:"difficulty_sha256d"`
	DifficultyScrypt  float64 `json:"difficulty_scrypt"`
	Testnet           bool    `json:"testnet"`
	Keypoololdest     int64   `json:"keypoololdest"`
	Keypoolsize       int     `json:"keypoolsize"`
	Paytxfee          int64   `json:"paytxfee"`
	UnlockedUntil     int64   `json:"unlocked_until"`
	Errors            string  `json:"errors"`
}

func isValidJSONArray(s string) (o []string, i bool) {
	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		return nil, false
	}
	return o, true
}
func isValidJSONMap(s string) (o map[string]interface{}, i bool) {
	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		return nil, false
	}
	return o, true
}

func isInteger(s string) (int, bool) {
	r, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return r, true
}

func isFloat(s string) (float64, bool) {
	r, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return r, true
}

func isEmpty(s string) bool {
	if s == "" {
		return true
	}
	return false
}

func failure(mode string, erri error, s ...interface{}) (r string, err error) {
	err = erri
	switch mode {
	case "cli", "":
		r += "\nERROR: "
		for i := range s {
			r += s[i].(string) + "  "
		}
		r += "\n"
	case "json":
		r += "{\"error\":\""
		for i := range s {
			r += s[i].(string) + "  "
		}
		r += "}"
	}
	if err == nil {
		err = errors.New(r)
	}
	return
}

// Cmd is the map of the API commands
var Cmd map[string]func(string, []string, error) (string, error)

func init() {
	Cmd = map[string]func(string, []string, error) (string, error){

		"addmultisigaddress": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				n, i := isInteger(f[0])
				if !i {
					return failure(mode, err, "First parameter was not a number")
				}
				if !json.Valid([]byte(f[1])) {
					return failure(mode, err, "List of addresses was not formatted correctly")
				}
				list, i := isValidJSONArray(f[1])
				if !i {
					return failure(mode, err, "JSON was valid but not an array")
				}
				r += "  Number of signers required:" + strconv.Itoa(n) + "\n  Addresses: "
				for i := range list {
					if i > 0 {
						r += ", "
					}
					r += list[i]
				}
				r += "\n"
				if f[2] != "" {
					r += "  Account address will be added to: '" + f[2] + "'"
				}
			}
			return
		},

		"addnode": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				node := f[0]
				if isEmpty(node) {
					return failure(mode, err, "<node> is a required parameter")
				}
				action := f[1]
				if isEmpty(action) {
					return failure(mode, err, "Must have one of <add|remove|onetry>")
				}
				if !(action == "add" || f[1] == "remove" || f[1] == "onetry") {
					return failure(mode, err, "unrecognised option ", action)
				}
				switch action {
				case "add":
					r += "Adding "
				case "remove":
					r += "Removing "
				case "onetry":
					r += "Trying once "
				}
				r += "'" + node + "'"
			}
			return
		},

		"backupwallet": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				wallet := f[0]
				if len(f[0]) != 1 {
					return failure(mode, err, "ERROR: <destination> is a required parameter")
				}
				r += "Backing up wallet to '" + wallet + "'\n"
			}
			return
		},

		"createmultisig": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				n, N := isInteger(f[0])
				list, L := isValidJSONArray(f[1])
				switch {
				case len(f) > 2:
					return failure(mode, err, "Excess parameters, wanted <nrequired> <'[\"key\",\"key\"]'>, got ", f)
				case isEmpty(f[0]):
					return failure(mode, err, "<nrequired> is a required parameter")
				case !N:
					return failure(mode, err, "First parameter was not a number")
				case !json.Valid([]byte(f[1])):
					return failure(mode, err, "List of addresses was not formatted correctly")
				case !L:
					return failure(mode, err, "List of addresses was not formatted correctly")
				}
				r += "Creating multisig address\n"
				r += "  Number of signers required: " + strconv.Itoa(n) + "\n"
				r += "  Addresses:\n    "
				for i := range list {
					if i > 0 {
						r += ", \n    "
					}
					r += list[i]
				}
				r += "\n"
			}
			return
		},

		"createrawtransaction": func(mode string, f []string, erri error) (r string, err error) {
			tx, j := isValidJSONMap(f[0])
			addrs, a := isValidJSONArray(f[1])
			err = erri
			switch mode {
			case "cli", "":
				switch {
				case len(f[0]) != 1:
					return failure(mode, err, "<transaction> is a required parameter")
				case !json.Valid([]byte(f[0])):
					return failure(mode, err, "Transaction JSON was not formatted correctly")
				case !j:
					return failure(mode, err, "Transaction was not correct JSON ", f[0])
				case isEmpty(f[1]):
					return failure(mode, err, "<addresses> is a required parameter")
				case !json.Valid([]byte(f[1])):
					return failure(mode, err, "List of addresses was not formatted correctly")
				case !a:
					return failure(mode, err, "List of addresses was not correct JSON")
				}
				r += "Transaction: " + fmt.Sprint(tx) + "\n"
				r += "Addresses: " + fmt.Sprint(addrs) + "\n"
			}
			return
		},

		"decoderawtransaction": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f[0]) != 1 {
					return failure(mode, err, "<hex string> is a required parameter")
				}
				r += "Decoding " + f[0]
			}
			return
		},

		"dumpprivkey": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if f[0] == "" {
					return failure(mode, err, "<parallelcoinaddress> is a required parameter")
				}
				r += "Dumping address '" + f[0] + "'\n"

			}
			return
		},

		"dumpwallet": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 2 {
					return failure(mode, err, "Excess arguments")
				}
				if wallet.Db.Filename == "" {
					wallet.Db.SetFilename(*args.DataDir + "/" + *args.Wallet)
				}
				logger.Debug(wallet.Db.Filename)
				if db, err := wallet.NewDB(); err != nil {
					return failure(mode, err, "unable to open wallet")
				} else if err := db.Open(); err != nil {
					return failure(mode, err, "")
				} else if dump, err := db.Dump(); err != nil {
					return failure(mode, err, "")
				} else {
					r += dump + "\n"
					if len(f) == 1 {
						if wallet.Db.Filename == "" {
							wallet.Db.SetFilename(*args.DataDir + "/" + *args.Wallet)
							dumpfile, err := os.OpenFile(f[0], os.O_RDWR|os.O_CREATE, 0600)
							if err != nil {
								return failure(mode, err, "")
							}
							defer dumpfile.Close()
							// Write r to dump file
						}
					}
				}
			}
			return
		},

		"encryptwallet": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f[0]) != 1 {
					return failure(mode, err, "<passphrase> is a required parameter")
				}
				r += "Encypting wallet"
			}
			return
		},

		"getaccount": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f[0]) != 1 {
					return failure(mode, err, "<parallelcoinaddress> is a required parameter")
				}
				r += "Getting account '" + f[0] + "'..."
			}
			return
		},

		"getaccountaddress": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f[0]) != 1 {
					return failure(mode, err, "<account> is a required parameter")
				}
				r += "Getting address from account '" + f[0] + "'...\n"
			}
			return
		},

		"getaddednodeinfo": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				switch {
				case !(f[0] == "true" || f[0] == "false"):
					return failure(mode, err, "<dns> must be set to true or false")
				case f[0] == "true":
					r += "Showing all DNS information\n"
				case f[0] == "false":
					r += "Showing no DNS information\n"
				case f[1] == "":
					r += "Showing info about all added peers\n"
				case f[1] != "":
					r += "Showing info about '" + f[1] + "'"
				}
			}
			return
		},

		"getaddressesbyaccount": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if isEmpty(f[0]) {
					return failure(mode, err, "<account> is a required parameter")
				}
				r += "Getting info about account '" + f[0] + "'\n"
			}
			return
		},

		"getbalance": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				minconf := 1
				account := ""
				if len(f) > 2 {
					return failure(mode, err, "Excess arguments: ", f)
				}
				for i := range f {
					if strings.HasPrefix(f[i], "minconf=") {
						split := strings.Split(f[i], "=")
						n, i := isInteger(split[1])
						if i {
							minconf = n
						} else {
							return failure(mode, err, "minconf value is not a number")
						}
					} else {
						account = f[i]
					}
				}
				r += "minconf =" + strconv.Itoa(minconf) + "\n"
				if account != "" {
					r += "account = '" + account + "'\n"
				}
			}
			return
		},

		"getbestblockhash": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				r += "Getting best block hash ...\n"
			}
			return
		},

		"getblock": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) == 0 {
					return failure(mode, err, "No arguments given for '", "'\n")
				}
				verbose := false
				hash := ""
				for i := range f {
					if f[i] == "verbose" {
						if verbose {
							return failure(mode, err, "verbose can only be set once")
						}
						verbose = true
					} else {
						if hash != "" {
							return failure(mode, err, "Too many arguments given for '", "'\n")
						}
						hash = f[i]
					}
				}
				if verbose {
					r += "verbose enabled\n"
				}
				if hash == "" {
					return failure(mode, err, "No hash given")
				}
				r += "Info about: '" + hash + "'\n"
			}
			return
		},

		"getblockcount": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required")
				}
				r += "Getting the number of blocks in the longest block chain ..."
			}
			return
		},

		"getblockhash": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 1 {
					return failure(mode, err, "No/too many arguments given for '", "'\n")
				}
				n, i := isInteger(f[0])
				if i {
					r += "Getting block hash of block number " + strconv.Itoa(n) + "\n"
				} else {
					return failure(mode, err, "Argument is not a number")
				}
			}
			return
		},

		"getblocktemplate": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					if !json.Valid([]byte(f[0])) {
						return failure(mode, err, "Transaction JSON was not formatted correctly")
					}
					m, i := isValidJSONMap(f[0])
					if !i {
						return failure(mode, err, "Parameters incorrectly formatted ", f[0])
					}
					params, _ := json.Marshal(m)
					r += "Parameters " + string(params) + "\n"
				}
			}
			return
		},

		"getconnectioncount": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				r += "Getting connection count ..."
			}
			return
		},

		"getdifficulty": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				r += "Getting current Proof of Work difficulty ..."
			}
			return
		},

		"getgenerate": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Getting mining status ...")
			}
			return
		},

		"gethashespersec": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Getting miner's current hashrate ...")
			}
			return
		},

		"getinfo": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				info := GetInfo{
					Version:           version.Client,
					Protocolversion:   version.Protocol,
					Walletversion:     server.Walletdb.Version(),
					Balance:           server.Walletdb.GetBalance(),
					Timeoffset:        util.GetTimeOffset(),
					Connections:       len(net.Nodes),
					Proxy:             *args.Proxy,
					PowAlgoID:         algos.Code(*args.Algo),
					PowAlgo:           *args.Algo,
					Difficulty:        rpc.GetDifficulty(block.ChainIndex, algos.Code(*args.Algo)),
					DifficultySha256D: rpc.GetDifficulty(block.ChainIndex, algos.SHA256D),
					DifficultyScrypt:  rpc.GetDifficulty(block.ChainIndex, algos.SCRYPT),
					Testnet:           *args.TestNet,
					Keypoololdest:     server.Walletdb.GetOldestKeyPoolTime(),
					Keypoolsize:       server.Walletdb.GetKeyPoolSize(),
					Paytxfee:          state.TransactionFee,
					UnlockedUntil:     wallet.Db.UnlockedUntil,
					Errors:            server.GetWarnings("statusbar"),
				}
				reply, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					return failure(mode, err, "Unable to marshal struct into JSON")
				}
				r += string(reply) + "\n"
			}
			return
		},

		"getmininginfo": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				r += "Getting info about mining on this server ...\n"
			}
			return
		},

		"getnetworkhashps": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				blocks := 120
				height := -1
				for i := range f {
					if strings.HasPrefix(f[i], "blocks=") {
						value := strings.Split(f[i], "=")[1]
						var i bool
						blocks, i = isInteger(value)
						if !i {
							return failure(mode, err, "Value after blocks= not a number")
						}
					}
					if strings.HasPrefix(f[i], "height=") {
						value := strings.Split(f[i], "=")[1]
						var i bool
						height, i = isInteger(value)
						if !i {
							return failure(mode, err, "Value after height= not a number")
						}
					}
				}
				r += "Getting estimated hashes per second with [blocks] average at [height] block \n"
				r += "blocks = " + strconv.Itoa(blocks) + " height = " + strconv.Itoa(height) + "\n"
			}
			return
		},

		"getnewaddress": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				switch {
				case len(f) > 1:
					return failure(mode, err, "Excess arguments after [account]")
				case f[0] != "":
					r += "Adding address to account " + f[0] + "\n"
				default:
					r += "Adding new address to wallet\n"
				}
			}
			return
		},

		"getpeerinfo": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				r += "Getting peer information ...\n"
			}
			return
		},

		"getrawmempool": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Getting raw mempool information ...")
			}
			return
		},

		"getrawtransaction": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				nargs := len(f)
				verbose := false
				txid := ""
				switch {
				case nargs < 1:
					return failure(mode, err, "A <txid> must be specified.")
				case nargs > 2:
					return failure(mode, err, "Too many arguments specified, want <txid> [verbose]")
				}
				for i := range f {
					if f[i] == "verbose" {
						if verbose {
							return failure(mode, err, "Already set verbose flag")
						}
						verbose = true
					} else if txid != "" {
						if f[i] != "" {
							return failure(mode, err, "Already set <txid>")
						}
						txid = f[i]
					} else {
						txid = f[i]
					}
				}
				r += "Returning txid '" + txid + "'\n"
				if verbose {
					r += "  in raw hex encoded tx\n"
				} else {
					r += "  in JSON encoded form\n"
				}
			}
			return
		},

		"getreceivedbyaccount": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				minconf := 1
				account := ""
				for i := range f {
					if strings.HasPrefix(f[i], "minconf=") {
						split := strings.Split(f[i], "=")
						n, i := isInteger(split[1])
						if i {
							minconf = n
						} else {
							return failure(mode, err, "minconf value is not a number")
						}
					} else {
						account = f[i]
					}
				}
				r += "Getting balance with minconf=" + strconv.Itoa(minconf)
				if account != "" {
					r += " from account='" + account
				}
				r += "\n"
			}
			return
		},

		"getreceivedbyaddress": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				minconf := 1
				address := ""
				for i := range f {
					if strings.HasPrefix(f[i], "minconf=") {
						split := strings.Split(f[i], "=")
						n, i := isInteger(split[1])
						if i {
							minconf = n
						} else {
							return failure(mode, err, "minconf value is not a number")
						}
					} else {
						address = f[i]
					}
				}
				r += "Getting balance with minconf=" + strconv.Itoa(minconf)
				r += " from '" + address + "'\n"
			}
			return
		},

		"gettransaction": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 1 {
					return failure(mode, err, "Wrong number of arguments, wanted <txid> got ", f)
				}
				r += "Getting transaction with txid '" + f[0] + "'\n"
			}
			return
		},

		"gettxout": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				includemempool := false
				n := -1
				txid := ""
				cond := false
				if len(f) > 3 {
					return failure(mode, err, "Too many arguments, want <txid> <n> [includemempool] got: ", f)
				}
				if len(f) < 2 {
					return failure(mode, err, "Require at minimum <txid> <n> ", len(f), f)
				}
				for i := range f {
					if f[i] == "includemempool" {
						if includemempool {
							return failure(mode, err, "Already set includemempool flag")
						}
						includemempool = true
					} else {
						if n == -1 {
							n, cond = isInteger(f[i])
							if !cond {
								txid = f[i]
								n = -1
							}
						} else {
							_, cond = isInteger(f[i])
							if cond {
								return failure(mode, err, "Cannot set <n> twice")
							}
							if txid != "" {
								return failure(mode, err, "cannot set <txid> twice")
							}
							txid = f[i]
						}
					}
				}
				if n == -1 {
					return failure(mode, err, "did not set txout number")
				}
				if txid == "" {
					return failure(mode, err, "did not set txid")
				}
				fmt.Println("getting details about txid:", txid, "number:", n)
				if includemempool {
					fmt.Println("  including mempool")
				}
			}
			return
		},

		"gettxoutsetinfo": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Getting txout set info ...")
			}
			return
		},

		"getwork": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 1 {
					return failure(mode, err, "Too many arguments, want [data] got: ", f)
				}
				if len(f) == 1 {
					fmt.Print("Using data\n", f[0], "\n")
				}
				fmt.Println("Getting work ...")
			}
			return
		},

		"help": func(mode string, f []string, erri error) (r string, err error) {
			fmt.Println("Help")
			err = erri
			switch mode {
			case "cli", "":
				subcmd.HelpCommand()
			}
			return
		},

		"importprivkey": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				rescan := false
				label := ""
				if f[0] == "" {
					return failure(mode, err, "<parallelcoinprivkey> is a mandatory parameter")
				}
				for i := 1; i < len(f); i++ {
					if f[i] == "rescan" {
						if rescan {
							return failure(mode, err, "Cannot set rescan more than once")
						}
						rescan = true
					} else {
						if label == "" {
							label = f[i]
						} else {
							return failure(mode, err, "Already set label, cannot set twice")
						}
					}
				}
				if rescan {
					fmt.Println("Rescanning enabled")
				}
				if label != "" {
					fmt.Println("Adding label", label)
				}
				fmt.Println("Importing private key", f[0], " ...")
			}
			return
		},

		"importwallet": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 1 {
					return failure(mode, err, "parameter <filename> is mandatory")
				}
				fmt.Printf("Importing wallet from '%s'", f[0])
			}
			return
		},

		"keypoolrefill": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Refilling keypool ...")
			}
			return
		},

		"listaccounts": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				minconf := 1
				if len(f) > 1 {
					return failure(mode, err, "Too many arguments, only [minconf=x] permitted")
				}
				if len(f) > 0 {
					if strings.HasPrefix(f[0], "minconf=") {
						split := strings.Split(f[0], "=")
						n, i := isInteger(split[1])
						if i {
							minconf = n
						} else {
							return failure(mode, err, "minconf value is not a number")
						}
					} else {
						return failure(mode, err, "Urecognised argument ", f[0])
					}
				}
				fmt.Println("Listing accounts with balances (minconf=" + strconv.Itoa(minconf) + ")")
			}
			return
		},

		"listaddressgroupings": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Listing address groupings ...")
			}
			return
		},

		"listlockunspent": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Starting miner ...")
			}
			return
		},

		"listreceivedbyaccount": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				minconf := 1
				includeempty := false
				if len(f) > 2 {
					return failure(mode, err, "Too many arguments ")
				}
				for i := range f {
					if strings.HasPrefix(f[i], "minconf=") {
						split := strings.Split(f[i], "=")
						n, i := isInteger(split[1])
						if i {
							minconf = n
						} else {
							return failure(mode, err, "minconf value is not a number")
						}
					} else if f[i] == "includeempty" {
						if includeempty {
							return failure(mode, err, "includeempty already set")
						}
						includeempty = true
					} else {
						return failure(mode, err, "Urecognised argument ", f[0])
					}
				}
				if includeempty {
					fmt.Print("Including empty, l")
				} else {
					fmt.Print("L")
				}
				fmt.Println("isting received by account (minconf=" + strconv.Itoa(minconf) + ")")
			}
			return
		},

		"listreceivedbyaddress": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				minconf := 1
				includeempty := false
				if len(f) > 2 {
					return failure(mode, err, "Too many arguments ")
				}
				for i := range f {
					if strings.HasPrefix(f[i], "minconf=") {
						split := strings.Split(f[i], "=")
						n, i := isInteger(split[1])
						if i {
							minconf = n
						} else {
							return failure(mode, err, "minconf value is not a number")
						}
					} else if f[i] == "includeempty" {
						if includeempty {
							return failure(mode, err, "includeempty already set")
						}
						includeempty = true
					} else {
						return failure(mode, err, "Urecognised argument ", f[0])
					}
				}
				if includeempty {
					fmt.Print("Including empty, l")
				} else {
					fmt.Print("L")
				}
				fmt.Println("isting received by address (minconf=" + strconv.Itoa(minconf) + ")")
			}
			return
		},

		"listsinceblock": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				blockhash := ""
				targetconfirmations := 0
				for i := range f {
					n, t := isInteger(f[i])
					if t {
						if targetconfirmations == 0 {
							targetconfirmations = n
						} else {
							return failure(mode, err, "[target-confirmations] already set")
						}
					} else {
						if blockhash == "" {
							blockhash = f[i]
						} else {
							return failure(mode, err, "Cannot set more than one blockhash")
						}
					}
				}
				if blockhash == "" {
					fmt.Println("Listing all transactions(!) ...")
				} else {
					fmt.Println("Listing since block ", blockhash)
				}
				if targetconfirmations > 0 {
					fmt.Println("with minimum confirmations", targetconfirmations)
				}
			}
			return
		},

		"listtransactions": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				countset := false
				count := 10
				from := 0
				if len(f) > 3 {
					return failure(mode, err, "Excess arguments")
				}
				if len(f) > 0 {
					for i := range f {
						if strings.HasPrefix(f[i], "count=") {
							split := strings.Split(f[i], "=")
							n, i := isInteger(split[1])
							if i {
								if !countset {
									count = n
									countset = true
								} else {
									return failure(mode, err, "Cannot set count more than once")
								}
							} else {
								return failure(mode, err, "count value is not a number")
							}
						}
						if strings.HasPrefix(f[i], "from=") {
							split := strings.Split(f[i], "=")
							n, i := isInteger(split[1])
							if i {
								if from == 0 {
									from = n
								} else {
									return failure(mode, err, "Cannot set from more than once")
								}
							} else {
								return failure(mode, err, "from value is not a number")
							}
						}
					}
				}
				fmt.Println("Listing", count, "most recent transactions")
				if from != 0 {
					fmt.Println("from", from)
				}
			}
			return
		},

		"listunspent": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				t := false
				minconf := 1
				minconfset := false
				maxconf := 9999999
				maxconfset := false
				var addresses []string
				if len(f) > 3 {
					return failure(mode, err, "Excess arguments")
				}
				if len(f) > 0 {
					for i := range f {
						if strings.HasPrefix(f[i], "minconf=") {
							split := strings.Split(f[i], "=")
							n, i := isInteger(split[1])
							if i {
								if !minconfset {
									minconf = n
									minconfset = true
								} else {
									return failure(mode, err, "Cannot set minconf more than once")
								}
							} else {
								return failure(mode, err, "minconf value is not a number")
							}
						} else if strings.HasPrefix(f[i], "maxconf=") {
							split := strings.Split(f[i], "=")
							n, i := isInteger(split[1])
							if i {
								if !maxconfset {
									maxconf = n
									maxconfset = true
								} else {
									return failure(mode, err, "Cannot set maxconf more than once")
								}
							} else {
								return failure(mode, err, "maxconf is not a number")
							}
						} else {
							if !json.Valid([]byte(f[1])) {
								return failure(mode, err, "List of addresses was not formatted correctly")
							}
							addresses, t = isValidJSONArray(f[i])
							if !t {
								return failure(mode, err, "List of addresses not correct JSON array")
							}
						}
					}
				}
				fmt.Println("Listing unspent transactions outputs")
				if addresses != nil {
					fmt.Println("from addresses")
					for i := range addresses {
						fmt.Println(" ", addresses[i])
					}
				}
				fmt.Println("with minconf = ", minconf, "and maxconf = ", maxconf)
			}
			return
		},

		"lockunspent": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 2 {
					return failure(mode, err, "Excess arguments")
				}
				unlock := false
				var arrayofobjects string
				for i := range f {
					if f[i] == "unlock" {
						if unlock {
							return failure(mode, err, "Cannot set unlock more than once")
						}
						unlock = true
					} else {
						arrayofobjects = f[i]
					}
				}
				if arrayofobjects == "" {
					return failure(mode, err, "No list of outputs provided")
				}
				if unlock {
					fmt.Println("Unlocking", arrayofobjects)
				} else {
					fmt.Println("Locking", arrayofobjects)
				}
			}
			return
		},

		"makekeypair": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				fmt.Println("Making public/private keypair ...")
				switch {
				case len(f) == 0:
				case len(f) == 1:
					fmt.Println("Using prefix", f[0])
				default:
					return failure(mode, err, "Excess arguments")
				}
				newKey := key.Priv{}
				newKey.New()
				pubkey := newKey.GetPub()
				privkey := newKey.GetPriv()
				b58, err := base58check.Encode("B2", hex.EncodeToString(privkey.Get()))
				if err != nil {
					return failure(mode, err, "Base58check encoding failure")
				}
				keymap := map[string]string{
					"PublicKey":  hex.EncodeToString(pubkey.Key()),
					"PrivateKey": b58,
				}
				jsonbytes, _ := json.MarshalIndent(keymap, "", "  ")
				fmt.Println(string(jsonbytes))
			}
			return
		},

		"move": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				var fromacct, toacct string
				var amount float64
				minconf := 1
				minconfset := false
				t := false
				comment := ""
				if len(f) < 3 {
					return failure(mode, err, "Insufficient arguments")
				}
				fromacct, toacct = f[0], f[1]
				amount, t = isFloat(f[2])
				if !t {
					return failure(mode, err, "amount is not a number")
				}
				length := len(f)
				for i := 3; i < length; i++ {
					if strings.HasPrefix(f[i], "minconf=") {
						split := strings.Split(f[i], "=")
						n, i := isInteger(split[1])
						if i {
							if !minconfset {
								minconf = n
								minconfset = true
							} else {
								return failure(mode, err, "Cannot set minconf more than once")
							}
						} else {
							return failure(mode, err, "minconf value is not a number")
						}
					} else {
						if comment == "" {
							comment = f[i]
						} else {
							return failure(mode, err, "Cannot set comment twice")
						}
					}
				}
				fmt.Println("Moving amount", amount, "from", fromacct, "to", toacct, "with minconf =", minconf)
				if comment != "" {
					fmt.Println("Comment: '" + comment + "'")
				}
			}
			return
		},

		"sendalert": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) < 6 {
					return failure(mode, err, "Insufficient arguments")
				}
				message := f[0]
				privatekey := f[1]
				minver, t := isInteger(f[2])
				if !t {
					return failure(mode, err, "minver is not a number")
				}
				maxver, t := isInteger(f[3])
				if !t {
					return failure(mode, err, "maxver is not a number")
				}
				priority, t := isInteger(f[4])
				if !t {
					return failure(mode, err, "priority is not a number")
				}
				id := f[5]
				cancelupto := 0
				if len(f) == 7 {
					cancelupto, t = isInteger(f[6])
					if !t {
						return failure(mode, err, "[cancelupto] is not a number")
					}
				}
				if len(f) > 7 {
					return failure(mode, err, "Excess arguments ", f[7:])
				}
				fmt.Println("message: '" + message + "'")
				fmt.Println("private key: '" + privatekey + "'")
				fmt.Println("minver:", minver, "maxver:", maxver)
				fmt.Println("priority:", priority)
				fmt.Println("ID: '" + id + "'")
				if cancelupto != 0 {
					fmt.Println("cancel up to:", cancelupto)
				}
			}
			return
		},

		"sendfrom": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) < 3 {
					return failure(mode, err, "Insufficient arguments")
				}
				fromaccount := f[0]
				toaddress := f[1]
				amount, t := isFloat(f[2])
				if !t {
					return failure(mode, err, "amount is not a number")
				}
				minconf := 1
				comment := ""
				commentto := ""
				if len(f) > 3 {
					for i := 3; i < len(f); i++ {
						switch {
						case strings.HasPrefix(f[i], "minconf="):
							if minconf != 1 {
								return failure(mode, err, "Already set minconf")
							}
							minconf, t = isInteger(strings.Split(f[i], "=")[1])
							if !t {
								return failure(mode, err, "minconf value is not a number")
							}
							if minconf == 1 {
								return failure(mode, err, "no point setting minconf to default")
							}
						case strings.HasPrefix(f[i], "comment="):
							if comment != "" {
								return failure(mode, err, "Cannot set comment more than once")
							}
							comment = strings.Split(f[i], "=")[1]
						case strings.HasPrefix(f[i], "comment-to="):
							if commentto != "" {
								return failure(mode, err, "Cannot set comment-to more than once")
							}
							commentto = strings.Split(f[i], "=")[1]
						default:
							return failure(mode, err, "Unrecognised parameter '"+f[i]+"'")
						}
					}
				}
				fmt.Println("From account '"+fromaccount+"' to address '"+toaddress+"' amount", amount, "minconf =", minconf)
				if comment != "" {
					fmt.Println("Comment: '" + comment + "'")
				}
				if commentto != "" {
					fmt.Println("Comment-to: '" + commentto + "'")
				}
			}
			return
		},

		"sendmany": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) < 2 {
					return failure(mode, err, "<fromaccount> and address list are required")
				}
				fromaccount := f[0]
				if !json.Valid([]byte(f[1])) {
					return failure(mode, err, "Transaction JSON was not formatted correctly")
				}
				addresses, t := isValidJSONMap(f[1])
				if !t {
					return failure(mode, err, "list is not in correct JSON format")
				}
				comment := ""
				minconf := 1
				if len(f) > 2 {
					for i := 2; i < len(f); i++ {
						switch {
						case strings.HasPrefix(f[i], "minconf="):
							if minconf != 1 {
								return failure(mode, err, "Already set minconf")
							}
							minconf, t = isInteger(strings.Split(f[i], "=")[1])
							if !t {
								return failure(mode, err, "minconf value is not a number")
							}
							if minconf == 1 {
								return failure(mode, err, "no point setting minconf to default")
							}
						default:
							if comment != "" {
								return failure(mode, err, "Already set comment")
							}
							comment = f[i]
						}
					}
				}
				fmt.Println("From account:", fromaccount, "to accounts:", addresses, "minconf =", minconf, "comment: '"+comment+"'")
			}
			return
		},

		"sendrawtransaction": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 1 {
					return failure(mode, err, "Must have <hex string> and nothing else")
				}
				fmt.Println("Sending raw transaction\n", f[0])
			}
			return
		},

		"sendtoaddress": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) < 2 {
					return failure(mode, err, "Insufficient arguments")
				}
				address := f[0]
				amount, t := isFloat(f[1])
				if !t {
					return failure(mode, err, "<amount> '"+f[1]+"' needs to be a number")
				}
				var comment, commentto string
				if len(f) > 2 {
					for i := 2; i < len(f); i++ {
						switch {
						case strings.HasPrefix(f[i], "comment="):
							if comment != "" {
								return failure(mode, err, "Cannot set comment more than once")
							}
							comment = strings.Split(f[i], "=")[1]
						case strings.HasPrefix(f[i], "comment-to="):
							if commentto != "" {
								return failure(mode, err, "Cannot set comment-to more than once")
							}
							commentto = strings.Split(f[i], "=")[1]
						default:
							return failure(mode, err, "Unrecognised parameter '"+f[i]+"'")
						}
					}
				}
				fmt.Println("Sending to address", address, "amount", amount)
				if comment != "" {
					fmt.Println("comment '" + comment + "'")
				}
				if commentto != "" {
					fmt.Println("comment-to '" + commentto + "'")
				}
			}
			return
		},

		"setaccount": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				switch {
				case len(f) < 2:
					return failure(mode, err, "Insufficient arguments")
				case len(f) > 2:
					return failure(mode, err, "Excess arguments")
				}
				parallelcoinaddress := f[0]
				account := f[1]
				fmt.Println("Setting account", account, "to associate with address", parallelcoinaddress)
			}
			return
		},

		"setgenerate": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Starting miner ...")
			}
			return
		},

		"settxfee": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 1 {
					return failure(mode, err, "Incorrect number of arguments")
				}
				txfee, t := isFloat(f[0])
				if !t {
					return failure(mode, err, "<amount> is not a number")
				}
				fmt.Println("Setting txfee to", txfee)
			}
			return
		},

		"signmessage": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 2 {
					return failure(mode, err, "Incorrect number of arguments")
				}
				address := f[0]
				message := f[1]
				fmt.Println("Signing message '"+message+"' with key from address", address)
			}
			return
		},

		"signrawtransaction": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) < 1 {
					return failure(mode, err, "Incorrect number of arguments")
				}
				sighashtype := "ALL"
				sighashtypeset := false
				if !json.Valid([]byte(f[0])) {
					return failure(mode, err, "Transaction JSON was not formatted correctly")
				}
				data, _ := isValidJSONMap(f[0])
				dataset := false
				var privkeys []string
				privkeysset := false
				if len(f) > 1 {
					for i := 1; i < len(f); i++ {
						if strings.HasPrefix(f[i], "sighashtype=") {
							if sighashtypeset {
								return failure(mode, err, "Already set sighashtype")
							}
							sighashtype = strings.Split(f[i], "=")[1]
							switch sighashtype {
							case "ALL", "NONE", "SINGLE", "ALL|ANYONECANPAY", "NONE|ANYONECANPAY", "SINGLE|ANYONECANPAY":
								sighashtypeset = true
							default:
								return failure(mode, err, "sighashtype not one of available options")
							}

						}
						if !json.Valid([]byte(f[i])) {
							return failure(mode, err, "Transaction JSON was not formatted correctly")
						}
						datain, t := isValidJSONMap(f[i])
						if t {
							if !dataset {
								dataset = true
								data = datain
							}
						} else if !privkeysset {
							privkeys = strings.Split(f[i], ",")
							privkeysset = true
						}
					}
				}
				fmt.Println("Signing raw transaction ...")
				if dataset {
					fmt.Println("data: ", data)
				}
				if privkeysset {
					fmt.Println("private keys:", privkeys)
				}
				fmt.Println("sighashtype =", sighashtype)
			}
			return
		},

		"stop": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 0 {
					return failure(mode, err, "No arguments are required ")
				}
				fmt.Println("Stopping miner ...")
			}
			return
		},

		"submitblock": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) < 1 {
					return failure(mode, err, "Incorrect number of arguments")
				}
				fmt.Println("Submitting block '" + f[0] + "'")
			}
			return
		},

		"validateaddress": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 1 {
					return failure(mode, err, "Incorrect number of arguments")
				}
				fmt.Println("Validating address", f[0])
			}
			return
		},

		"verifychain": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) > 1 {
					return failure(mode, err, "only one optional parameter allowed")
				}
				var checklevel int
				var t bool
				if f[0] != "" {
					checklevel, t = isInteger(f[0])
					if !t {
						return failure(mode, err, "checklevel is not an integer")
					}
				}
				if t {
					fmt.Println("checklevel set to", checklevel)
				}
				fmt.Println("Verifying chain ...")
			}
			return
		},

		"verifymessage": func(mode string, f []string, erri error) (r string, err error) {
			err = erri
			switch mode {
			case "cli", "":
				if len(f) != 3 {
					return failure(mode, err, "Incorrect number of args")
				}
				address := f[0]
				signature := f[1]
				message := f[2]
				fmt.Println("Verifying message '"+message+"' with signature '"+signature+"' using address", address)
			}
			return
		},
	}
}
