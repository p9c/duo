package Uint
import (
	"golang.org/x/crypto/ripemd160"
	"math/big"
)
const (
	// Bitwidth160 is the number of bits in a U160
	Bitwidth160 = 160
	// Bytewidth160 is the number of bytes in a U160
	Bytewidth160 = 20
)
// U160 stores the 160 bit integer value in a big.Int
type U160 struct {
	big.Int
}
type u160 interface {
	Zero256() *U256
	Assign(*U160) *U160
	ToString() string
	ToBytes() []byte
	FromUint64(uint64) *U160
	FromString(string) *U160
	FromBytes([]byte) *U160
	EQ(*U160) bool
	NEQ(*U160) bool
	GT(*U160) bool
	LT(*U160) bool
	GToE(*U160) bool
	LToE(*U160) bool
	Not(*U160) *U160
	And(*U160) *U160
	Or(*U160) *U160
	Xor(*U160) *U160
	Add(*U160) *U160
	Sub(*U160) *U160
}
func (u *U160) truncate(input *U160) *U160 {
	u.Int = *u.SetBytes(input.Bytes()[:Bytewidth160])
	return u
}
// Returns an empty uint256
func Zero160() *U160 {
	return &U160{}
}
// Stores a new value in the current variable as well returning the value
func (u *U160) Assign(input *U160) *U160 {
	if input.BitLen() <= Bitwidth160 {
		u = input
	} else {
		u.truncate(input)
	}
	return u
}
// Converts a U256 into a string containing a decimal representing the value
func (u *U160) ToString() string {
	if u.BitLen() > Bitwidth160 {
		u.truncate(u)
	}
	return u.String()
}
// Returns a byte slice containing the U256 raw binary data
func (u *U160) ToBytes() []byte {
	if u.BitLen() > Bitwidth160 {
		u.truncate(u)
	}
	return u.Bytes()
}
// Puts a Uint64 value into the U256 and returns a U160
func (u *U160) FromUint64(input uint64) *U160 {
	u.SetUint64(input)
	return u
}
// Converts a string, attempting to autodetect base, into a U256, truncates if the number is too large for 160 bits
func (u *U160) FromString(input string) *U160 {
	u.SetString(input, 0)
	if u.BitLen() > Bitwidth160 {
		u.truncate(u)
	}
	return u
}
// Converts a byte slice into a U160, and returns the U160
func (u *U160) FromBytes(input []byte) *U160 {
	if len(input) > Bytewidth160 {
		input = input[:Bytewidth160]
	}
	u.SetBytes(input)
	return u
}
// Returns true if the operand is equal to the value in the receiver
func (u *U160) EQ(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) == 0
}
// Returns true if the operand is not equal to the value in the receiver
func (u *U160) NEQ(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) != 0
}
// Returns true if the operand is greater than the value in the receiver
func (u *U160) GT(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) == 1
}
// Returns true if the operand is less than the value in the receiver
func (u *U160) LT(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) == -1
}
// Returns true if the operand is greater than or equal to the value in the receiver
func (u *U160) GToE(operand *U160) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == 1
}
// Returns true if the operand is less than or equal to the value in the receiver
func (u *U160) LToE(operand *U160) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == -1
}
// Returns the bitwise inversion of the value in the receiver and stores it
func (u *U160) Not() *U160 {
	u.Int.Not(&u.Int)
	return u
}
// Returns the binary AND with the operand and stores it
func (u *U160) And(operand *U160) *U160 {
	u.Int.And(&u.Int, &operand.Int)
	return u
}
// Returns the binary OR with the operand and stores it
func (u *U160) Or(operand *U160) *U160 {
	u.Int.Or(&u.Int, &operand.Int)
	return u
}
// Returns the binary XOR with another U160 and stores it
func (u *U160) Xor(operand *U160) *U160 {
	u.Int.Xor(&u.Int, &operand.Int)
	return u
}
// Returns the sum of the value in the receiver and the operand and stores it
func (u *U160) Add(operand *U160) *U160 {
	u.Int.Add(&u.Int, &operand.Int)
	return u
}
// Returns the difference of the value in the receiver and the operand and stores it
func (u *U160) Sub(operand *U160) *U160 {
	u.Int.Sub(&u.Int, &operand.Int)
	return u
}
// Returns the RIPEMD160 hash of a list of byte slices
func RIPEMD160(b ...[]byte) *U160 {
	var data []byte
	for i := range b {
		data = append(data, b[i]...)
	}
	digest := ripemd160.New()
	sum := digest.Sum(data)
	var out U160
	out.Int.SetBytes(sum)
	return &out
}
