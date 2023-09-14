package safemutex

import (
	"reflect"
	"sync"
)

// Mutex contains guarded value inside, access to value allowed inside callbacks only
// it allow to guarantee not access to the value without lock the mutex
// zero value is usable as mutex with default options and zero value of guarded type
type Mutex[T any] struct {
	m           sync.Mutex
	value       T
	options     MutexOptions
	initialized bool
	errWrap     errWrap
}

// New create Mutex with initial value and default options.
// New call internal checks for T and panic if checks failed, see MutexOptions for details
func New[T any](value T) Mutex[T] {
	return NewWithOptions(value, MutexOptions{})
}

// NewWithOptions create Mutex with initial value and custom options.
// MutexOptions allow to reduce default security when it needs.
// NewWithOptions call internal checks for T and panic if checks failed, see MutexOptions for details
func NewWithOptions[T any](value T, options MutexOptions) Mutex[T] {
	res := Mutex[T]{
		value:   value,
		options: options,
	}

	res.validateLocked()

	//nolint:govet
	//goland:noinspection GoVetCopyLock
	return res
}

// Lock - call f within locked mutex.
// it will panic if value type not pass internal checks
// it will panic with ErrPoisoned if previous call exited without return value:
// with panic or runtime.Goexit()
func (m *Mutex[T]) Lock(f ReadWriteCallback[T]) {
	m.m.Lock()
	defer m.m.Unlock()

	m.callLocked(f)
}

func (m *Mutex[T]) callLocked(f ReadWriteCallback[T]) {
	m.validateLocked()

	hasPanic := true

	defer func() {
		if hasPanic && !m.options.AllowPoisoned {
			m.errWrap = errWrap{ErrPoisoned}
		}
	}()

	m.value = f(m.value)
	hasPanic = false

}

func (m *Mutex[T]) validateLocked() {
	if m.errWrap.err != nil {
		panic(m.errWrap)
	}

	if m.initialized {
		return
	}

	m.initialized = true

	if !m.options.AllowPointers {
		if checkTypeCanContainPointers(reflect.TypeOf(m.value)) {
			panic(errContainPointers)
		}
	}
}

type MutexOptions struct {
	AllowPointers bool
	AllowPoisoned bool
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
