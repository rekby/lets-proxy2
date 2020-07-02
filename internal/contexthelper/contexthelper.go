package contexthelper

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type CombinedContext struct {
	contexts []context.Context

	deadline   time.Time
	deadlineOk bool

	mu   sync.RWMutex
	done chan struct{}
	err  error
}

type droppedCancelContext struct {
	ctx context.Context
}

// Deadline return minimum of contextx deadlines, if present
func (cc *CombinedContext) Deadline() (deadline time.Time, ok bool) {
	return cc.deadline, cc.deadlineOk
}

// Done return channel, that will close when first of initial contexts will closed
func (cc *CombinedContext) Done() <-chan struct{} {
	return cc.done
}

// Err return error, of context, which close handled first
func (cc *CombinedContext) Err() error {
	cc.mu.RLock()
	res := cc.err
	cc.mu.RUnlock()

	return res
}

// Value return value of key. It iterate  over initial context one by one and return first not nil value.
// If all of contextx return nil - return nil too
func (cc *CombinedContext) Value(key interface{}) interface{} {
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
func CombineContext(contexts ...context.Context) *CombinedContext {
	res := &CombinedContext{
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
	// handlepanic: no external call
	go res.waitCloseAny()
	return res
}

func (cc *CombinedContext) waitCloseAny() {
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

func DropCancelContext(ctx context.Context) context.Context {
	return droppedCancelContext{ctx}
}

func (d droppedCancelContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (d droppedCancelContext) Done() <-chan struct{} {
	return nil
}

func (d droppedCancelContext) Err() error {
	return nil
}

func (d droppedCancelContext) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}
