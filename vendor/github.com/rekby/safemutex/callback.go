package safemutex

type ReadWriteCallback[T any] func(value T) (newValue T)
