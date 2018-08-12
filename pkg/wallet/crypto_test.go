package wallet

import (
	"bytes"
	"crypto/cipher"
	"crypto/aes"
	"fmt"
	"testing"
)
func TestCrypto(t *testing.T) {
   db, err := NewDB(f)
	if err != nil {
		t.Error("Failed to open")
	}
	db.Net = "mainnet"
	for i := 0; i<Flast; i++ {
		db.NewTable(KeyNames[i])
	}
	imp, err := Import()
	if err != nil {
		t.Error("failed to import wallet", err)
   }
   fmt.Println("encrypted masterkey: ", imp.MKeys[0].EncryptedKey)
   ckey, iv, err := imp.MKeys[0].DeriveCipher(passwd)
   fmt.Println("Master key ciphertext", ckey.Buffer())
   block, err := aes.NewCipher(ckey.Buffer()[:32])
   dec := cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
   enc := cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])
   // mcipher, err := memguard.NewMutable(64)
   // block.Decrypt(mcipher.Buffer(), append(imp.MKeys[0].EncryptedKey, imp.MKeys[0].Salt...))
   // fmt.Println("decrypted masterkey: ", mcipher.Buffer())
   // block, err = aes.NewCipher(mcipher.Buffer())
   for i := range imp.CKeys {
      fmt.Println("encrypted     ", imp.CKeys[i].Priv)
      rePriv := make([]byte, len(imp.CKeys[i].Priv))
      dPriv := make([]byte, len(imp.CKeys[i].Priv))
      // mode.CryptBlocks(dPub, imp.CKeys[i].Pub)
      dec.CryptBlocks(dPriv, imp.CKeys[i].Priv)
      enc.CryptBlocks(rePriv, dPriv)
      fmt.Println("re-encrypted  ", imp.CKeys[i].Priv)
      if bytes.Compare(imp.CKeys[i].Priv, rePriv) != 0 {
         fmt.Println("encrypted decrypted wallet key did not match ", dPriv, rePriv)
      }
   }
}