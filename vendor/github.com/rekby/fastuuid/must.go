package fastuuid

// MustUUIDv4 generate uuid without allocations
// It will panic if it can't read random data.
//
// It is safe calling this function from concurrent goroutines.
func MustUUIDv4() [16]byte {
	res, err := UUIDv4()
	if err != nil {
		panic(err)
	}
	return res
}

// MustUUIDv4String generate uuid random string with only one allocation
// It will panic if it can't read random data.
//
// It is safe calling this function from concurrent goroutines.
func MustUUIDv4String() string {
	res, err := UUIDv4String()
	if err != nil {
		panic(err)
	}
	return res
}

// MustUUIDv4StringBytes generate uuid and render it as string to dst buffer without allocations
// It will panic if it can't read random data.
//
// It is safe calling this function from concurrent goroutines.
func MustUUIDv4StringBytes(dst []byte) {
	err := UUIDv4StringBytes(dst)
	if err != nil {
		panic(err)
	}
}
