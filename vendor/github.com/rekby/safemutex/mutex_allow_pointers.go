package safemutex

// MutexWithPointers contains guarded value inside, access to value allowed inside callbacks only
// it allow to guarantee not access to the value without lock the mutex
// zero value is usable as mutex with default options and zero value of guarded type
type MutexWithPointers[T any] struct {
	mutexBase[T]
}

// NewWithPointers create Mutex with initial value and default options.
// NewWithPointers call internal checks for T and panic if checks failed, see MutexOptions for details
func NewWithPointers[T any](value T) MutexWithPointers[T] {
	res := MutexWithPointers[T]{
		mutexBase: mutexBase[T]{
			value: value,
		},
	}

	//nolint:govet
	//goland:noinspection GoVetCopyLock
	return res
}

// Lock - call f within locked mutex.
// it will panic with ErrPoisoned if previous call exited without return value:
// with panic or runtime.Goexit()
func (m *MutexWithPointers[T]) Lock(f ReadWriteCallback[T]) {
	m.m.Lock()
	defer m.m.Unlock()

	m.baseValidateLocked()
	m.callLocked(f)
}
