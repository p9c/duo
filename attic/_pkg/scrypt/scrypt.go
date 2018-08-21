// An interface for the Scrypt hash function
package scrypt
const (
	// ScratchpadSize is
	ScratchpadSize = 131072 + 63
)
var (
	// HMACSHA256Ctx is
	HMACSHA256Ctx HMACSHA256Context
)
// HMACSHA256Context is
type HMACSHA256Context struct {
}
