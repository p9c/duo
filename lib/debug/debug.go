package dbg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
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

// Note is a debug-only structure used for miscellaneous log notes
type Note struct {
	Value string
}

// NewNote creates a new note with a given string
func NewNote(s string) *Note {
	return &Note{s}
}

// Freeze renders the note into JSON
func (r *Note) Freeze() (R string) {
	return `"Note":"` + r.Value + `"`
}

// Thaw unmarshals the note from JSON
func (r *Note) Thaw(string) (R interface{}) {
	return
}

// Dbg is a debug instance
type Dbg struct {
	File    *os.File
	Counter int
	Last    int64
}

var (
	// LogFileName is the location where the log is output for viewing
	LogFileName = "/tmp/duo.json"
	// D is the central repository for the debugger
	D       Dbg
	initlen int64
)

// Append an entry to the log
func Append(key string, object Debugger) (out string) {
	out = ","
	out += `"` + fmt.Sprint(time.Now().UTC().Format("06-01-02 15:04:05.0000000")) + "\":{"
	pc, fil, line, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc).Name()
	out += `"line":"file://` + fil + ":" + fmt.Sprint(line) + "\","
	out += `"func":` + strconv.Quote(fun) + `,`
	fileHandle, _ := os.Open(fil)
	fileScanner := bufio.NewScanner(fileHandle)
	i := 1
	for fileScanner.Scan() {
		if i == line {
			out += `"code":` + strconv.Quote(fileScanner.Text()) + `,`
		}
		i++
	}
	fileHandle.Close()
	if key != "" {
		out += `"` + key + `":`
	}
	out += object.Freeze() + `}`
	D.File.Write([]byte(out))
	D.Counter++
	return out
}

// Close the json string. Should be run as part of a defer recover() closure in main
func (r *Dbg) Close() {
	out := "}}"
	D.File.Write([]byte(out))
	if err := D.File.Close(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Init()
}

// Init is exposed here for testing purposes
func Init() {
	var err error
	if _, err = os.Stat(LogFileName); os.IsNotExist(err) {
		D.File, err = os.OpenFile(LogFileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println(err)
		}
		D.File.Write([]byte("{"))
	} else {
		D.File, err = os.OpenFile(LogFileName, os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			fmt.Println(err)
		}
		fi, err := D.File.Stat()
		if err != nil {
		}
		initlen = fi.Size()
		switch {
		case initlen == 0:
			D.File.Write([]byte("{"))
		case initlen > 0:
			offset := fi.Size() - 1
			D.File.Truncate(offset)
			D.File.Write([]byte(","))
		}
	}
	D.File.Write([]byte("\"" + time.Now().UTC().Format("06-01-02 15:04:05") + "\":{"))
	s := strconv.Quote(strings.Join(os.Args, " "))
	s = s[1 : len(s)-2]
	D.File.Write([]byte(`"Args":"` + s + `"`))
	fi, err := D.File.Stat()
	if err != nil {
	}
	initlen = fi.Size() + 1
}
