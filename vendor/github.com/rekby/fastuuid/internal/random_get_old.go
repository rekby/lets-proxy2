//go:build !go1.18
// +build !go1.18

package internal

import "github.com/valyala/fastrand"

func GetRandomReader() *RandomReader {
	for {
		index := fastrand.Uint32n(uint32(readersCount))
		randReaders[index].m.Lock()
		return &randReaders[index]
	}
}
