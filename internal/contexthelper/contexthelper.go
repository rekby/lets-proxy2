package contexthelper

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type combinedContext struct {
	contexts []context.Context

	deadline   time.Time
	deadlineOk bool

	mu   sync.RWMutex
	done chan struct{}
	err  error
}

// Deadline return minimum of contextx deadlines, if present
func (cc *combinedContext) Deadline() (deadline time.Time, ok bool) {
	return cc.deadline, cc.deadlineOk
}

// Done return channel, that will close when first of initial contexts will closed
func (cc *combinedContext) Done() <-chan struct{} {
	return cc.done
}

// Err return error, of context, which close handled first
func (cc *combinedContext) Err() error {
	cc.mu.RLock()
	res := cc.err
	cc.mu.RUnlock()

	return res
}

// Value return value of key. It iterate  over initial context one by one and return first not nil value.
// If all of contextx return nil - return nil too
func (cc *combinedContext) Value(key interface{}) interface{} {
	for _, ctx := range cc.contexts {
		val := ctx.Value(key)
		if val != nil {
			return val
		}
	}

	return nil
}

// CombineContext return combined context: common deadline, values, close.
// Minimum one of underly contexts MUST be closed for prevent memory leak.
func CombineContext(contexts ...context.Context) *combinedContext {
	res := &combinedContext{
		contexts: contexts,
		done:     make(chan struct{}),
	}

	var deadlineMin time.Time
	var deadlineOk bool

	for _, ctx := range contexts {
		deadline, ok := ctx.Deadline()
		if !ok {
			continue
		}
		if deadlineOk && deadlineMin.Before(deadline) {
			continue
		}
		deadlineOk = true
		deadlineMin = deadline
	}
	res.deadline, res.deadlineOk = deadlineMin, deadlineOk
	go res.waitCloseAny()
	return res
}

func (cc *combinedContext) waitCloseAny() {
	selectCases := make([]reflect.SelectCase, len(cc.contexts))
	for i, ctx := range cc.contexts {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		}
	}

	index, _, _ := reflect.Select(selectCases)
	err := cc.contexts[index].Err()
	cc.mu.Lock()
	cc.err = err
	close(cc.done)
	cc.mu.Unlock()
}
