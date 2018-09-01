package interfaces

/////////////////////////////////////////
// Coding interface
/////////////////////////////////////////

// CodeTypes is the types of encoding available, append only to add new ones for compatibility
var CodeType = []string{
	"byte",
	"string",
	"decimal",
	"hex",
	"base32",
	"base58",
	"base58check",
	"base64",
}