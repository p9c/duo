package types

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"
)

func TestCrypt(t *testing.T) {

}

func TestKDF(t *testing.T) {
	p := "testpassword"
	password := NewPassword().FromString(&p)
	iv := NewBytes().WithSize(16)
	rand.Read(*iv.Buffer())
	iterations := 400000
	tstart := time.Now()
	LB, IV, err := KDF(password, iv, iterations)
	tend := time.Now()
	elapsed := tend.Sub(tstart)
	fmt.Println(elapsed, "elapsed", elapsed.Nanoseconds()/int64(iterations), "nanoseconds/iteration")
	fmt.Println(LB, IV, err)
}
