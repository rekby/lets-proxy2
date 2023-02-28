//go:build go1.20
// +build go1.20

package internal

import "unsafe"

func UnsafeConvertToString(src []byte) string {
	return unsafe.String(&src[0], len(src))
}
