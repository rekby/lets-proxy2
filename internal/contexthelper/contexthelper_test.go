package contexthelper

import (
	"context"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
)

var _ context.Context = &CombinedContext{}

func TestCombinedContext_Deadline(t *testing.T) {
	td := testdeep.NewT(t)
	var ctx *CombinedContext
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
	var ctx *CombinedContext

	const waitPause = 10 * time.Millisecond

	wait := func() { time.Sleep(waitPause) }

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

	ctxTimeout1, ctxTimeout1Cancel := context.WithTimeout(context.Background(), 10*waitPause)
	defer ctxTimeout1Cancel()

	ctx = CombineContext(ctx1, ctxTimeout1)
	wait()
	td.CmpDeeply(ctx.Err(), ctx1.Err())

	ctxTimeout2, ctxTimeout2Cancel := context.WithTimeout(context.Background(), 10*waitPause)
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
	type ctxKeyType int

	const one ctxKeyType = 1

	td := testdeep.NewT(t)
	var ctx = CombineContext()

	td.CmpDeeply(ctx.Value(one), nil)

	ctx1 := context.WithValue(context.Background(), one, 2)
	ctx = CombineContext(ctx1)
	td.CmpDeeply(ctx.Value(one), 2)

	ctx2 := context.WithValue(context.Background(), one, 3)
	ctx = CombineContext(ctx1, ctx2)
	td.CmpDeeply(ctx.Value(one), 2)

	ctx = CombineContext(ctx2, ctx1)
	td.CmpDeeply(ctx.Value(one), 3)
}

func TestCombinedContext_Done(t *testing.T) {
	td := testdeep.NewT(t)

	var ctx *CombinedContext

	ctx1, ctx1Cancel := context.WithCancel(context.Background())

	ctx = CombineContext(ctx1)
	td.CmpNoError(ctx.Err())

	select {
	case <-ctx.Done():
		t.Error()
	default:
		// pass
	}

	ctx1Cancel()

	time.Sleep(time.Millisecond * 100)

	td.CmpError(ctx.Err())

	select {
	case <-ctx.Done():
		// pass
	default:
		t.Error()
	}
}

func TestDropCancelContext(t *testing.T) {
	td := testdeep.NewT(t)
	type keyType string
	const key keyType = "key"
	const val = "val"

	ctx, ctxCancel := context.WithCancel(context.WithValue(context.Background(), key, val))
	dropCancel := DropCancelContext(ctx)
	ctxCancel()

	deadline, ok := dropCancel.Deadline()
	td.Cmp(deadline, time.Time{})
	td.False(ok)

	td.CmpNoError(dropCancel.Err())
	td.Nil(dropCancel.Done())
	td.Cmp(dropCancel.Value(key), val)
}
