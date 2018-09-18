package blockcrypt

import (
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
}
