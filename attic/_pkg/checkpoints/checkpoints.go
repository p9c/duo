// A library for managing checkpoints for the DUO token ledger blockchain
package checkpoints
import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
)
const (
	// SigcheckVerificationFactor is
	SigcheckVerificationFactor float64 = 5.0
)
var (
	// Enabled is true if checkpoints are enabled (default)
	Enabled = true
	// Checkpoints is the list of block hashes and heights that must match our expectations (it's a kinda cheap centralised security provided by developers)
	Checkpoints = Map{
		0: new(Uint.U256).FromString("0x000009f0fcbad3aac904d3660cfdcf238bf298cfe73adf1d39d14fc5c740ccc7"),
	}
	// MetaData is metadata about the checkpoints
	MetaData = Data{
		CheckpointsMap:            Checkpoints,
		TimeLastCheckpoint:        1405742300,
		TransactionLastCheckpoint: 0,
		TransactionsPerDay:        2880,
	}
	// TestnetCheckpoints is the list of checkpoints for the testnet
	TestnetCheckpoints = Map{
		0: new(Uint.U256).FromString("0x000009f0fcbad3aac904d3660cfdcf238bf298cfe73adf1d39d14fc5c740ccc7"),
	}
	// TestnetData is the metadata about the checkpoints for the testnet
	TestnetData = Data{
		CheckpointsMap:            TestnetCheckpoints,
		TimeLastCheckpoint:        1405741700,
		TransactionLastCheckpoint: 0,
		TransactionsPerDay:        2880,
	}
)
type Map map[int]*Uint.U256
type Data struct {
	CheckpointsMap                                Map
	TimeLastCheckpoint, TransactionLastCheckpoint int64
	TransactionsPerDay                            float64
}
