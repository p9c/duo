package tac

import (
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/key"
)

// A URL is a specification of a destination or local network binding
type URL string

// Func is the type for any TaC API call
type Func func(string, ...interface{}) []interface{}

// Function is a specification of the identifying string for a function call, and the function itself, that can be invoked by an authorised client, or, passed directly as a function call instead of marshalled to a wire format. Name format allows 3 levels, package/module/function. They are arbitrary strings, but should be unique and meaningful.
type Function struct {
	Name string
	Call *Func
}

// Peer is the collection of inbound connection endpoints and public key for contacting a peer
type Peer struct {
	PubKey    []byte
	Addresses []URL
}

// Session is an open connection to a peer
type Session struct {
	core.State
	Endpoint    URL
	Peer        Peer
	Key         buf.Secure
	Established int64
	Expires     int64
	Handle      interface{}
}

// Node is an endpoint used to control an application
type Node struct {
	core.State
	PubKey     key.Pub
	PrivKey    key.Priv
	Entrypoint *func() error
	Listeners  []URL
	Sessions   []Session
	Functions  []Function
	Configs    []Configuration
}

// Configuration is a key and value pair that represents a configuration item for a node
type Configuration struct {
	Key   string
	Value string
}
