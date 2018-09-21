package key

import (
	"crypto/sha256"
	"fmt"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/hash160"
	"testing"
)

func TestCrypt(t *testing.T) {
	fmt.Println()
	pass := []byte("password")
	BC := blockcrypt.New().Generate(buf.NewSecure().Copy(&pass).(*buf.Secure)).Arm()
	priv := NewPriv().WithBC(BC)
	priv.Make().SetCoding("string")
	priv.pub.Decompress()
	address := hash160.Sum(priv.pub.Bytes())
	origPub := priv.PubKey().Bytes()
	fmt.Println("public key")
	fmt.Println("private key as plaintext", *priv.Bytes())
	fmt.Println("content of buffer", *priv.Crypt.Bytes())
	var emptypriv *Priv
	emptypriv.WithBC(BC)
	fmt.Println(emptypriv.IsValid(), priv.IsValid())
	fmt.Println(emptypriv.Bytes())

	fmt.Println(emptypriv.Invalidate().Bytes())
	fmt.Println(emptypriv.Zero())
	fmt.Println(emptypriv.Make())
	fmt.Println(emptypriv.AsEC())
	fmt.Println(priv.AsEC())
	fmt.Println(emptypriv.PubKey())
	fmt.Println(priv.PubKey())
	privB, pubB := priv.Bytes(), priv.pub.Bytes()
	fmt.Println(privB, pubB)
	fmt.Println(priv.Invalidate())
	fmt.Println(priv.SetKey(privB, pubB))
	out := priv.Bytes()
	fmt.Println("key as plaintext", *out)
	fmt.Println("content of buffer", *priv.Crypt.Bytes())
	emptypriv.Sign(nil)
	emptypriv.SignCompact(nil)
	message := []byte("Testing signatures")
	messageHash := sha256.Sum256(message)
	mh := messageHash[:]
	full := priv.Sign(&mh)
	compact := priv.SignCompact(&mh)
	fmt.Println("full signature", full.Len(), full.Bytes())
	fmt.Println("compact signature", compact.Len(), compact.Bytes())
	var emptysig *Sig
	fmt.Println(emptysig.AsEC())
	fmt.Println(full.AsEC())
	fmt.Println(full.Error())
	fmt.Println(compact.AsEC())
	fmt.Println(compact.Error())
	fmt.Println("original public key ", origPub)
	c := compact.Recover(&mh, address)
	fmt.Println("From compact signature\nrecovered public key", c.Bytes())
	f := full.Recover(&mh, address)
	fmt.Println("From full signature\nrecovered public key", f.Bytes())
}
