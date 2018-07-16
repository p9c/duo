package addr

var (
	// TriedBucketCount is the number of buckets tried
	TriedBucketCount = 64
	// TriedBucketSize is the size of the buckets
	TriedBucketSize = 64
	// NewBucketCount is the count of new buckets
	NewBucketCount = 256
	// NewBucketSize is the size of new buckets
	NewBucketSize = 64
	// TriedBucketsPerGroup is the number of buckets tried per group
	TriedBucketsPerGroup = 4
	// NewBucketsPerSourceGroup is the number of new buckets per source group
	NewBucketsPerSourceGroup = 32
	// NewBucketsPerAddress is the number of buckets allocated per address
	NewBucketsPerAddress = 4
	// TriedEntriesInspectOnEvict is the number of entries ot try before evicting a record
	TriedEntriesInspectOnEvict = 4
	// HorizonDays is the number of days before a record is purged
	HorizonDays = 60
	// Retries is the number of times we retry
	Retries = 3
	// MaxFailures is the maximum number of failures before dropping a connection
	MaxFailures = 20
	// MinFailDays is the number of days in which the failures can count up
	MinFailDays = 7
	// GetAddrMaxPCT is
	GetAddrMaxPCT = 23
	// GetAddrMax is
	GetAddrMax = 2500
)
