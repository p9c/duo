// An interface to a bloom filter for optimising various types of data access
package bloom
const (
	// UpdateNone updates none
	UpdateNone = 0
	// UpdateAll updates all
	UpdateAll = 1
	// UpdateP2PubkeyOnly updates only P2Pubkeys
	UpdateP2PubkeyOnly = 2
	// UpdateMask is the bits involved in a bloom update flag
	UpdateMask = 3
)
type Filter struct {
	data             []byte
	HashFuncs, Tweak uint
	Flags            byte
}
