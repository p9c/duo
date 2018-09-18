package blockcrypt

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"testing"
)

func TestBlockCrypt(t *testing.T) {
	bc := New()
	p := []byte("testingpassword123!")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	bc.Generate(pass)
	bc.Arm()
	fmt.Println(p)
	pp := bc.Encrypt(&p)
	fmt.Println(*pp)
	ppp := bc.Decrypt(pp)
	fmt.Println(*ppp)
	if bytes.Compare(p, *ppp) != 0 {
		t.Fatal("Did not correctly encrypt and decrypt")
	}
	IV := bc.IV
	Crypt := bc.Crypt
	Iterations := bc.Iterations
	nc := New()
	nc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), Iterations)
	nc.Unlock(pass)
	nc.Arm()
	fmt.Println(p)
	pp = nc.Encrypt(&p)
	fmt.Println(*pp)
	ppp = nc.Decrypt(pp)
	fmt.Println(*ppp)
	if bytes.Compare(p, *ppp) != 0 {
		t.Fatal("Did not correctly encrypt and decrypt")
	}
	Gen(nil, nil, 0)
	Gen(nc.Password, nil, 0)
	Gen(nc.Password, nc.IV, 0)
	var ec *BlockCrypt
	err := errors.New("")
	ec.SetStatus("test")
	ec.SetStatusIf(err)
	ec.OK()
	ec.UnsetStatus()
	_ = ec.Error()
	nc.OK()
	_ = nc.Error()
	nc.Disarm()
	nc.Arm()
}
