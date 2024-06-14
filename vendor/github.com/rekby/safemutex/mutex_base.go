package safemutex

import "sync"

type mutexVariant interface {
	sync.Mutex | sync.RWMutex
}

type mutexBase[T any, M mutexVariant] struct {
	m       M
	value   T
	errWrap errWrap
}

func (m *mutexBase[T, M]) baseValidateLocked() {
	if m.errWrap.err != nil {
		panic(m.errWrap)
	}
}

func (m *mutexBase[T, M]) callLocked(f ReadWriteCallback[T]) {
	m.errWrap.err = ErrPoisoned
	m.value = f(m.value)
	m.errWrap.err = nil
}

func (m *mutexBase[T, M]) callReadLocked(f ReadCallback[T]) {
	f(m.value)
}
