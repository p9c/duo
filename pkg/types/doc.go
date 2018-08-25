// Package types contains implementations for all of the fundamental data structures used in Duo.
/*
Security and ease of use are the key goals of this library. All operations that copy data in erase the source byte-wise before allocating a new buffer, thus they are all labeled 'FromXXX()' to imply that the operation is a move.

All types that have inter-convertible data have functions for all of the relevant types, generally an move operation from another variable or a copy operation that creates a new variable in the relevant type.

All of the types are based on structs in order to simplify error handling (each type stores an error value from the last operation that can produce one) and if one might want to perform a subsequent operation on the same variable, because it returns the pointer to itself, the resultant expression can have further pointer receiver methods chained together.

Each type in the library has distinct *move*, *reference* and *copy* functions. Most type conversions only provide a move function. This potentially could be a side effect for users of the library, but because this pattern is consistent in this library it would be obvious pretty quickly when the side effects rear their ugly heads.

It is also important that consumers of the library be aware of these features because of the use of memguard LockedBuffers which have a very limited supply available. Thus also it is preferable if one function does not use it any more but passes it to another function, it is the same LockedBuffer and the receiver can delete it in confidence that it will not affect other parts of the application.

Of course it is up to the programmer using this library but it is highly recommended that one uses the functional programming principles of clean functions and sharing by communication (which is a Go principle too) are used when using this library because the LockedBuffers make it vitally important that resources are managed properly, and will fail if side effects are not avoided stringently.

The secondary aspect as previously mentioned, is that the reason for using these fenced memory buffers is security, which will of course also benefit from no side effects in the code. Furthermore, if for some reason (using unsafe library, presumably) the application tries to write over the canaries around each buffer, it will also cause a panic and a lot of headache.

*/
package types
