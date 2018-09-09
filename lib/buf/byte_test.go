package buf

import (
	"fmt"
	"gitlab.com/parallelcoin/duo/lib/debug"
	"runtime"
	"testing"
)

func TestFreezeThaw(t *testing.T) {
	defer func() {
		if recover() != nil {
			pc, fil, line, _ := runtime.Caller(0)
			fmt.Println(pc, fil, line)
			dbg.D.Close()
		}
	}()
	b := NewByte()
	frozen := b.Freeze()
	dbg.Append(`"parameters":""`, frozen)
	dbg.D.Close()
	// o7
}
