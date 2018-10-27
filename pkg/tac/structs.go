package tac

// An Address is a specification of a destination or local network binding
type Address struct {
	Network string
	Address string
}

// Function is a specification of the identifying string for a function call, and the function itself, that can be invoked by an authorised client, or, if Local is true, passed directly as a function call instead of marshalled to a wire format. Name format allows 3 levels, package/module/function
type Function struct {
	Name     string
	Function *func(...interface{}) []interface{}
	Local    bool
}

// Session is an open connection to a peer
type Session struct {
	Peer        *[]Address
	Key         []byte
	Established int64
	Expires     int64
}

// Node is an endpoint used to control an application
type Node struct {
	Listeners []Address
	Sessions  []Session
	Functions []Function
}

// Configuration is a key and value pair that represents a configuration item for a node
type Configuration struct {
	Key   string
	Value string
}
