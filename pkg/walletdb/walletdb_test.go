package walletdb

import (
	"encoding/hex"
	"fmt"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/proto"
	"testing"
)

func TestMasterKey(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	teststring := "if you can read this the encryption worked"
	testbytes := []byte(teststring)
	testcipher := BC.Encrypt(&testbytes)

	fmt.Println("\nMESSAGE ENCRYPTION")
	fmt.Println("string '" + teststring + "'")
	fmt.Println("plain", len(testbytes), hex.EncodeToString(testbytes))
	fmt.Println("ciph ", len(*testcipher), hex.EncodeToString(*testcipher))

	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.WriteMasterKey(BC)
	crypt := BC.Crypt.Bytes()
	idx := proto.Hash64(crypt)
	defer wdb.EraseMasterKey(idx)
	BCs := wdb.ReadMasterKeys()
	for i := range BCs {
		crypt = BCs[i].Crypt.Bytes()
		iv := BCs[i].IV.Bytes()
		iterations := BCs[i].Iterations
		idx = proto.Hash64(crypt)

		fmt.Println("\nMASTER KEY")
		fmt.Println("idx  ", len(*idx), hex.EncodeToString(*idx))
		fmt.Println("crypt", len(*crypt), hex.EncodeToString(*crypt))
		fmt.Println("iv   ", len(*iv), hex.EncodeToString(*iv))
		fmt.Println("iters", iterations)
		BCs[i].Unlock(buf.NewSecure().Copy(&p).(*buf.Secure)).Arm()

		fmt.Println("\nMESSAGE DECRYPTION")
		fmt.Println("ciph ", len(*testcipher), hex.EncodeToString(*testcipher))
		fmt.Println("decrp", len(*BCs[i].Decrypt(testcipher)), hex.EncodeToString(*BCs[i].Decrypt(testcipher)))
		fmt.Println("decrs '" + string(*BCs[i].Decrypt(testcipher)) + "'")

		k := key.NewPriv()
		k.WithBC(BCs[i])
		k.Make()
		kh := []byte(k.GetID())
		pidx := proto.Hash64(&kh)

		fmt.Println("\nORIGINAL KEY")
		fmt.Println("idx   ", len(*pidx), hex.EncodeToString(*pidx))
		fmt.Println("id    ", len(kh), hex.EncodeToString(kh))
		fmt.Println("prvkey", len(*k.Bytes()), hex.EncodeToString(*k.Bytes()))
		fmt.Println("crypt ", len(*k.Crypt.Val), hex.EncodeToString(*k.Crypt.Val))
		fmt.Println("pubkey", len(*k.PubKey().Bytes()), hex.EncodeToString(*k.PubKey().Bytes()))
		wdb.WithBC(BCs[i])
		wdb.WriteKey(k)
		defer wdb.EraseKey(&kh)

		fmt.Println("\nRECOVERED KEY")
		pk := wdb.ReadKey(&kh)
		fmt.Println("prvkey", len(*pk.Bytes()), hex.EncodeToString(*pk.Bytes()))
		fmt.Println("pubkey", pk.PubKey().(*buf.Byte).Len(), hex.EncodeToString(*pk.PubKey().Bytes()))

		fmt.Println("\nWRITE NAME")
		k.WithBC(BCs[i])
		k.Make()
		id := []byte(k.GetID())
		label := []byte("some random thing")
		wdb.WriteName(&id, &label)
		fmt.Println("addr  ", hex.EncodeToString(id))
		fmt.Println("label ", string(label))
		defer wdb.EraseName(&id)

		fmt.Println("\nREAD NAME")
		rName := wdb.ReadName(&id)
		fmt.Println("addr  ", hex.EncodeToString(rName.Address))
		fmt.Println("label ", rName.Label)

		fmt.Println("\nWRITE ACCOUNT")
		k.WithBC(BCs[i])
		k.Make()
		address := []byte(k.GetID())
		pub := k.PubKey().Bytes()
		wdb.WriteAccount(&address, pub)
		fmt.Println("addr  ", hex.EncodeToString(address))
		fmt.Println("pub   ", hex.EncodeToString(*pub))
		defer wdb.EraseAccount(&address)

		fmt.Println("\nREAD ACCOUNT")
		rAccount := wdb.ReadAccount(&address)
		fmt.Println("addr  ", hex.EncodeToString(rAccount.Address))
		fmt.Println("pub   ", hex.EncodeToString(rAccount.Pub))

		fmt.Println("\nREMOVE BLOCKCRYPT")
		wdb.RemoveBC()

		fmt.Println("\nRECOVERED KEY")
		pk = wdb.ReadKey(&kh)
		fmt.Println("prvkey", len(*pk.Bytes()), hex.EncodeToString(*pk.Bytes()))
		fmt.Println("pubkey", pk.PubKey().(*buf.Byte).Len(), hex.EncodeToString(*pk.PubKey().Bytes()))

		fmt.Println("\nREAD NAME")
		rName = wdb.ReadName(&id)
		fmt.Println("addr  ", hex.EncodeToString(rName.Address))
		fmt.Println("label ", rName.Label)

		fmt.Println("\nREAD ACCOUNT")
		rAccount = wdb.ReadAccount(&address)
		fmt.Println("addr  ", hex.EncodeToString(rAccount.Address))
		fmt.Println("pub   ", hex.EncodeToString(rAccount.Pub))

		wdb.dump()
	}
}
