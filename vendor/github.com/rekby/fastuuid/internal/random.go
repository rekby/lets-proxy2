package internal

import (
	"crypto/rand"
	"github.com/rekby/fastuuid/internal/ibufio"
	"io"
	"runtime"
	"sync"
)

const bufSize = 1024

var readersCount = runtime.NumCPU() * 4
var randReaders []RandomReader
var randomSource = rand.Reader

func init() {
	initializeRandReaders(readersCount)
}

func initializeRandReaders(count int) {
	readersCount = count
	randReaders = make([]RandomReader, count)
	for i := 0; i < count; i++ {
		randReaders[i].reader = *ibufio.NewReaderSize(randomSource, bufSize)
	}
}

// SetRandomSource is NOT thread-safe
func SetRandomSource(reader io.Reader) {
	if reader == nil {
		reader = rand.Reader
	}
	randomSource = reader
	initializeRandReaders(readersCount)
}

type RandomReader struct {
	m      sync.Mutex
	reader ibufio.Reader
	_      [64]byte // padding, for prevent share same cache line by neighbours
}

func (r *RandomReader) ReadFull(dst []byte) error {
	// modified copy of io.ReadAtLeast, need for prevent work with interface
	min := len(dst)
	var n int
	for n < min {
		nn, err := r.reader.Read(dst[n:])
		if err != nil {
			return err
		}
		n += nn
	}
	if n < min {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (r *RandomReader) Release() {
	r.m.Unlock()
}
