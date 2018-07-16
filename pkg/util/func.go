package util

func Reverse(bytes []byte) (result []byte) {
	result = make([]byte, len(bytes))
	position := len(bytes) - 1
	for i := range bytes {
		result[position-i] = bytes[i]
	}
	return
}
