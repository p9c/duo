package sync

func removeTrailingZeroes(in []byte) []byte {
	length := 0
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] != 0 {
			length = i
			break
		}
	}
	return in[:length+1]
}

func removeLeadingZeroes(in []byte) []byte {
	nonzerostart := 0
	for i := range in {
		if in[i] != 0 {
			nonzerostart = i
			break
		}
	}
	return in[nonzerostart:]
}
