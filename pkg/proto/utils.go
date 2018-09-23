package proto

// Zero makes all the bytes in a slice zero
func Zero(b *[]byte) {
	B := *b
	for i := range B {
		B[i] = 0
	}
}
