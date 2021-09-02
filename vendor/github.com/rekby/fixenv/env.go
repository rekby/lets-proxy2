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

type EnvT struct {
	t T
	c *cache

	m      sync.Locker
	scopes map[string]*scopeInfo
}

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

func (e *EnvT) T() T {
	return e.t
}

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
		e.t.Fatalf("failed to call fixture func: %v", err)
		// return not reachable after Fatalf
		return nil
	}

	return res
}

func (e *EnvT) Cleanup(f func()) {
	e.T().Cleanup(f)
}

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
