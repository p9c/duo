package proto

// StringCodings is a list of the names of encoding types for the stringer
var StringCodings = []string{
	"bytes",       // golang format, square brackets and space separated decimal up to 255
	"string",      // the raw bytes as a UTF-8 string
	"decimal",     // the value of the bytes as a single number in decimal
	"hex",         // the value of the bytes as a single number in hexadecimal
	"base32",      // the value of the bytes as a single number in base32
	"base58check", // the value of the bytes as a single number in base58 check with the first byte as the prefix
	"base64",      // the value of the bytes as a single number in base64
}

// CommonErrors are common error values from library functions
type CommonErrors struct {
	NilRec, NilParam, NilBuf, ZeroLen, OutOfBounds, InvalidType, ZeroLenBuf string
}

// Errors gives a short readable reference to indicate a common error string
var Errors = CommonErrors{
	"nil receiver", "nil parameter", "nil buffer", "zero length parameter", "index out of bounds", "invalid type", "zero length buffer",
}
