// Package state gathers every variable and constant required for a parallelcoin node. In the original codebase, these are scattered through several files.
package state

import (
	"os"
	"sync"

	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/alert"
	"gitlab.com/parallelcoin/duo/pkg/algos"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/serialize"
	"gitlab.com/parallelcoin/duo/pkg/tx"
	"gitlab.com/parallelcoin/duo/pkg/txdb"
	"gitlab.com/parallelcoin/duo/pkg/wallet"
)

const (
	// Coin is the denomination of one token
	Coin int64 = 100000000
	// Cent is 1/100th of one token
	Cent int64 = 1000000
	// MaxBlockSize is the maximum block size
	MaxBlockSize = 2000000
	// MaxBlockSizeGen -
	MaxBlockSizeGen = MaxBlockSize / 2
	// MaxStandardTxSize -
	MaxStandardTxSize = MaxBlockSizeGen / 5
	// MaxBlockSigOps is the maximum number of signature operations in a block
	MaxBlockSigOps = MaxBlockSize / 50
	// MaxOrphanTransactians is the maximum number of orphan transactions stored
	MaxOrphanTransactians = MaxBlockSize / 100
	// Megabyte is the number of bytes in one megabyte
	Megabyte = 0x100000
	// MaxBlockfileSize is the maximum size of a block file
	MaxBlockfileSize = Megabyte * 128
	// MaxBlockfileChunkSize is the largest size of a chunk in the Blockfile
	MaxBlockfileChunkSize = Megabyte * 16
	// UndofileChunkSize is thte size of the undo buffer chunks
	UndofileChunkSize = Megabyte
	// MempoolHeight is a special code signifying a transaction is in the mempool
	MempoolHeight = 0x7FFFFFFF
	// MaxMoney is the maximum number of tokens ever
	MaxMoney int64 = 1000000 * Coin
	// CoinbaseMaturity is the number of blocks that must be confirmed before a newly mined token can be spent
	CoinbaseMaturity = 100
	// LocktimeThreshold is
	LocktimeThreshold = 500000000
	// MaxScriptcheckThreads is the maximum number of threads that will be spawned to check scripts
	MaxScriptcheckThreads = 16
	// DefaultBlockPrioritySize is
	DefaultBlockPrioritySize = 27000
	// MinDiskSpace is the minimum amount of disk space remaining before the server will cease operating
	MinDiskSpace = 52428800
	// BlockAlgoWorkWeightStart is
	BlockAlgoWorkWeightStart = 142000
	// MessageMagic is the marker of a message signed or encrypted using parallelcoin keys
	MessageMagic = "Parallelcoin Signed Message:\n"
	// StartSubsidy is the number of tokens won in mining rewards from genesis
	StartSubsidy = 2 * Coin
	// MinSubsidy is the 'fair release' reduced reward in early blocks
	MinSubsidy = Coin / 200
	// TargetTimespan is the number of milliseconds between blocks (5 minutes)
	TargetTimespan int64 = 30000
	// TargetSpacing is
	TargetSpacing int64 = 300
	// Interval is
	Interval int64 = 100
	// AveragingInterval is the number of blocks averaged to compute the difficulty
	// Long to stop difficulty getting stuck high with wide difference in hash power
	AveragingInterval int64 = 2047
	// AveragingTargetTimespan is
	AveragingTargetTimespan int64 = AveragingInterval * TargetSpacing
	// MaxAdjustDown is the maximum amount the difficulty can adjust down after a block
	MaxAdjustDown int64 = 20 // 8%
	// MaxAdjustUp is the maximum amount the difficulty can adjust up after a block
	MaxAdjustUp int64 = 80 // 8%
	// TargetTimespanAdjDown is
	TargetTimespanAdjDown int64 = TargetTimespan * (100 + MaxAdjustDown) / 100
	// MinActualTimespan is
	MinActualTimespan = AveragingTargetTimespan * (100 - MaxAdjustUp) / 100
	// MaxActualTimespan is
	MaxActualTimespan = AveragingTargetTimespan * (100 + MaxAdjustDown) / 100
	// GMFRelay  - GetMinFee mode
	GMFRelay = iota
	// GMFSend is MinFee send
	GMFSend
	// MedianTimeSpan is
	MedianTimeSpan = 11
	// ModeValid -
	ModeValid = iota
	// ModeInvalid -
	ModeInvalid
	// ModeError -
	ModeError
	// P2SH -
	P2SH = "/P2SH/"
	// BindFlagsNone -
	BindFlagsNone = 0
	// BindFlagsExplicit -
	BindFlagsExplicit = 1
	// BindFlagsReportError -
	BindFlagsReportError = 2
)

