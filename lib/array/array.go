package array

// Array is an interface that allows a type to present access to uniform sized parts of the parent structure
type Array interface {
	Elem(int) interface{}
	Len() int
	SetElem(int, interface{}) Array
}
