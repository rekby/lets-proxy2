package safemutex

import (
	"reflect"
	"sync"
)

// Mutex contains guarded value inside, access to value allowed inside callbacks only
// it allow to guarantee not access to the value without lock the mutex
// zero value is usable as mutex with default options and zero value of guarded type
// Mutex deny to save value with any type of pointers, which allow accidentally change internal state.
// it will panic if T contains any pointer.
type Mutex[T any] struct {
	mutexBase[T, sync.Mutex]
	initOnce    sync.Once
	initialized bool // for tests only
}

// New create Mutex with initial value and default options.
// New call internal checks for T and panic if checks failed, see MutexOptions for details
func New[T any](value T) Mutex[T] {
	res := Mutex[T]{
		mutexBase: mutexBase[T, sync.Mutex]{
			value: value,
		},
	}

	res.validateLocked()

	//nolint:govet
	//goland:noinspection GoVetCopyLock
	return res
}

// Lock - call f within locked mutex.
// it will panic if value type not pass internal checks
// it will panic with ErrPoisoned if previous locked call exited without return value:
// with panic or runtime.Goexit()
func (m *Mutex[T]) Lock(f ReadWriteCallback[T]) {
	m.m.Lock()
	defer m.m.Unlock()

	m.validateLocked()
	m.callLocked(f)
}

func (m *Mutex[T]) validateLocked() {
	m.baseValidateLocked()

	m.initOnce.Do(m.initLocked)
}

func (m *Mutex[T]) initLocked() {
	if checkTypeCanContainPointers(reflect.TypeOf(m.value)) {
		m.errWrap.err = errContainPointers
		panic(m.errWrap)
	}
	m.initialized = true
}

// checkTypeCanContainPointers check the value for potential contain pointers
// return true only of value guaranteed without any pointers and false in other cases (has pointers or unknown)
func checkTypeCanContainPointers(t reflectType) bool {
	if t == nil {
		return true
	}
	switch t.Kind() {
	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Bool, reflect.Complex64, reflect.Complex128, reflect.Float32, reflect.Float64,
		reflect.String:
		return false
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			structField := t.Field(i)
			if checkTypeCanContainPointers(structField.Type) {
				return true
			}
		}
		return false
	case reflect.Array:
		return checkTypeCanContainPointers(t.Elem())
	case reflect.Pointer, reflect.UnsafePointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.Interface,
		reflect.Func:
		return true
	default:
		return true
	}
}

type reflectType interface {
	Kind() reflect.Kind
	NumField() int
	Field(i int) reflect.StructField
	Elem() reflect.Type
}
