// Package types contains implementations for all of the fundamental data structures used in Duo.
/*
Security and ease of use are the key goals of this library. All operations that copy data in erase the source byte-wise before allocating a new buffer, thus they are all labeled 'FromXXX()' to imply that the operation is a move.

All types that have inter-convertible data have functions for all of the relevant types, generally an move operation from another variable or a copy operation that creates a new variable in the relevant type.

All of the types are based on structs in order to simplify error handling (each type stores an error value from the last operation that can produce one) and if one might want to perform a subsequent operation on the same variable, because it returns the pointer to itself, the resultant expression can have further pointer receiver methods chained together.

Because the variables are all stored in structs with no exported variable symbols, functions must exist that equate to assignment operations, though with the caveat mentioned at the top - unless one specifically uses the copy function (each type will have a ToThisSameType() function, which functions as a copy), data will always be purged before new data is stored. All other operations will function destructively on the source by default as this is best practice to limit the results of bugs in the code leaking information unintentionally.

*/
package types
