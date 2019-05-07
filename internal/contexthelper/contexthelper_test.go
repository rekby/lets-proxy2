package contexthelper

import (
	"context"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
)

var _ context.Context = &combinedContext{}

func TestCombinedContext_Deadline(t *testing.T) {
	td := testdeep.NewT(t)
	var ctx *combinedContext
	var deadline time.Time
	var ok bool

	ctx = CombineContext()
	_, ok = ctx.Deadline()
	td.False(ok)

	ctx = CombineContext(context.Background(), context.Background())
	_, ok = ctx.Deadline()
	td.False(ok)

	ctx = CombineContext(context.Background(), context.Background())
	_, ok = ctx.Deadline()
	td.False(ok)

	contextDeadline1, cancel1 := context.WithDeadline(context.Background(), time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	defer cancel1()

	ctx = CombineContext(contextDeadline1)
	deadline, ok = ctx.Deadline()
	td.True(ok)
	td.CmpDeeply(deadline, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))

	contextDeadline2, cancel2 := context.WithDeadline(context.Background(), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC))
	defer cancel2()

	ctx = CombineContext(contextDeadline1, contextDeadline2)
	deadline, ok = ctx.Deadline()
	td.True(ok)
	td.CmpDeeply(deadline, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))

	ctx = CombineContext(contextDeadline2, contextDeadline1)
	deadline, ok = ctx.Deadline()
	td.True(ok)
	td.CmpDeeply(deadline, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
}

func TestCombinedContext_Err(t *testing.T) {
	td := testdeep.NewT(t)
	var ctx *combinedContext

	wait := func() { time.Sleep(time.Millisecond) }

	ctx = CombineContext()
	td.CmpNoError(ctx.Err())

	ctx = CombineContext(context.Background())
	td.CmpNoError(ctx.Err())

	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx = CombineContext(ctx1)
	cancel1()
	wait()
	td.CmpDeeply(ctx.Err(), ctx1.Err())

	ctx2, cancel2 := context.WithDeadline(context.Background(), time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	ctx = CombineContext(context.Background(), ctx2)
	cancel2()
	wait()
	td.CmpDeeply(ctx.Err(), ctx2.Err())

	ctxTimeout1, ctxTimeout1Cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer ctxTimeout1Cancel()

	ctx = CombineContext(ctx1, ctxTimeout1)
	wait()
	td.CmpDeeply(ctx.Err(), ctx1.Err())

	ctxTimeout2, ctxTimeout2Cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer ctxTimeout2Cancel()

	ctx3, ctx3Cancel := context.WithCancel(context.Background())

	ctx = CombineContext(ctxTimeout2, ctx3)
	wait()
	td.CmpNoError(ctx.Err())
	ctx3Cancel()
	wait()

	td.CmpDeeply(ctx.Err(), ctx3.Err())
}

func TestCombinedContext_Value(t *testing.T) {
	td := testdeep.NewT(t)
	var ctx *combinedContext

	ctx = CombineContext()
	td.CmpDeeply(ctx.Value(1), nil)

	ctx1 := context.WithValue(context.Background(), 1, 2)
	ctx = CombineContext(ctx1)
	td.CmpDeeply(ctx.Value(1), 2)

	ctx2 := context.WithValue(context.Background(), 1, 3)
	ctx = CombineContext(ctx1, ctx2)
	td.CmpDeeply(ctx.Value(1), 2)

	ctx = CombineContext(ctx2, ctx1)
	td.CmpDeeply(ctx.Value(1), 3)
}

func TestCombinedContext_Done(t *testing.T) {
	td := testdeep.NewT(t)
	var ctx *combinedContext
	var done bool

	wait := func() { time.Sleep(10 * time.Millisecond) }

	ctx1, ctx1Cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer ctx1Cancel()

	done = false
	ctx = CombineContext(ctx1)
	go func() {
		<-ctx.Done()
		done = true
	}()

	wait()
	td.True(done)

}