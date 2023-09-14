package safemutex

// ReadWriteCallback receive current value, saved in mutex and return new value
type ReadWriteCallback[T any] func(synced T) T
