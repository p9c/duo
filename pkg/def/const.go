package def

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	// "time"
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
		red := color.New(color.FgRed).Add(color.Bold)
		green := color.New(color.FgGreen)
		// header := color.New(color.FgWhite)
		hyperlink := color.New(color.FgBlue).Add(color.Underline)
		pc, fil, line, _ := runtime.Caller(3)
		// header.Print("[", time.Now().Format(time.RFC3339), "] ")
		var temp string
		for i := range s {
			temp += s[i].([]interface{})[0].(string)
		}
		green.Print(fmt.Sprint(runtime.FuncForPC(pc).Name(), "() "))
		red.Print(temp + " ")
		hyperlink.Print(fil, ":", line)
		fmt.Println()
	}
}
