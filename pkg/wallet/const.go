package wallet

const (
	// CurrentVersion is the version number from this source repository
	CurrentVersion = 1
	// FeatureBase is the base version number for a wallet
	FeatureBase = 10500
	// FeatureWalletCrypt indicates if the wallet enables encrypted keys
	FeatureWalletCrypt = 40000
	// FeatureCompressedPubKey indicates if the wallet enables compressed public keys
	FeatureCompressedPubKey = 60000
	// FeatureLatest is the newest version of the wallet
	FeatureLatest = 60000
)

var (
	// AccountingEntryNumber is
	AccountingEntryNumber = 0
)
