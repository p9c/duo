package key

import (
	"crypto/sha256"
	"fmt"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"testing"
)

func TestCrypt(t *testing.T) {
	fmt.Println()
	pass := []byte("password")
	BC := blockcrypt.New().Generate(buf.NewSecure().Copy(&pass).(*buf.Secure)).Arm()
	priv := NewPriv().WithBC(BC)
	priv.Make().SetCoding("string")
	fmt.Println("key as plaintext", *priv.Bytes())
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
	// fmt.Println(priv.Copy(priv.Bytes()))
	out := priv.Bytes()
	fmt.Println("key as plaintext", *out)
	fmt.Println("content of buffer", *priv.Crypt.Bytes())
	emptypriv.Sign(nil)
	emptypriv.SignCompact(nil)
	message := []byte("Testing signatures")
	messageHash := sha256.Sum256(message)
	mh := messageHash[:]
	fmt.Println(priv.Sign(&mh))
	fmt.Println(priv.SignCompact(&mh))
}
