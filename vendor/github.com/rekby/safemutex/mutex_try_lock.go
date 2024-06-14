//go:build go1.19
// +build go1.19

package safemutex

// TryLock - call f within locked mutex if locked successfully.
// returned true if locked successfully
// return true if mutex already locked
// it will panic if value type not pass internal checks
// it will panic with ErrPoisoned if locked successfully and previous call exited without return value:
// with panic or runtime.Goexit()
//
// Available since go 1.19 only
func (m *Mutex[T]) TryLock(f ReadWriteCallback[T]) bool {
	locked := m.m.TryLock()
	if !locked {
		return false
	}
	defer m.m.Unlock()

	m.validateLocked()
	m.callLocked(f)
	return true
}

// TryLock - call f within locked mutex if locked successfully.
// returned true if locked successfully
// return true if mutex already locked
// it will panic if value type not pass internal checks
// it will panic with ErrPoisoned if locked successfully and previous call exited without return value:
// with panic or runtime.Goexit()
//
// Available since go 1.19 only
func (m *RWMutex[T]) TryLock(f ReadWriteCallback[T]) bool {
	locked := m.m.TryLock()
	if !locked {
		return false
	}
	defer m.m.Unlock()

	m.validateLocked()
	m.callLocked(f)
	return true
}

// TryRLock - call f within read locked mutex if locked successfully.
// returned true if locked successfully
// return true if mutex already locked
// it will panic if value type not pass internal checks
// it will panic with ErrPoisoned if locked successfully and previous call exited without return value:
// with panic or runtime.Goexit()
//
// Available since go 1.19 only
func (m *RWMutex[T]) TryRLock(f ReadCallback[T]) bool {
	locked := m.m.TryRLock()
	if !locked {
		return false
	}
	defer m.m.RUnlock()

	m.validateLocked()
	m.callReadLocked(f)
	return true
}

// TryLock - call f within locked mutex if locked successfully.
// returned true if locked successfully
// return true if mutex already locked
// it will panic with ErrPoisoned if locked successfully and previous call exited without return value:
// with panic or runtime.Goexit()
//
// Available since go 1.19 only
func (m *MutexWithPointers[T]) TryLock(f ReadWriteCallback[T]) bool {
	locked := m.m.TryLock()
	if !locked {
		return false
	}
	defer m.m.Unlock()

	m.baseValidateLocked()
	m.callLocked(f)
	return true
}

// TryLock - call f within locked mutex if locked successfully.
// returned true if locked successfully
// return true if mutex already locked
// it will panic with ErrPoisoned if locked successfully and previous call exited without return value:
// with panic or runtime.Goexit()
//
// Available since go 1.19 only
func (m *RWMutexWithPointers[T]) TryLock(f ReadWriteCallback[T]) bool {
	locked := m.m.TryLock()
	if !locked {
		return false
	}
	defer m.m.Unlock()

	m.baseValidateLocked()
	m.callLocked(f)
	return true
}

// TryRLock - call f within read locked mutex if locked successfully.
// returned true if locked successfully
// return true if mutex already locked
// it will panic with ErrPoisoned if locked successfully and previous call exited without return value:
// with panic or runtime.Goexit()
//
// Available since go 1.19 only
func (m *RWMutexWithPointers[T]) TryRLock(f ReadCallback[T]) bool {
	locked := m.m.TryRLock()
	if !locked {
		return false
	}
	defer m.m.RUnlock()

	m.baseValidateLocked()
	m.callReadLocked(f)
	return true
}
