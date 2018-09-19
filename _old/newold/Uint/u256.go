package Uint
import (
	"crypto/sha256"
	"math/big"
)
const (
	// Bitwidth256 is the number of bits in a U256
	Bitwidth256 = 256
	// Bytewidth256 is the number of bytes in a U256
	Bytewidth256 = 32
)
// U256 stores a 256 bit value in a big.Int
type U256 struct {
	big.Int
}
type u256 interface {
	Zero256() *U256
	Assign(*U256) *U256
	ToString() string
	ToBytes() []byte
	FromUint64(uint64) *U256
	FromString(string) *U256
	FromBytes([]byte) *U256
	EQ(*U256) bool
	NEQ(*U256) bool
	GT(*U256) bool
	LT(*U256) bool
	GToE(*U256) bool
	LToE(*U256) bool
	Not(*U256) *U256
	And(*U256) *U256
	Or(*U256) *U256
	Xor(*U256) *U256
	Add(*U256) *U256
	Sub(*U256) *U256
}
// Zero256 returns an empty uint256
func Zero256() *U256 {
	return &U256{}
}
func (u *U256) truncate(input *U256) *U256 {
	u.Int = *u.SetBytes(input.Bytes()[:Bytewidth256])
	return u
}
// Stores a U256 in the receiver and returns it
func (u *U256) Assign(input *U256) *U256 {
	if input.BitLen() <= Bitwidth256 {
		u = input
	} else {
		u.truncate(input)
	}
	return u
}
// Returns a base 10 representation of the value in the receiver as a string
func (u *U256) ToString() string {
	if u.BitLen() > Bitwidth256 {
		u.truncate(u)
	}
	return u.String()
}
// Returns the 256 bit integer in the receiver as a byte slice
func (u *U256) ToBytes() []byte {
	if u.BitLen() > Bitwidth256 {
		u.truncate(u)
	}
	return u.Bytes()
}
// Converts a uint64 to U256, stores it and returns it
func (u *U256) FromUint64(input uint64) *U256 {
	u.SetUint64(input)
	return u
}
// FromString converts a string, detecting hexadecimal or octal prefix as required and converting to U256 and returning to the caller
func (u *U256) FromString(input string) *U256 {
	u.SetString(input, 0)
	if u.BitLen() > Bitwidth256 {
		u.truncate(u)
	}
	return u
}
// Converts a byte slice to U256, truncating or padding as required, storing in the receiver, and returning the derived value
func (u *U256) FromBytes(input []byte) *U256 {
	if len(input) > Bytewidth256 {
		input = input[:Bytewidth256]
	}
	u.SetBytes(input)
	return u
}
// Returns true if the operand is equal to the value in the receiver
func (u *U256) EQ(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) == 0
}
// Returns true if the operand is not equal to the value in the receiver
func (u *U256) NEQ(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) != 0
}
// Returns true if the operand is greater than the value in the receiver
func (u *U256) GT(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) == 1
}
// Returns true if the operand is lesser than the value in the receiver
func (u *U256) LT(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) == -1
}
// Returns true if the operand is greater than or equal to the value in the receiver
func (u *U256) GToE(operand *U256) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == 1
}
// Returns true if the operand is less than or equal to the value in the receiver
func (u *U256) LToE(operand *U256) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == -1
}
// Returns the bit-inverse of a 256 bit integer and stores it in the receiver
func (u *U256) Not() *U256 {
	u.Int.Not(&u.Int)
	return u
}
// Performs a logical AND between one U256 and another and stores the value in the receiver
func (u *U256) And(operand *U256) *U256 {
	u.Int.And(&u.Int, &operand.Int)
	return u
}
// Returns the logical OR between one 256 bit value and another and stores the result in the receiver
func (u *U256) Or(operand *U256) *U256 {
	u.Int.Or(&u.Int, &operand.Int)
	return u
}
// Returns the logical Exclusive OR between one 256 value and another and stores the value in the receiver
func (u *U256) Xor(operand *U256) *U256 {
	u.Int.Xor(&u.Int, &operand.Int)
	return u
}
// Adds a value to a 256 bit integer and stores the result in the receiver
func (u *U256) Add(operand *U256) *U256 {
	u.Int.Add(&u.Int, &operand.Int)
	return u
}
// Subtracts a value from a 256 bit integer and stores the result in the receiver
func (u *U256) Sub(operand *U256) *U256 {
	u.Int.Sub(&u.Int, &operand.Int)
	return u
}
// Returns the SHA256 hash of a one or more byte slices
func SHA256(b ...[]byte) *U256 {
	var data []byte
	for i := range b {
		data = append(data, b[i]...)
	}
	sum := sha256.Sum256(data)
	var out U256
	out.Int.SetBytes(sum[:])
	return &out
}
