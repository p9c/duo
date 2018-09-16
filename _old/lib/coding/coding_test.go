package coding_test

import (
	"bytes"
	"crypto/rand"
	"gitlab.com/parallelcoin/duo/lib/coding"
	"testing"
)

func TestCoding(t *testing.T) {
	wantB := []byte{113}
	wantW := []byte("this is an average size sentence.")
	wantL := make([]byte, 256)
	rand.Read(wantL)
	for i := range coding.Codings {
		if gotB := coding.Decode(coding.Encode(wantB, coding.Codings[i]), coding.Codings[i]); bytes.Compare(gotB, wantB) != 0 {
			t.Errorf("Encode/Decode for %s testfailed, got %v, want %v", coding.Codings[i], gotB, wantB)
		}
		if gotW := coding.Decode(coding.Encode(wantW, coding.Codings[i]), coding.Codings[i]); bytes.Compare(gotW, wantW) != 0 {
			t.Errorf("Encode/Decode for %s testfailed, got %v, want %v", coding.Codings[i], gotW, wantW)
		}
		if gotL := coding.Decode(coding.Encode(wantL, coding.Codings[i]), coding.Codings[i]); bytes.Compare(gotL, wantL) != 0 {
			t.Errorf("Encode/Decode for %s testfailed, got %v, want %v", coding.Codings[i], gotL, wantL)
		}
	}
	coding.Decode("", "hex")
	coding.Decode("abcdefg", "base58check")
}
