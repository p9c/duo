package walletdb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/proto"
	"github.com/parallelcointeam/duo/pkg/walletdb/entries"
)

// ReadName reads a name entry out of the database
func (r *DB) ReadName(id *[]byte) (out *rec.Name) {
	out = new(rec.Name)
	k := []byte(rec.Tables["Name"])
	idx := proto.Hash64(id)
	k = append(k, *idx...)
	if r.BC != nil {
		id = r.BC.Encrypt(id)
	}
	k = append(k, *id...)
	fmt.Println("search", hex.EncodeToString(k))
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	var V []byte
	err := r.DB.View(func(txn *badger.Txn) error {
		item, er := txn.Get(k)
		if er != nil {
			return er
		}
		V, er = item.Value()
		if er != nil {
			return er
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
		if r.BC != nil {
			out.Address = *r.BC.Decrypt(id)
			out.Label = string(*r.BC.Decrypt(&V))
		} else {
			out.Address = *id
			out.Label = string(V)
		}
		out.Idx = *idx
	}
	return
}

// WriteName writes a name entry to the database
func (r *DB) WriteName(address, label *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if address == nil || label == nil {
		r.SetStatus(er.NilParam)
	}
	idx := proto.Hash64(address)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
		label = r.BC.Encrypt(label)
	}
	k := []byte(rec.Tables["Name"])
	k = append(k, *idx...)
	k = append(k, *address...)
	fmt.Println("write  ", hex.EncodeToString(k))
	v := *label
	txn := r.DB.NewTransaction(true)
	err := txn.Set(k, v)
	if r.SetStatusIf(err).OK() {
		r.SetStatusIf(txn.Commit(nil))
	}
	return r
}

// EraseName removes a name entry from the database
func (r *DB) EraseName(address *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	idx := proto.Hash64(address)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
	}
	k := []byte(rec.Tables["Name"])
	k = append(k, *idx...)
	k = append(k, *address...)
	txn := r.DB.NewTransaction(true)
	if r.SetStatusIf(txn.Delete(k)).OK() {
		txn.Commit(nil)
	}
	return r
}

// ReadTx reads a transaction entry from the database
func (r *DB) ReadTx() {}

// WriteTx writes a transaction entry from the database
func (r *DB) WriteTx() {}

// EraseTx deletes a transaction entry from the database
func (r *DB) EraseTx() {}

// ReadKey reads a key entry from the database
func (r *DB) ReadKey(id *[]byte) (out *key.Priv) {
	k := []byte(rec.Tables["Key"])
	idx := proto.Hash64(id)
	k = append(k, *idx...)
	if r.BC != nil {
		id = r.BC.Encrypt(id)
	}
	k = append(k, *id...)
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	var V []byte
	err := r.DB.View(func(txn *badger.Txn) error {
		item, er := txn.Get(k)
		if er != nil {
			return er
		}
		V, er = item.Value()
		if er != nil {
			return er
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
		var priv, pub *[]byte
		if r.BC != nil {
			encpriv := V[:48]
			encpub := V[48:]
			out = key.NewPriv()
			out.WithBC(r.BC)
			priv = r.BC.Decrypt(&encpriv)
			pub = r.BC.Decrypt(&encpub)
		} else {
			pr, pu := V[:32], V[32:]
			priv, pub = &pr, &pu
		}
		out.SetKey(priv, pub)
	}
	return
}

// WriteKey writes a key entry to the database
func (r *DB) WriteKey(priv *key.Priv) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if priv.Crypt.Len() < 1 {
		r.SetStatus("zero length crypt")
		return r
	}
	I := []byte(priv.GetID())
	id := &I
	idx := proto.Hash64(id)
	var pk, pp *[]byte
	if r.BC != nil {
		id = priv.BC.Encrypt(id)
		pk = priv.BC.Encrypt(priv.Bytes())
		pp = priv.BC.Encrypt(priv.PubKey().Bytes())
	} else {
		pk = priv.Bytes()
		pp = priv.PubKey().Bytes()
	}
	k := []byte(rec.Tables["Key"])
	k = append(k, *idx...)
	k = append(k, *id...)
	v := *pk
	v = append(v, *pp...)
	txn := r.DB.NewTransaction(true)
	err := txn.Set(k, v)
	if r.SetStatusIf(err).OK() {
		r.SetStatusIf(txn.Commit(nil))
	}
	return r
}

// EraseKey deletes a key entry
func (r *DB) EraseKey(id *[]byte) *DB {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	idx := proto.Hash64(id)
	search := append(rec.Tables["Key"], *idx...)
	ID := buf.NewSecure().Copy(id).(*buf.Secure)
	encid := r.Encrypt(ID)
	search = append(search, *encid.Bytes()...)
	txn := r.DB.NewTransaction(true)
	if r.SetStatusIf(txn.Delete(search)).OK() {
		txn.Commit(nil)
	}
	return r
}

// ReadMasterKeys returns all of the masterkey entries in the database
func (r *DB) ReadMasterKeys() (BC []*bc.BlockCrypt) {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := item.Key()
			value, _ := item.Value()
			table := key[:8]
			if bytes.Compare(table, rec.Tables["MasterKey"]) == 0 {
				idx := key[8:16]
				crypt := value[:48]
				cryptHash := proto.Hash64(&crypt)
				if bytes.Compare(idx, *cryptHash) != 0 {
					r.SetStatus("index of crypt was incorrect")
				}
				iv := value[48:60]
				iterations := value[60:68]
				var it int64
				proto.BytesToInt(&it, &iterations)
				BC = append(BC,
					&bc.BlockCrypt{
						Crypt:      buf.NewByte().Copy(&crypt).(*buf.Byte),
						IV:         buf.NewByte().Copy(&iv).(*buf.Byte),
						Iterations: it,
					})
			}
		}
		return nil
	})
	r.SetStatusIf(err)
	return
}

