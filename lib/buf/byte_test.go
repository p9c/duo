package buf

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"gitlab.com/parallelcoin/duo/lib/debug"
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

	dbg.Append(frozen)
	dbg.D.Close()

}

func TestByte_Null(t *testing.T) {
	type fields struct {
		buf    byte
		status string
		coding string
	}
	tests := []struct {
		name   string
		fields fields
		want   Buf
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Byte{
				buf:    tt.fields.buf,
				status: tt.fields.status,
				coding: tt.fields.coding,
			}
			if got := r.Null(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Byte.Null() = %v, want %v", got, tt.want)
			}
		})
	}
}
