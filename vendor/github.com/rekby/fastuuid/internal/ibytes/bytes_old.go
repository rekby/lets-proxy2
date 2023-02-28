//go:build !go1.20
// +build !go1.20

package ibytes

// Clone returns a copy of b[:len(b)].
// The result may have additional unused capacity.
// Clone(nil) returns nil.
// copy from golang 1.20
func Clone(b []byte) []byte {
	if b == nil {
		return nil
	}
	return append([]byte{}, b...)
}
