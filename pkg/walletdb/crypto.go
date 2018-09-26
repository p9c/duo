package walletdb

import (
	"github.com/parallelcointeam/duo/pkg/buf"
)

// Encrypt transparently uses a BlockCrypt if available to encrypt the data before it is written to the database, or it writes plaintext
func (r *DB) Encrypt(in *buf.Secure) (out *buf.Byte) {
	r = r.NewIf()
	switch {
	case !r.OK():
		return &buf.Byte{}
	case r.BC != nil:
		out = out.Copy(r.BC.Encrypt(in.Bytes())).(*buf.Byte)
	default:
		out = out.Copy(in.Bytes()).(*buf.Byte)
	}
	return
}

// Decrypt transparently uses a BlockCrypt if available to deecrypt the data before it is returned to the caller, or it writes plaintext
func (r *DB) Decrypt(in *buf.Byte) (out *buf.Secure) {
	r = r.NewIf()
	switch {
	case !r.OK():
		return &buf.Secure{}
	case r.BC.GCM != nil:
		out = out.Copy(r.BC.Decrypt(in.Bytes())).(*buf.Secure)
	default:
		out = out.Copy(in.Bytes()).(*buf.Secure)
	}
	return
}
