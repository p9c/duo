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

// U256 is a 256 bit unsigned integer
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

// Assign copies a 256 bit integer from another variable
func (u *U256) Assign(input *U256) *U256 {
	if input.BitLen() <= Bitwidth256 {
		u = input
	} else {
		u.truncate(input)
	}
	return u
}

// ToString converts the 256 bit integer into a decimal 10 number as a string
func (u *U256) ToString() string {
	if u.BitLen() > Bitwidth256 {
		u.truncate(u)
	}
	return u.String()
}

// ToBytes converts a 256 bit integer to a byte slice
func (u *U256) ToBytes() []byte {
	if u.BitLen() > Bitwidth256 {
		u.truncate(u)
	}
	return u.Bytes()
}

// FromUint64 converts from an unsigned 64 bit integer to a 256 bit integer
func (u *U256) FromUint64(input uint64) *U256 {
	u.SetUint64(input)
	return u
}

// FromString converts a string, detecting hexadecimal or octal prefix as required and converting to 256 bit integer
func (u *U256) FromString(input string) *U256 {
	u.SetString(input, 0)
	if u.BitLen() > Bitwidth256 {
		u.truncate(u)
	}
	return u
}

// FromBytes converts from a byte slice to a 256 bit integer
func (u *U256) FromBytes(input []byte) *U256 {
	if len(input) > Bytewidth256 {
		input = input[:Bytewidth256]
	}
	u.SetBytes(input)
	return u
}

// EQ returns true if the operand is equal to the value stored in a 256 bit integer
func (u *U256) EQ(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) == 0
}

// NEQ returns true if the operand is not equal to the value stored in a 256 bit integer
func (u *U256) NEQ(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) != 0
}

// GT returns true if the operand is greater than the value stored in a 256 bit integer
func (u *U256) GT(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) == 1
}

// LT returns true if the operand is less than the value stored in a 256 bit integer
func (u *U256) LT(operand *U256) bool {
	return u.Int.Cmp(&operand.Int) == -1
}

// GToE returns true if the operand is greater than or equal than the value stored in a 256 bit integer
func (u *U256) GToE(operand *U256) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == 1
}

// LToE returns true if the operand is less than or equal to the value stored in a 256 bit integer
func (u *U256) LToE(operand *U256) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == -1
}

// Not returns the bit-inverse of a 256 bit integer and also changes itself to this inverse
func (u *U256) Not() *U256 {
	u.Int.Not(&u.Int)
	return u
}

// And performs a logical AND between one 256 bit value and another and also changes itself to the result
func (u *U256) And(operand *U256) *U256 {
	u.Int.And(&u.Int, &operand.Int)
	return u
}

// Or returns the logical OR between one 256 bit value and another and also changes itself to the result
func (u *U256) Or(operand *U256) *U256 {
	u.Int.Or(&u.Int, &operand.Int)
	return u
}

// Xor returs the logical Exclusive OR between one 256 value and anothter and also changes itself to the result
func (u *U256) Xor(operand *U256) *U256 {
	u.Int.Xor(&u.Int, &operand.Int)
	return u
}

// Add adds a value to a 256 bit integer and and also changes itself to the result
func (u *U256) Add(operand *U256) *U256 {
	u.Int.Add(&u.Int, &operand.Int)
	return u
}

// Sub subtracts a value from a 256 bit integer and and also changes itself to the result
func (u *U256) Sub(operand *U256) *U256 {
	u.Int.Sub(&u.Int, &operand.Int)
	return u
}

// SHA256 returns the SHA256 hash of a byte slice
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
