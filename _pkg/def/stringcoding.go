package def

// StringCoding is a coding type for converting data to strings
type StringCoding int

// Get returns the current coding type
func (r StringCoding) Get() string {
	switch {
	case int(r) > len(StringCodingTypes):
		Dbg.Print("code higher than maximum")
		r = 0
	case int(r) < 0:
		Dbg.Print("negative coding")
		r = 0
		return StringCodingTypes[r]
	}
	return StringCodingTypes[r]
}

// Set sets the coding type
func (r StringCoding) Set(s string) Coding {
	r = 0
	for i := range StringCodingTypes {
		if s == StringCodingTypes[i] {
			r = StringCoding(i)
			break
		}
	}
	return r
}

// List returns an array of strings representing available coding
func (r StringCoding) List() (R []string) {
	for i := range StringCodingTypes {
		R = append(R, StringCodingTypes[i])
	}
	return
}
