// A library for querying proof of work algorithm metadata
package algos
const (
	// SHA256D is proof of work using double SHA256 hashes
	SHA256D = iota
	// SCRYPT is proof of work using Scrypt hashes
	SCRYPT
	// NumAlgos is the number of algorithms currently supported
	NumAlgos
	// BlockVersionDefault is the current default block version
	BlockVersionDefault = 2
	// BlockVersionAlgo is
	BlockVersionAlgo = (1 << 9)
	// BlockVersionScrypt is
	BlockVersionScrypt = (1 << 9)
)
// Returns the string identifier of the PoW algorithm code number submitted
func Name(algo int) string {
	switch algo {
	case SHA256D:
		return "sha256d"
	case SCRYPT:
		return "scrypt"
	}
	return "unknown"
}
// Accepts the string name of an algorithm and returns the code that identifies it
func Code(s string) int {
	switch s {
	case "sha256d":
		return SHA256D
	case "scrypt":
		return SCRYPT
	}
	return 0
}
