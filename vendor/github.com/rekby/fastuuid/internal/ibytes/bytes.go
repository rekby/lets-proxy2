//go:build go1.20
// +build go1.20

package ibytes

import "bytes"

func Clone(b []byte) []byte {
	return bytes.Clone(b)
}
