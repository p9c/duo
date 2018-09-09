package buf

import (
	"fmt"
	"runtime"
	"testing"

	"gitlab.com/parallelcoin/duo/lib/debug"
)

func TestFreezeThaw(t *testing.T) {

	defer func() {
		if recover() != nil {
			fmt.Println("wooty")
			pc, fil, line, _ := runtime.Caller(0)
			fmt.Println(pc, fil, line)
			dbg.D.Close()
		}
	}()
	b := NewByte()
	frozen := b.Freeze()
	dbg.Append(frozen)
	dbg.D.Close()
}
