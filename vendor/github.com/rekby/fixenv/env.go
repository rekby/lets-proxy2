package fixenv

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const packageScopeName = "TestMain"

var (
	globalCache               = newCache()
	globalEmptyFixtureOptions = &FixtureOptions{}

	globalMutex     = &sync.Mutex{}
	globalScopeInfo = make(map[string]*scopeInfo)
)

// EnvT - fixture cache and cleanup engine
// it created from test and pass to fixture function
// manage cache of fixtures, depends from fixture, param, test, scope.
// and call cleanup, when scope closed.
// It can be base to own, more powerful local environments.
type EnvT struct {
	t T
	c *cache

	m      sync.Locker
	scopes map[string]*scopeInfo
}

// NewEnv create EnvT from test
func NewEnv(t T) *EnvT {
	env := newEnv(t, globalCache, globalMutex, globalScopeInfo)
	env.onCreate()
	return env
}

func newEnv(t T, c *cache, m sync.Locker, scopes map[string]*scopeInfo) *EnvT {
	return &EnvT{
		t:      t,
		c:      c,
		m:      m,
		scopes: scopes,
	}
}

// T return test from EnvT created
func (e *EnvT) T() T {
	return e.t
}

// Cache call from fixture and manage call f and cache it.
// Cache must be called direct from fixture - it use runtime stacktrace for
// detect called method - it is part of cache key.
// params - part of cache key. Usually - parameters, passed to fixture.
//          it allow use parametrized fixtures with different results.
//          params must be json serializable.
// opt - fixture options, nil for default options.
// f - callback - fixture body.
// Cache guarantee for call f exactly once for same Cache called and params combination.
func (e *EnvT) Cache(params interface{}, opt *FixtureOptions, f FixtureCallbackFunc) interface{} {
	if opt == nil {
		opt = globalEmptyFixtureOptions
	}
	key, err := makeCacheKey(e.t.Name(), params, opt, false)
	if err != nil {
		e.t.Fatalf("failed to create cache key: %v", err)
		// return not reacheble after Fatalf
		return nil
	}
	wrappedF := e.fixtureCallWrapper(key, f, opt)
	res, err := e.c.GetOrSet(key, wrappedF)
	if err != nil {
		if errors.Is(err, ErrSkipTest) {
			e.T().SkipNow()
		} else {
			e.t.Fatalf("failed to call fixture func: %v", err)
		}

		// panic must be not reachable after SkipNow or Fatalf
		panic("fixenv: must be unreachable code after err check in fixture cache")
	}

	return res
}

// tearDown called from base test cleanup
// it clean env cache and call fixture's cleanups for the scope.
func (e *EnvT) tearDown() {
	e.m.Lock()
	defer e.m.Unlock()

	testName := e.t.Name()
	if si, ok := e.scopes[testName]; ok {
		cacheKeys := si.Keys()
		e.c.DeleteKeys(cacheKeys...)
		delete(e.scopes, testName)
	} else {
		e.t.Fatalf("unexpected call env tearDown for test: %q", testName)
	}
}

// onCreate register env in internal stuctures.
func (e *EnvT) onCreate() {
	e.m.Lock()
	defer e.m.Unlock()

	testName := e.t.Name()
	if _, ok := e.scopes[testName]; ok {
		e.t.Fatalf("Env exist already for scope: %q", testName)
	} else {
		e.scopes[testName] = newScopeInfo(e.t)
		e.t.Cleanup(e.tearDown)
	}
}

// makeCacheKey generate cache key
// must be called from first level of env functions - for detect external caller
func makeCacheKey(testname string, params interface{}, opt *FixtureOptions, testCall bool) (cacheKey, error) {
	externalCallerLevel := 4
	var pc = make([]uintptr, externalCallerLevel)
	var extCallerFrame runtime.Frame
	if externalCallerLevel == runtime.Callers(0, pc) {
		frames := runtime.CallersFrames(pc)
		frames.Next()                     // callers
		frames.Next()                     // the function
		frames.Next()                     // caller of the function
		extCallerFrame, _ = frames.Next() // external caller
	}
	scopeName := scopeName(testname, opt.Scope)
	return makeCacheKeyFromFrame(params, opt.Scope, extCallerFrame, scopeName, testCall)
}

func makeCacheKeyFromFrame(params interface{}, scope CacheScope, f runtime.Frame, scopeName string, testCall bool) (cacheKey, error) {
	switch {
	case f.Function == "":
		return "", errors.New("failed to detect caller func name")
	case f.File == "":
		return "", errors.New("failed to detect caller func file")
	default:
		// pass
	}

	key := struct {
		Scope        CacheScope  `json:"scope"`
		ScopeName    string      `json:"scope_name"`
		FunctionName string      `json:"func"`
		FileName     string      `json:"fname"`
		Params       interface{} `json:"params"`
	}{
		Scope:        scope,
		ScopeName:    scopeName,
		FunctionName: f.Function,
		FileName:     f.File,
		Params:       params,
	}
	if testCall {
		key.FileName = ".../" + filepath.Base(key.FileName)
	}

	keyBytes, err := json.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("failed to serialize params to json: %v", err)
	}
	return cacheKey(keyBytes), nil

}

func (e *EnvT) fixtureCallWrapper(key cacheKey, f FixtureCallbackFunc, opt *FixtureOptions) FixtureCallbackFunc {
	return func() (res interface{}, err error) {
		scopeName := scopeName(e.t.Name(), opt.Scope)

		e.m.Lock()
		si := e.scopes[scopeName]
		e.m.Unlock()

		if si == nil {
			e.t.Fatalf("Unexpected scope. Create env for test %q", scopeName)
			// not reachable
			return nil, nil
		}

		defer func() {
			si.AddKey(key)
		}()

		res, err = f()

		if opt.CleanupFunc != nil {
			si.t.Cleanup(opt.CleanupFunc)
		}

		return res, err
	}
}

func scopeName(testName string, scope CacheScope) string {
	switch scope {
	case ScopePackage:
		return packageScopeName
	case ScopeTest:
		return testName
	case ScopeTestAndSubtests:
		parts := strings.SplitN(testName, "/", 2)
		return parts[0]
	default:
		panic(fmt.Sprintf("Unknown scope: %v", scope))
	}
}
