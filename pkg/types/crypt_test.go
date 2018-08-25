package types

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"
)

func TestCrypt(t *testing.T) {
	a := NewCrypt()
	p := "testpassword"
	a.Generate(NewPassword().FromString(&p))
}

func TestKDF(t *testing.T) {
	p := "testpassword"
	password := NewPassword().FromString(&p)
	iv := NewBytes().WithSize(12)
	rand.Read(*iv.Buffer())
	iterations := 400000
	tstart := time.Now()
	LB, IV, err := kdf(password, iv, iterations)
	tend := time.Now()
	elapsed := tend.Sub(tstart)
	fmt.Println(elapsed, "elapsed", elapsed.Nanoseconds()/int64(iterations), "nanoseconds/iteration")
	fmt.Println(LB, IV, err)
}

func TestKDFBench(t *testing.T) {
	iter := KDFBench(time.Second)
	fmt.Println("Performed", iter, "iterations in 1 seecond")
}
