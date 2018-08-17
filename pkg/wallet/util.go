package wallet

import (
	"fmt"
	"github.com/awnumar/memguard"
)

var AllocatedBufferCount int
var AllocatedBufferTotalSize int
var AllLockedBuffers []*memguard.LockedBuffer

// Allocates a new memguard LockedBuffer of a given size and keeps track of it
func NewBuffer(size int) (B *memguard.LockedBuffer, err error) {
	B, err = memguard.NewMutable(size)
	if err != nil {
		fmt.Println("Failed to allocate LockedBuffer")
		return
	}
	AllLockedBuffers = append(AllLockedBuffers, B)
	AllocatedBufferCount++
	AllocatedBufferTotalSize += B.Size()
	return
}

// Allocates a new memguard Lockbuffer and fills it with the contents of a byte slice and keeps track of it
func NewBufferFromBytes(b []byte) (B *memguard.LockedBuffer, err error) {
	B, err = memguard.NewMutableFromBytes(b)
	if err != nil {
		fmt.Println("Failed to allocate LockedBuffer")
		return
	}
	AllLockedBuffers = append(AllLockedBuffers, B)
	AllocatedBufferCount++
	AllocatedBufferTotalSize += B.Size()
	return
}

// Deallocates a memguard LockedBuffer and removes it from its register
func DeleteBuffer(b *memguard.LockedBuffer) {
	if AllocatedBufferCount < 2 {
		fmt.Println("Who is trying to delete an item we aren't tracking???")
		memguard.SafeExit(1)
	}
	var me int
	for i := range AllLockedBuffers {
		if AllLockedBuffers[i] == b {
			me = i
			break
		}
	}
	switch {
	case me == len(AllLockedBuffers):
		AllLockedBuffers = AllLockedBuffers[:len(AllLockedBuffers)-1]
	case me == 0:
		AllLockedBuffers = AllLockedBuffers[1:]
	default:
		AllLockedBuffers = append(AllLockedBuffers[:me], AllLockedBuffers[me+1:]...)
	}
	AllocatedBufferCount--
	AllocatedBufferTotalSize -= b.Size()
	b.Destroy()
}

// Deallocates a slice of buffers
func DeleteBuffers(b ...*memguard.LockedBuffer) {
	for i := range b {
		DeleteBuffer(b[i])
	}
}

// // FormatString prepends a byte with the length of a string for wire/storage formatting
// func FormatString(s string) (result []byte) {
// 	result = append([]byte{byte(len(s))}, s...)
// 	return
// }

// // FormatBytes prepends a byte with the length of a byte slice for wire/storage formatting
// func FormatBytes(b []byte) (result []byte) {
// 	result = append([]byte{byte(len(b))}, b...)
// 	return
// }
// func BytesToUint32(in []byte) (result uint32) {
// 	resultB := bytes.NewBuffer(in)
// 	binary.Read(resultB, binary.LittleEndian, &result)
// 	return
// }
// func Uint32ToBytes(in uint32) (result []byte) {
// 	resultB := bytes.NewBuffer([]byte{})
// 	binary.Write(resultB, binary.LittleEndian, &in)
// 	return resultB.Bytes()
// }
// func BytesToInt64(in []byte) (result int64) {
// 	resultB := bytes.NewBuffer(in)
// 	binary.Read(resultB, binary.LittleEndian, &result)
// 	return
// }
// func Int64ToBytes(in int64) (result []byte) {
// 	result = make([]byte, 8)
// 	binary.LittleEndian.PutUint64(result, uint64(in))
// 	return
// }
// func BytesToUint64(in []byte) (result uint64) {
// 	resultB := bytes.NewBuffer(in)
// 	binary.Read(resultB, binary.LittleEndian, &result)
// 	return
// }
// func Uint64ToBytes(in uint64) (result []byte) {
// 	result = make([]byte, 8)
// 	binary.LittleEndian.PutUint64(result, in)
// 	return
// }

// // Append a byte slice to a byte slice in the caller's scope
// func Append(b *[]byte, B ...[]byte) {
// 	for i := range B {
// 		*b = append(*b, B[i]...)
// 	}
// }
// func ToPub(pubEC *ec.PublicKey) (pub *key.Pub) {
// 	pub = &key.Pub{}
// 	pub.SetPub(*pubEC)
// 	return
// }
// func ParsePub(pub []byte) (key *ec.PublicKey, err error) {
// 	return ec.ParsePubKey(pub, ec.S256())
// }

// // func SetPriv(priv []byte) (result *key.Priv) {
// // 	result = &key.Priv{}
// // 	result.Set(priv)
// // 	return
// // }
// func PubToHex(pub interface{}) string {
// 	return hex.EncodeToString(pub.(*key.Pub).GetPub().SerializeUncompressed())
// }

// // func PrivToHex(priv interface{}) string {
// // 	return hex.EncodeToString(priv.(*key.Priv).Get())
// // }
// func BytesToHex(b []byte) string {
// 	return hex.EncodeToString(b)
// }
// func StringToUint64(s string) (uint64, error) {
// 	return strconv.ParseUint(s, 10, 64)
// }
// func StringToInt64(s string) (int64, error) {
// 	return strconv.ParseInt(s, 10, 64)
// }
// func StringToUint32(s string) (uint32, error) {
// 	u, err := strconv.ParseUint(s, 10, 32)
// 	return uint32(u), err
// }
