package th

import (
	"context"
	"github.com/rekby/safemutex"
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
	tWrap := &testWrapper{
		T: td,
	}
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

	m safemutex.MutexWithPointers[testWrapperSynced]
}

type testWrapperSynced struct {
	cleanups        []func()
	cleanupsStarted bool
}

func (t *testWrapper) Cleanup(f func()) {
	t.m.Lock(func(synced testWrapperSynced) testWrapperSynced {
		synced.cleanups = append(synced.cleanups, f)
		return synced
	})
}

func (t *testWrapper) startCleanups() {
	var started bool
	var cleanups []func()

	t.m.Lock(func(synced testWrapperSynced) testWrapperSynced {
		started := synced.cleanupsStarted
		if !started {
			synced.cleanupsStarted = true
		}
		cleanups = synced.cleanups
		return synced
	})

	if started {
		return
	}

	for i := len(cleanups) - 1; i >= 0; i-- {
		f := cleanups[i]
		f()
	}
}
