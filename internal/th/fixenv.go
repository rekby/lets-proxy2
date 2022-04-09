package th

import (
	"context"
	"sync"
	"testing"

	"github.com/maxatome/go-testdeep"

	"github.com/rekby/fixenv"
)

type Env struct {
	Ctx context.Context
	*fixenv.EnvT
	TD
}

func NewEnv(t *testing.T) (env *Env, ctx context.Context, cancel func()) {
	td := testdeep.NewT(t)
	ctx, ctxCancel := TestContext(td)
	tWrap := &testWrapper{T: td}
	env = &Env{
		EnvT: fixenv.NewEnv(tWrap),
		Ctx:  ctx,
		TD:   TD{T: td},
	}
	tWrap.Cleanup(ctxCancel)
	return env, ctx, tWrap.startCleanups
}

func (e *Env) T() fixenv.T {
	return e.EnvT.T()
}

// TD struct need for rename embedded field in Env
type TD struct {
	*testdeep.T
}

type testWrapper struct {
	*testdeep.T

	m               sync.Mutex
	cleanups        []func()
	cleanupsStarted bool
}

func (t *testWrapper) Cleanup(f func()) {
	t.m.Lock()
	defer t.m.Unlock()

	t.cleanups = append(t.cleanups, f)
}

func (t *testWrapper) startCleanups() {
	t.m.Lock()
	started := t.cleanupsStarted
	if !started {
		t.cleanupsStarted = true
	}
	t.m.Unlock()

	if started {
		return
	}

	for i := len(t.cleanups) - 1; i >= 0; i-- {
		f := t.cleanups[i]
		f()
	}
}
