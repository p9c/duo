package main

import (
	"fmt"
)

var (
	lalpha = "abcdefghijklmnopqrstuvwxyz"
	ualpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number = "0123456789"
	symbol = `()[]{}"\/'*.,;:`
	whole  = `package testprototype

	import (
		. "gitlab.com/parallelcoin/duo/pkg/bytes"
	)
	
	type TestPrototype struct {
		val  *[]byte
		set  bool
		utf8 bool
		err  error
	}
	
	type testPrototype interface {
		Buf() []byte
		Copy(*Bytes) *Bytes
		Delete()
		Error() string
		IsSet() bool
		IsUTF8() bool
		Len() int
		Link(*Bytes) *Bytes
		Load(*[]byte) *Bytes
		MarshalJSON() ([]byte, error)
		Move(*Bytes) *Bytes
		New(int) *Bytes
		Null() *Bytes
		Rand(int) *Bytes
		SetBin() *Bytes
		SetError(string) *Bytes
		SetUTF8() *Bytes
		String() string
	}`
)

type Lex struct {
	symtype string
	content string
}

func isWhiteSpace(b byte) bool {
	switch b {
	case ' ', '\n', '\t':
		return true
	}
	return false
}

func isSymbol(b byte) bool {
	for i := range symbol {
		if b == symbol[i] {
			return true
		}
	}
	return false
}

func isAlpha(b byte) bool {
	for i := range lalpha {
		if b == lalpha[i] {
			return true
		}
	}
	for i := range ualpha {
		if b == ualpha[i] {
			return true
		}
	}
	return false
}

func isNum(b byte) bool {
	for i := range number {
		if b == number[i] {
			return true
		}
	}
	return false
}

func main() {
	lexed := [][]byte{}
	// whole, _ := ioutil.ReadAll(os.Stdin)
	i := 0
	for {
		b := []byte{}
		switch {
		case isWhiteSpace(whole[i]):
			i++
		case isAlpha(whole[i]):
			b = append(b, whole[i])
			i++
			for isAlpha(whole[i]) || isNum(whole[i]) {
				b = append(b, whole[i])
				i++
			}
			lexed = append(lexed, b)
		case isNum(whole[i]):
			b = append(b, whole[i])
			i++
			for isNum(whole[i]) || whole[i] == '.' || whole[i] == ',' {
				b = append(b, whole[i])
				i++
			}
			lexed = append(lexed, b)
		case isSymbol(whole[i]):
			lexed = append(lexed, []byte{whole[i]})
			i++
		}
		if i >= len(whole) {
			break
		}
	}
	for i := range lexed {
		if string(lexed[i]) == "package" {
			i++
			fmt.Println("package", string(lexed[i]))
		}
		if string(lexed[i]) == `"` {
			fmt.Println("Open quotes")
			a := lexed[i]
			i++
			for string(lexed[i]) != `"` {
				a = append(a, lexed[i]...)
				i++
			}
			a = append(a, lexed[i]...)
		}
	}
}
