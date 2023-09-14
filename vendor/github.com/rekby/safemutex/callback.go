package safemutex

// ReadCallback receive current value, saved in mutex
type ReadCallback[T any] func(synced T)

// ReadWriteCallback receive current value, saved in mutex and return new value
type ReadWriteCallback[T any] func(synced T) T
