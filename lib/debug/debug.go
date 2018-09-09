package dbg

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// Debugger is an interface for keeping track of state changes for debugging
type Debugger interface {
	// Freeze prints the full content of the variable into JSON format
	Freeze() string
	// Thaw takes a description made by the above and fills a variable with it
	Thaw(string) interface{}
}

type Dbg struct {
	File    *os.File
	Counter int
	Last    int64
}

var (
	// LogFileName is the location where the log is output for viewing
	LogFileName = "/tmp/duo.json"
	// D is the central repository for the debugger
	D Dbg
)

// Append an entry to the log
func Append(entry ...string) string {
	pc, fil, line, _ := runtime.Caller(0)
	fun := runtime.FuncForPC(pc).Name()
	out := "\"datestamp\":\"" + fmt.Sprint(time.Now().UTC().Format("06-01-02 15:04:05.00000")) + "\""
	out += ",\"func\":\"" + fun + "\","
	out += "\"line\":\"" + fil + ":" + fmt.Sprint(line) + "\","
	out += strings.Join(entry, ",")
	// out += "}"
	D.File.Write([]byte(out))
	return out
}

// Close the json string. Should be run as part of a defer recover() closure in main
func (r *Dbg) Close() {
	out := "}"
	D.File.Write([]byte(out))
	if err := D.File.Close(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error
	if _, err = os.Stat(LogFileName); os.IsNotExist(err) {
		D.File, err = os.OpenFile(LogFileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		D.File.Write([]byte(""))
	} else {
		D.File, err = os.OpenFile(LogFileName, os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		offset, _ := D.File.Seek(-1, os.SEEK_END)
		D.File.Truncate(offset)
		D.File.Write([]byte(","))
	}
	D.File.Write([]byte("\"" + time.Now().UTC().Format("06-01-02 15:04:05") + "\":{"))
	D.File.Write([]byte("\"Args\":\"" + strings.Join(os.Args, " ") + "\","))
}
