package main

import (
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
)

func main() {
	a := buf.NewUnsafe()
	fmt.Println(a.Buf())
}
