package key

import (
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/proto"
)

// Store is a keychain for public and private keys
type Store struct {
	BC     *blockcrypt.BlockCrypt
	privs  map[ID]*Priv
	pubs   map[ID]*Pub
	Status string
}

// NewStore creates a new Store
func NewStore() *Store {
	return new(Store)
}

// Encrypt sets the store to encrypt private keys
func (r *Store) Encrypt(bc *blockcrypt.BlockCrypt) *Store {
	if r == nil {
		r = NewStore().SetStatus(er.NilRec).(*Store)
	} else {
		for i := range r.privs {
			r.privs[i].BC = bc
			r.privs[i].Copy(r.privs[i].Bytes())
		}
	}
	return r
}

// Decrypt sets the store to not encrypt private keys
func (r *Store) Decrypt() *Store {
	if r == nil {
		r = NewStore().SetStatus(er.NilRec).(*Store)
	} else {
		for i := range r.privs {
			tmp := r.privs[i].Bytes()
			r.privs[i].BC = nil
			r.privs[i].Copy(tmp)
		}
	}
	return r
}

// Find returns the key with matching ID as requested. The return type is Priv but if there is no private key the field is empty
func (r *Store) Find(id *ID) (out *Priv) {
	if r == nil {
		r = NewStore().SetStatus(er.NilRec).(*Store)
	}
	out = new(Priv)
	for i := range r.privs {
		if r.privs[i].IsEqual(id.Bytes()) {
			out = r.privs[i]
			return
		}
	}
	for i := range r.pubs {
		if r.pubs[i].IsEqual(id.Bytes()) {
			out.pub = r.pubs[i]
			return
		}
	}
	return
}

// SetStatus is a
func (r *Store) SetStatus(s string) proto.Status {
	if r == nil {
		r = NewStore().SetStatus(er.NilRec).(*Store)
	} else {
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *Store) SetStatusIf(err error) proto.Status {
	if r == nil {
		r = NewStore().SetStatus(er.NilRec).(*Store)
	} else {
		if err != nil {
			r.Status = err.Error()
		}
	}
	return r
}

// UnsetStatus is a
func (r *Store) UnsetStatus() proto.Status {
	if r == nil {
		r = NewStore().SetStatus(er.NilRec).(*Store)
	} else {
		r.Status = ""
	}
	return r
}

// OK returns true if there is no error
func (r *Store) OK() bool {
	if r == nil {
		r = NewStore().SetStatus(er.NilRec).(*Store)
		return false
	}
	return r.Status == ""
}
