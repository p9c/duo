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

// U160 is a 160 bit unsigned integer, used for storing ID hashes of keys and scripts
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

// Zero160 returns an empty uint256
func Zero160() *U160 {
	return &U160{}
}

// Assign stores a new value in thte current variable
func (u *U160) Assign(input *U160) *U160 {
	if input.BitLen() <= Bitwidth160 {
		u = input
	} else {
		u.truncate(input)
	}
	return u
}

// ToString converts a U256 into a string containing a decimal representing the value
func (u *U160) ToString() string {
	if u.BitLen() > Bitwidth160 {
		u.truncate(u)
	}
	return u.String()
}

// ToBytes returns a byte slice containing the U256
func (u *U160) ToBytes() []byte {
	if u.BitLen() > Bitwidth160 {
		u.truncate(u)
	}
	return u.Bytes()
}

// FromUint64 puts a Uint64 value into the U256
func (u *U160) FromUint64(input uint64) *U160 {
	u.SetUint64(input)
	return u
}

// FromString converts a string, attempting to autodetect base, into a U256
func (u *U160) FromString(input string) *U160 {
	u.SetString(input, 0)
	if u.BitLen() > Bitwidth160 {
		u.truncate(u)
	}
	return u
}

// FromBytes converts a byte slice into a U256
func (u *U160) FromBytes(input []byte) *U160 {
	if len(input) > Bytewidth160 {
		input = input[:Bytewidth160]
	}
	u.SetBytes(input)
	return u
}

// EQ returns true if another U160 is equal to this one
func (u *U160) EQ(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) == 0
}

// NEQ returns true if another U160 is not equal to this one
func (u *U160) NEQ(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) != 0
}

// GT returns true if another U160 is greater than this one
func (u *U160) GT(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) == 1
}

// LT returns true if another U160 is less than this one
func (u *U160) LT(operand *U160) bool {
	return u.Int.Cmp(&operand.Int) == -1
}

// GToE returns true if another U160 is greater than or equal to this one
func (u *U160) GToE(operand *U160) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == 1
}

// LToE returns true if another U160 is less thtan oor equal to this one
func (u *U160) LToE(operand *U160) bool {
	i := u.Int.Cmp(&operand.Int)
	return i == 0 || i == -1
}

// Not returs the binary inverse and changes the current value to this
func (u *U160) Not() *U160 {
	u.Int.Not(&u.Int)
	return u
}

// And returns the binary AND with another U160 and changes its value to this
func (u *U160) And(operand *U160) *U160 {
	u.Int.And(&u.Int, &operand.Int)
	return u
}

// Or returns the binary OR with another U160 and changes its value to this
func (u *U160) Or(operand *U160) *U160 {
	u.Int.Or(&u.Int, &operand.Int)
	return u
}

// Xor returns the binary XOR with another U160 and changes its value to this
func (u *U160) Xor(operand *U160) *U160 {
	u.Int.Xor(&u.Int, &operand.Int)
	return u
}

// Add adds another U160 to this one
func (u *U160) Add(operand *U160) *U160 {
	u.Int.Add(&u.Int, &operand.Int)
	return u
}

// Sub subtracts another U160 to this one
func (u *U160) Sub(operand *U160) *U160 {
	u.Int.Sub(&u.Int, &operand.Int)
	return u
}

// RIPEMD160 takes a byte slice and returs a U160 containing the hash of the slice
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
