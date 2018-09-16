package def

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"strings"
	"time"
)

// Debug is an object that can be embedded into types to provide a debug print service for error status reporting
type Debug struct {
	enabled bool
}

// Dbg is a variable available to any importing this package to enabling debug printing for error status prints and miscellaneous
var Dbg = NewDebug()

// NewDebug creates a new debug object
func NewDebug() *Debug {
	r := new(Debug)
	r.enabled = true
	return r
}

// Disable turns off debug printing
func (r *Debug) Disable() interface{} {
	if r == nil {
		r = NewDebug()
	}
	r.enabled = false
	return r
}

// Enable turns on debug printing
func (r *Debug) Enable() interface{} {
	if r == nil {
		r = NewDebug()
	}
	r.enabled = true
	return r
}

// Enabled returns true if debug printing is enabled
func (r *Debug) Enabled() bool {
	if r == nil {
		r = NewDebug()
	}
	return r.enabled
}

// Print prints out debugging information
func (r *Debug) Print(s ...interface{}) interface{} {
	if r.enabled {
		red := color.New(color.FgRed).Add(color.Bold)
		green := color.New(color.FgGreen)
		header := color.New(color.FgWhite)
		hyperlink := color.New(color.FgBlue).Add(color.Underline)
		pc, fil, line, _ := runtime.Caller(2)
		var temp string
		for i := range s {
			temp += fmt.Sprint(s[i])
		}
		// temp = fmt.Sprintf(" %-32v", temp)
		header.Print(time.Now().UTC().Format("06-01-02, 15:04.05 ((("))
		red.Print(temp + "))) ")
		S := strings.Split(runtime.FuncForPC(pc).Name(), "gitlab.com/parallelcoin/duo/")
		green.Print(fmt.Sprint(S[1], "() "))
		hyperlink.Print(""+fil, ":", line, " ")
		// pc, fil, line, _ = runtime.Caller(3)
		// hyperlink.Print(fil, ":", line)
		fmt.Println()

	}
	return r
}