// WriteMasterKey adds a master key entry to the database
func (r *DB) WriteMasterKey(BC *bc.BlockCrypt) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if BC.Crypt.Len() < 1 {
		r.SetStatus("zero length crypt")
		return r
	}
	out := proto.Hash64(BC.Crypt.Bytes())
	key := append(rec.Tables["MasterKey"], *out...)
	value := *BC.Crypt.Val
	value = append(value, *BC.IV.Bytes()...)
	value = append(value, *proto.IntToBytes(BC.Iterations)...)
	txn := r.DB.NewTransaction(true)
	err := txn.Set(key, value)
	r.SetStatusIf(err)
	txn.Commit(nil)
	return r
}

// EraseMasterKey deletes a masterkey entry from the database
func (r *DB) EraseMasterKey(idx *[]byte) *DB {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	search := append(rec.Tables["MasterKey"], *idx...)
	txn := r.DB.NewTransaction(true)
	if !r.SetStatusIf(txn.Delete(search)).OK() {
	}
	txn.Commit(nil)
	return r
}

// WriteScript writes a script entry to the database
func (r *DB) WriteScript() {}

// EraseScript deletes a script entry from the database
func (r *DB) EraseScript() {}

// WriteDefaultKey updates the default key used by interfaces when receiving payments
func (r *DB) WriteDefaultKey() {

}

// ReadDefaultKey returns the current set default key
func (r *DB) ReadDefaultKey() {

}

// WriteBestBlock gets the current best block entry
func (r *DB) WriteBestBlock() {}

// ReadBestBlock gets the current best block entry
func (r *DB) ReadBestBlock() {}

// ReadPool gets the oldest available pool entry and refreshes the pool after addresses are used
func (r *DB) ReadPool() {

}

// WritePool adds a new pool key to the wallet
func (r *DB) WritePool() {

}

// ErasePool removes a pool key
func (r *DB) ErasePool() {

}

// ReadMinVersion returns the minimum version required to read this database
func (r *DB) ReadMinVersion() {

}

// WriteMinVersion updates the minimum version
func (r *DB) WriteMinVersion() {

}

// ReadAccount finds an account stored due to being a correspondent account
func (r *DB) ReadAccount(address *[]byte) (out *rec.Account) {
	out = new(rec.Account)
	k := []byte(rec.Tables["Account"])
	idx := proto.Hash64(address)
	k = append(k, *idx...)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
	}
	k = append(k, *address...)
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	var V []byte
	err := r.DB.View(func(txn *badger.Txn) error {
		item, er := txn.Get(k)
		if er != nil {
			return er
		}
		V, er = item.Value()
		if er != nil {
			return er
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
		out.Idx = *idx
		if r.BC != nil {
			out.Address = *r.BC.Decrypt(address)
			if len(V) > 1 {
				out.Pub = *r.BC.Decrypt(&V)
			}
		} else {
			out.Address = *address
			out.Pub = V
		}
	}
	return
}

// WriteAccount writes a new account entry
func (r *DB) WriteAccount(address, pub *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if address == nil {
		r.SetStatus(er.NilParam)
		return r
	}
	idx := proto.Hash64(address)
	var k, v []byte
	if r.BC != nil {
		address = r.BC.Encrypt(address)
		if pub != nil {
			pub = r.BC.Encrypt(pub)
		}
	}
	k = []byte(rec.Tables["Account"])
	k = append(k, *idx...)
	k = append(k, *address...)
	if pub != nil {
		v = *pub
	} else {
		v = []byte{}
	}
	txn := r.DB.NewTransaction(true)
	err := txn.Set(k, v)
	if r.SetStatusIf(err).OK() {
		r.SetStatusIf(txn.Commit(nil))
	}
	return r
}

// EraseAccount deletes an account from the wallet database
func (r *DB) EraseAccount(address *[]byte) *DB {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	idx := proto.Hash64(address)
	search := append(rec.Tables["Account"], *idx...)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
	}
	search = append(search, *address...)
	txn := r.DB.NewTransaction(true)
	if r.SetStatusIf(txn.Delete(search)).OK() {
		txn.Commit(nil)
	}
	return r

}

// ReadAccountingEntry writes an accounting entry based on a transaction
func (r *DB) ReadAccountingEntry() {}

// WriteAccountingEntry writes an accounting entry based on a transaction
func (r *DB) WriteAccountingEntry() {}

// EraseAccountingEntry writes an accounting entry based on a transaction
func (r *DB) EraseAccountingEntry() {}

// GetAccountCreditDebit finds entries in the credit/debit records written related to each input transaction from a list of indexes of accounts of interest
func (r *DB) GetAccountCreditDebit() {}