var (
	// WalletFilename is the centrally stored filename of the wallet in a server instance
	WalletFilename string
	// WalletMain is the main wallet of a server instance
	WalletMain *wallet.Wallet
	// RequestShutdown is a trigger that represents whether a shutdown switch has been flipped
	RequestShutdown = false
	// CoinsDBView is
	CoinsDBView *txdb.CoinsViewDB
	// DebugFile is the file that the debug logs are written to
	DebugFile os.File
	// BlockStatus is a map of various block states
	BlockStatus = map[string]int{
		"ValidUnknown":      0,
		"ValidHeader":       1,
		"ValidTree":         2,
		"ValidTransactions": 3,
		"ValidChain":        4,
		"ValidScripts":      5,
		"ValidMask":         7,
		"HaveData":          8,
		"HaveUndo":          16,
		"HaveMask":          24,
		"FailedValid":       32,
		"FailedChild":       64,
		"FailedMask":        96,
	}
	// WalletRegisteredMutexSet is a set of mutexes for the wallet
	WalletRegisteredMutexSet sync.RWMutex
	// WalletRegisteredSet is the set of registered wallets
	WalletRegisteredSet wallet.Wallet
	// MainMutex is the main mutex of the server
	MainMutex sync.RWMutex
	// MemPool is the transaction mempool
	MemPool tx.MemPool
	// TransactionsUpdated is a count number of the number of updates done on the database since loading
	TransactionsUpdated uint64
	// BlockIndexMap is a map of block indices
	BlockIndexMap map[*Uint.U256]*block.Index
	// BlockIndexesByHeight is a blocks sorted by height
	BlockIndexesByHeight []block.Index
	// IndexGenesisBlock is the index of the genesis block
	IndexGenesisBlock *block.Index
	// BestHeight is the newest confirmed block
	BestHeight = -1
	// BestChainWork is
	BestChainWork Uint.U256
	// BestInvalidWork is
	BestInvalidWork Uint.U256
	// HashBestChain is the hash of the best block
	HashBestChain Uint.U256
	// IndexBest is the index of the best block
	IndexBest *block.Index
	// BlockIndexValidSet is the list of valid block indices
	BlockIndexValidSet map[*block.Index]bool
	// TimeBestReceived is the time the newest accepted block wast received
	TimeBestReceived int64
	// ScriptCheckThreads is the number of script checking threads
	ScriptCheckThreads int
	// Importing signifies if the server is importing blocks
	Importing = false
	// Reindex signifies if the server is reindexing blocks
	Reindex = false
	// Benchmark signifies if the server is running a benchmark
	Benchmark = false
	// TxIndex signifies if the server is indexing transactions
	TxIndex = false
	// CoinCacheSize is
	CoinCacheSize uint = 5000
	// Transaction is
	Transaction = tx.Transaction{
		MinTxFee:      10000,
		MinRelayTxFee: 10000,
	}
	// OrphanBlocksMap is a mak of orphan blocks
	OrphanBlocksMap map[*Uint.U256]*block.Block
	// OrphanBlocksByPrevMap is the orphans sorted by their previous block
	OrphanBlocksByPrevMap map[*Uint.U256]*block.Block
	// OrphanTransactionsMap is a map of orphan transactions
	OrphanTransactionsMap map[*Uint.U256]*ser.DataStream
	// OrphanTransactionsByPrevMap is a map of orphan transactions sorted by prev
	OrphanTransactionsByPrevMap map[*Uint.U256]*ser.DataStream
	// CoinbaseFlags is
	CoinbaseFlags key.Script
	// HashesPerSec is the current hashrate if the server is mining
	HashesPerSec float64
	// HPSTimerStart is
	HPSTimerStart int64
	// TransactionFee is
	TransactionFee int64
	// MiningAlgo is (defaults to SHA256D) the algorithm the server is mining
	MiningAlgo = algos.SHA256D
	// CoinsTip is
	CoinsTip block.CoinsViewCache
	// BlockTree is
	BlockTree txdb.BlockTreeDB
	// InOutPoints is
	InOutPoints tx.OutPoint
	// BlockIndexFBBHLast is
	BlockIndexFBBHLast block.Index
	// LastBlockFileMutex is
	LastBlockFileMutex sync.RWMutex
	// InfoLastBlockFile is
	InfoLastBlockFile block.FileInfo
	// LastBlockFile is
	LastBlockFile int
	// AlertsMap is
	AlertsMap map[*Uint.U256]alert.Alert
	// AlertsMutex is
	AlertsMutex sync.RWMutex
	// LastBlockTx is
	LastBlockTx uint64
	// LastBlockSize is
	LastBlockSize uint64
)

func Init() {}
