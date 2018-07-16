package key

import (
	"errors"
	"sync"
)

// Pair is a public/private key pair
type Pair struct {
	Public  *Pub
	Private *Priv
}

// Map of keys
type Map map[*ID]*Pair

// ScriptMap stores scripts
type ScriptMap map[*ScriptID]*Script

// Store is a basic key and script store
type Store struct {
	Mutex   sync.Mutex
	Keys    Map
	Scripts ScriptMap
}

type store interface {
	AddKeyPair(*Priv, *Pub) error
	AddPrivKey(*Priv) error
	HaveKey(*ID) bool
	GetPriv(*ID) (*Priv, error)
	GetPrivs([]*ID) ([]Priv, error)
	GetPub(*ID) *Pub
	AddScript(*Script) error
	HaveScript(*ScriptID) bool
	GetScript(*ScriptID) (*Script, error)
}

// NewStore creates a new basic key.Store
func NewStore() *Store {
	return &Store{}
}

// AddKeyPair adds a public key to a key.Store
func (s *Store) AddKeyPair(priv *Priv, pub *Pub) (err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	id := pub.GetID()
	s.Keys[id] = &Pair{pub, priv}
	return
}

// AddPrivKey adds a key
func (s *Store) AddPrivKey(priv *Priv) (err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.Keys[priv.GetPub().GetID()]; !ok {
		s.Keys[priv.GetPub().GetID()] = &Pair{priv.GetPub(), priv}
	} else {
		return errors.New("Key is already in store")
	}
	return
}

// HaveKey returns true if the key with the id is in the store
func (s *Store) HaveKey(id *ID) (r bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.Keys[id]; ok {
		r = true
	}
	return
}

// GetPriv gets the key matching an id
func (s *Store) GetPriv(id *ID) (priv *Priv, err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if key, ok := s.Keys[id]; ok {
		return key.Private, nil
	}
	return nil, errors.New("Key ID not found in key store")
}

// GetPrivs gets a set of keys matching the id's in a request
func (s *Store) GetPrivs(ids []*ID) (privs []Priv, err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	for i := range ids {
		if key, ok := s.Keys[ids[i]]; ok {
			privs = append(privs, *key.Private)
		}
	}
	return
}

// GetPub gets only the public key from a key with a specified ID
func (s *Store) GetPub(id *ID) (pub *Pub) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.Keys[id]; ok {
		pub.pub = s.Keys[id].Public.pub
	}
	return
}

// AddScript adds a script to the key.Store
func (s *Store) AddScript(script *Script) (err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Scripts[script.GetID()] = script
	return
}

// HaveScript returns whether the script with the ID is present in the store
func (s *Store) HaveScript(id *ScriptID) (r bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.Scripts[id]; ok {
		r = true
	}
	return
}

// GetScript retrieves a script given an ID
func (s *Store) GetScript(id *ScriptID) (script *Script) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.Scripts[id]; ok {
		script = s.Scripts[id]
	}
	return
}
