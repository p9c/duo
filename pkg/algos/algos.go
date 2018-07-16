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

// Get the type of PoW algorithm
func Get(version int) int {
	switch version & BlockVersionAlgo {
	case 0:
		return SHA256D
	case BlockVersionDefault:
		return SCRYPT
	}
	return SHA256D
}

// Name returns the string identifier of the PoW algorithm
func Name(algo int) string {
	switch algo {
	case SHA256D:
		return "sha256d"
	case SCRYPT:
		return "scrypt"
	}
	return "unknown"
}

// Code accepts the name and returns the code
func Code(s string) int {
	switch s {
	case "sha256d":
		return SHA256D
	case "scrypt":
		return SCRYPT
	}
	return 0
}
