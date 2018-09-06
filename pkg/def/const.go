package def

import (
	"fmt"
)

// StringCodingTypes is the types of encoding available, append only to add new ones for compatibility
var StringCodingTypes = []string{
	"byte",
	"string",
	"decimal",
	"hex",
	"base32",
	"base58",
	"base58check",
	"base64",
}

// Debug turns on error printing to stdout
var Debugging = true

func Debug(s ...interface{}) {
	if Debugging {
		fmt.Println(s...)
	}
}
