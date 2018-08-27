package byteprint

import (
	"fmt"
	"os"
)

// BP is a short prefix for invoking the byte slice print functions
var BP *BytePrint

// BytePrint outputs string data contained in byte slices directly to Stdout
type BytePrint struct{}
type bytePrint interface {
	Str(...interface{}) *BytePrint
	Apo(...interface{}) *BytePrint
	Quo(...interface{}) *BytePrint
	Brc(...interface{}) *BytePrint
	Brk(...interface{}) *BytePrint
	Par(...interface{}) *BytePrint
	CR()
	SP()
}

// Str prints byte slices and strings to stdout
func (r *BytePrint) Str(b ...interface{}) *BytePrint {
	Print(b...)
	return r
}

// Apo prints byte slices and strings to stdout surrounded by single quotes
func (r *BytePrint) Apo(b ...interface{}) *BytePrint {
	fmt.Print("'")
	Print(b...)
	fmt.Print("")
	return r
}

// Quo prints byte slices and strings to stdout surrounded by double quotes
func (r *BytePrint) Quo(b ...interface{}) *BytePrint {
	fmt.Print("\"")
	Print(b...)
	fmt.Print("\"")
	return r
}

// Brc prints byte slices and strings to stdout surrounded by braces ()
func (r *BytePrint) Brc(b ...interface{}) *BytePrint {
	fmt.Print("(")
	Print(b...)
	fmt.Print(")")
	return r
}

// Brk prints byte slices and strings to stdout surrounded by brackets []
func (r *BytePrint) Brk(b ...interface{}) *BytePrint {
	fmt.Print("[")
	Print(b...)
	fmt.Print("]")
	return r
}

// Par prints byte slices and strings to stdout surrounded by parenthesis {}
func (r *BytePrint) Par(b ...interface{}) *BytePrint {
	fmt.Print("{")
	Print(b...)
	fmt.Print("}")
	return r
}

// Print is same name but in package namespace can start a chain of print commands
func Print(b ...interface{}) *BytePrint {
	for i := range b {
		switch b[i].(type) {
		case *string:
			fmt.Print(*b[i].(*string))
		case string:
			fmt.Print(b[i].(string))
		case *[]byte:
			os.Stdout.Write(*b[i].(*[]byte))
		case rune:
			os.Stdout.Write([]byte{byte(int(b[i].(rune)))})
		}
	}
	return nil
}

// CR prints a carriage return
func (r *BytePrint) CR() *BytePrint {
	fmt.Println()
	return r
}

// SP prints a space
func (r *BytePrint) SP() *BytePrint {
	fmt.Print(" ")
	return r
}
