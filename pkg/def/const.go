package def

import (
	"fmt"
	"runtime"
	"time"
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

// Debugging turns on error printing to stdout
var Debugging = true

// Debug prints out debugging information
func Debug(s ...interface{}) {
	if Debugging {
		_, fil, line, _ := runtime.Caller(3)
		S := fmt.Sprint("[", time.Now().Format(time.RFC3339), "] ", fil, ":", line, " ")
		for i := range s {
			S = S + fmt.Sprint(s[i]) + " "
		}
		fmt.Println(S)
	}
}
