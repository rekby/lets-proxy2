package safemutex

import "errors"

var errContainPointers = errors.New("safe mutex: value type possible to contain pointers, use NewWithOptions for allow pointers into guarded value")

var ErrPoisoned = errors.New("safe mutex: mutex poisoned (exit from callback with panic), use NewWithOptions for allow use poisoned value")

// errWrap need for deny direct compare with returned errors
type errWrap struct {
	err error
}

func (e errWrap) Error() string {
	return e.err.Error()
}

func (e errWrap) Unwrap() error {
	return e.err
}
