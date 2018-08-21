// A library for keeping a cache based on a most recently used set
package mruset
// MRUSet - implement this for each type of set with a composition
type MRUSet struct {
	// MaxSize is the maximum size of an Most Recently Used list
	MaxSize int
	// set of items (map of some type of item
	// queue of items (sorted map)
}
