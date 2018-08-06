package logger
import (
	"fmt"
	"os"
	"time"
)
var (
	// PrintFlag indicates whether we are printing to the tty
	PrintFlag = true
	// DebugFlag indicates whether we will print debug messages
	DebugFlag = false
	// DebugNetFlag indicates whether we will print network debug messages
	DebugNetFlag = false
	// Output is the file the output will write to, by default it is the current tty
	Output = os.Stdout
)
func Debug(input ...interface{}) {
	if PrintFlag {
		fmt.Fprintln(Output, []string{time.Now().UTC().Format("2006-01-02 15:04:05 UTC")}, input)
	}
}
func DebugNet(input ...interface{}) {
	if DebugNetFlag {
		Debug("NET: ", input)
	}
}
