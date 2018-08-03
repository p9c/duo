package args

import (
	"flag"
	"fmt"
	"strings"

	"gitlab.com/parallelcoin/duo/pkg/homedir"
	"gitlab.com/parallelcoin/duo/pkg/iniflags"
)

// nolint
var (
	Version = flag.Bool("version",
		false,
		"Print version information")
	Help = flag.Bool("help",
		false,
		"Show help information")
	Daemon = flag.Bool("daemon",
		false,
		"Detach from current tty and run in the background")
	Gen = flag.Bool("gen",
		false,
		"Generate blocks")
	DataDir = flag.String("datadir",
		".duo",
		"Specify data directory (if relative, user's home directory is set as the base)")
	Conf = flag.String("conf",
		"config.ini",
		"Specify config filename (if relative ~/<datadir> is set as the base)")
	Wallet = flag.String("wallet",
		"wallet.dat",
		"Specify wallet file")
	WalletPassword = flag.String("walletpassword",
		"",
		"Password to decrypt wallet.dat")
	DBCache = flag.Uint64("dbcache",
		100,
		"Size of database cache in megabytes")
	Timeout = flag.Uint64("timeout",
		5000,
		"Specify connection timeout in milliseconds")
	Proxy = flag.String("proxy",
		"",
		"Connect through socks proxy")
	SocksVersion = flag.Uint("socks",
		5,
		"Select the version of socks proxy to use (4|5)")
	Tor = flag.String("tor",
		"",
		"Use proxy to reach tor hidden services (default: same as -proxy)")
	DNS = flag.Bool("dns",
		false,
		"Allow DNS lookups for -addnode, -seednode and -connect")
	Port = flag.Uint("port",
		11047,
		"Listen for connections on <port> (default: 11047 or testnet: 21047)")
	MaxConnections = flag.Uint("maxconnections",
		125,
		"Specify maximum number of connections to maintain")
	// These are possible to set multiply so must be processed in init()
	AddNodes,
	Connects,
	SeedNodes,
	RPCAllowIPs MultiArg
	ExternalIP = flag.String("externalip",
		"",
		"Specify your own public address")
	OnlyNet = flag.String("onlynet",
		"",
		"Only connect to nodes in network <net> (IPv4, IPv6 or Tor)")
	Discover = flag.Bool("discover",
		false,
		"Discover own IP address (default: true when listening and no -externalip)")
	Checkpoints = flag.Bool("checkpoints",
		true,
		"Only accept block chain matching built-in checkpoints")
	Listen = flag.Bool("listen",
		true,
		"Accept connections from outside (default: true if no -proxy or -connect)")
	Bind = flag.String("bind",
		"",
		"Bind to given address and always listen on it. Use [host]:port notation for IPv6")
	DNSSeed = flag.Bool("dnsseed",
		true,
		"Find peers using DNS lookup (default: true unless -connect)")
	BanScore = flag.Uint("banscore",
		100,
		"Threshold for disconnecting misbehaving peers")
	BanTime = flag.Uint64("bantime",
		86400,
		"Number of seconds to keep misbehaving peers from reconnecting")
	MaxReceiveBuffer = flag.Uint("maxreceivebuffer",
		5000,
		"Maximum per-connection receive buffer, x 1000 bytes")
	MaxSendBuffer = flag.Uint("maxsendbuffer",
		1000,
		"Maximum per-connection send buffer, x 1000 bytes")
	UPnP = flag.Bool("upnp",
		true,
		"Use UPnP to map the listening port (default: true when listening)")
	PayTxFee = flag.Uint64("paytxfee",
		0,
		"Fee per KB to add to transactions you send")
	MinTxFee = flag.Uint("mintxfee",
		1000,
		"Minimum fee per kilobyte of transactions")
	MinRelayTxFee = flag.Uint("minrelaytxfee",
		0,
		"Minimum fee to relay kilobyte of transactions")
	TestNet = flag.Bool("testnet",
		false,
		"Use the test network")

	Debug = flag.Bool("debug",
		false,
		"Output extra debugging information. Implies all other -debug* options")
	DebugNet = flag.Bool("debugnet",
		false,
		"Output extra network debugging information")
	LogTimestamps = flag.Bool("logtimestamps",
		true,
		"Prepends timestamps on debug output")
	ShrinkDebugFile = flag.Bool("shrinkdebugfile",
		true,
		"Shrink debug.log file on client startup (default: true when no -debug)")
	PrintToConsole = flag.Bool("printtoconsole",
		true,
		"Send trace/debug info to console instead of debug.log file")
	RegTest = flag.Bool("regtest",
		false,
		"Enter regression test mode, which uses a special chain in which blocks can be solved instantly.\n"+
			"This is intended for regression testing tools and app development.")
	RPCUser = flag.String("rpcuser",
		"",
		"Username for JSON-RPC connections")
	RPCPassword = flag.String("rpcpassword",
		"",
		"Password for JSON-RPC connections")
	RPCPort = flag.Uint("rpcport",
		11048,
		"Listen for JSON-RPC connections on <port> (default: 11048 or testnet: 21048)")
	RPCConnect = flag.String("rpcconnect",
		"127.0.0.1",
		"Send commands to node running on IP")
	RPCThreads = flag.Uint("rpcthreads",
		4,
		"Specify the number of threads to service RPC calls (default: 4)")
	BlockNotify = flag.String("blocknotify",
		"",
		"Execute command when the best block changes (%s in cmd is replaced by block hash)")
	WalletNotify = flag.String("walletnotify",
		"",
		"Execute command when a wallet transaction changes (%s in cmd is replaced by TxID)")
	UpgradeWallet = flag.Bool("upgradewallet",
		false,
		"Upgrade wallet to latest format")
	KeyPool = flag.Uint("keypool",
		100,
		"Specify key pool size")
	Rescan = flag.Bool("rescan",
		false,
		"Rescan the block chain for missing wallet transactions")
	SalvageWallet = flag.Bool("salvagewallet",
		false,
		"Attempt to recover private keys from a corrupt wallet.dat")
	CheckBlocks = flag.Uint64("checkblocks",
		288,
		"How many blocks to check at startup (0 = all)")
	CheckLevel = flag.Uint("checklevel",
		3,
		"How thorough the block verification is (0-4)")
	TxIndex = flag.Bool("txindex",
		false,
		"Maintain a full transaction index")
	LoadBlock = flag.String("loadblock",
		"",
		"Imports blocks from external blk000??.dat file")
	Reindex = flag.Bool("reindex",
		false,
		"Rebuild block chain index from current blk000??.dat files")
	Par = flag.Int("par",
		0,
		"Set the number of script verification threads (up to 16, 0 = auto, <0 = leave that many cores free)")
	Algo = flag.String("algo",
		"sha256d",
		"Mining algorithm: sha256d, scrypt")
	BlockMinSize = flag.Uint64("blockminsize",
		0,
		"Set minimum block size in bytes")
	BlockMaxSize = flag.Uint64("blockmaxsize",
		250000,
		"Set maximum block size in bytes")
	BlockPrioritySize = flag.Uint64("blockprioritysize",
		27000,
		"Set maximum size of high-priority/low-fee transactions in bytes")
	RPCSSL = flag.Bool("rpcssl",
		false,
		"Use OpenSSL (https) for JSON-RPC connections")
	RPCSSLCertificateChainFile = flag.String("rpcsslcertificatechainfile",
		"server.cert",
		"Server certificate file")
	RPCSSLPrivateKeyFile = flag.String("rpcsslprivatekeyfile",
		"server.pem",
		"Server private key")
	RPCSSLLCiphers = flag.String("rpcsslciphers",
		"TLSv1+HIGH:!SSLv2:!aNULL:!eNULL:!AH:!3DES:@STRENGTH",
		"Acceptable ciphers")
	CreateConf = flag.Bool("configure",
		false,
		"Create a default configuration file where -conf sets it")
	// commands
	// helpCommand = flag.NewFlagSet("help", flag.ExitOnError)
	// stopCommand = flag.NewFlagSet("stop", flag.ExitOnError)
)

// MultiArg is a list of argument strings split from a CLI input
type MultiArg struct {
	Value []string
}

// Set adds  an argument to a list of arguments
func (i *MultiArg) Set(value string) error {
	i.Value = append(i.Value, value)
	return nil
}

// String returns a list of arguments back into one string
func (i *MultiArg) String() string {
	return fmt.Sprint(*i)
}

func init() {
	flag.Var(&AddNodes, "addnode", "Add a node to connect to and attempt to keep the connection open")
	flag.Var(&Connects, "connect", "Connect only to the specified node(s)")
	flag.Var(&SeedNodes, "seednode", "Connect to a node to retrieve peer addresses, and disconnect")
	flag.Var(&RPCAllowIPs, "rpcallowip", "Allow JSON-RPC connections from specified IP address")
	if !strings.HasPrefix(*DataDir, "/") {
		homeDir, _ := homedir.Dir()
		*DataDir = homeDir + "/" + *DataDir
	}
	if !strings.HasPrefix(*Conf, "/") {
		*Conf = *DataDir + "/" + *Conf
		iniflags.SetConfigFile(*Conf)
	}

}
