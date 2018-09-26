// Package rec is a maybe unnecessary collection of specifications of the data formats used in the walletdb.
package rec

import (
	"github.com/parallelcointeam/duo/pkg/proto"
)

var (
	// TableNames are the list of table names
	TableNames = []string{"MasterKey", "Name", "Tx", "Seed", "Key", "Script", "Pool", "Setting", "Account", "Accounting", "CreditDebit", "BestBlock", "MinVersion", "DefaultKey"}
	// Tables are a map of 64 bit hashes formed from the exact variable names used here, this is used as a translation table
	Tables map[string]KeyPrefix
)

func init() {
	Tables = make(map[string]KeyPrefix)
	for i := range TableNames {
		t := []byte(TableNames[i])
		Tables[TableNames[i]] = KeyPrefix(*proto.Hash64(&t))
	}
}

// Idx is an 8 byte highwayhash of the encrypted data
type Idx []byte

// KeyPrefix is an 8 byte string that stores the table identifier
type KeyPrefix []byte

// MasterKey is the password-encrypted store for the symmetric ciphertext that seecures a wallet database
type MasterKey struct {
	Idx        Idx // in key
	Crypt      []byte
	IV         [12]byte
	Iterations int64
}

// Name is a human readable string label for an address. Idx is the Highwayhash 64 of the encrypted address, used in the accounts and keys for quick matching
type Name struct {
	Idx     Idx    // in key
	Address []byte // encrypt // in key
	Label   string // encrypt
}

// Tx is a transaction connected to an address in the wallet
type Tx struct {
	Idx                   Idx
	AcIdxs                []Idx          //in key
	ID                    [32]byte       // encrypt
	Data                  []byte         // encrypt
	Prev                  proto.MerkleTx // encrypt
	TimeRecvIsTxTime      int64          // encrypt
	TimeRecv              int64          // encrypt
	FromMe                bool
	Accounts              [][]byte // encrypt
	Spent                 []byte   // encrypt
	OrderPos              int64
	DebitCached           bool
	CreditCached          bool
	ImmatureCreditCached  bool
	AvailableCreditCached bool
	ChangeCached          bool
}

// Seed is a 'hierarchic deterministic' wallet seed that can spawn many subkeys
type Seed struct {
	Idx    Idx
	Secret [32]byte // encrypt
}

// Key is a key pair controlled by the user. The 64 bit highwayhash of the (encrypted) address and the encrypted address both live in the key so this and Account keys act as an index for all keys in the wallet in other fields, such as the Accounting, which has all of its addresses HWhashes
type Key struct {
	Idx     Idx    // in key
	Address []byte // encrypt      // in key
	Priv    []byte // encrypt
	Pub     []byte // encrypt
}

// Script is a script stored in the wallet database. The hwh64 index and key are in the key so it is quick to find and if necessary, copy back the address after decryption.
type Script struct {
	Idx  Idx      // in key
	ID   [20]byte // encrypt // in key
	Data []byte   // encrypt
}

// Pool is a wallet key pair that has not yet been put to use in a transaction. The Idx is the HWH of the encrypted key ID (address)
type Pool struct {
	Idx     Idx   // in key
	Seq     int64 // in key
	Key           // encrypt
	Created int64
	Expires int64
}

// Account is an address and public key of a counterparty account, for which the user does  not have a private key. The public key is filled in only when the address is found in the chain and the signature is available to recover the key.
type Account struct {
	Idx     Idx      // in key
	Address [20]byte // encrypt      // in key
	Pub     [33]byte // encrypt
}

// Accounting is an entry regarding internal movements of funds. The index contains the 64 bit highway hash of the encrypted addresses in the accounting entry, so scanning the ledger for relevant entries is fast and happens all in memory.
type Accounting struct {
	Idx          []Idx      // in key
	Sequence     int64      // in key
	Account      [][20]byte // encrypt
	CreditDebit  int64      // encrypt
	Timestamp    int64      // encrypt
	OtherAccount [][20]byte // encrypt
	Comment      string     // encrypt
	OrderPos     int64
	EntryNo      int64
	Extra        []byte // encrypt
}

// CreditDebit is an increase or decrease to the balance in the wallet
type CreditDebit struct {
	Idx     Idx      // in key
	Address [20]byte // encrypt
	Amount  int64    // encrypt
}

// BestBlock is the latest block at the time the wallet was last open
type BestBlock struct {
	Height uint64   // in key
	ID     [32]byte // in key
	Data   []byte
}

// MinVersion stores the minimum version number required to support the wallet
type MinVersion struct {
	Number int64
}

// DefaultKey sets the default key for receiving payments to use in interfaces
type DefaultKey struct {
	Address [20]byte // encrypt
}
