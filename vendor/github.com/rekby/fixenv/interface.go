package fixenv

// Env - fixture cache engine.
type Env interface {
	// T - return t object of current test/benchmark.
	T() T

	// Cache cache result of f calls
	// f call exactly once for every combination of scope and params
	// params must be json serializable (deserialize not need)
	Cache(params interface{}, opt *FixtureOptions, f FixtureCallbackFunc) interface{}

	// Cleanup add callback cleanup function
	// f called while env clean
	Cleanup(f func())
}

type CacheScope int

const (
	// ScopeTest mean fixture function with same parameters called once per every test and subtests. Default value.
	// Second and more calls will use cached value.
	ScopeTest CacheScope = iota

	// ScopePackage mean fixture function with same parameters called once per package
	// for use the scope with TearDown function developer must initialize global handler and cleaner at TestMain.
	ScopePackage CacheScope = iota

	// ScopeTestAndSubtests mean fixture cached for top level test and subtests
	ScopeTestAndSubtests CacheScope = iota
)

// FixtureCallbackFunc - function, which result can cached
// res - result for cache.
// if err not nil - T().Fatalf() will called with error message
// if res exit without return (panic, GoExit, t.FailNow, ...)
// then cache error about unexpected exit
type FixtureCallbackFunc func() (res interface{}, err error)

type FixtureCleanupFunc func()

type FixtureOptions struct {
	// Scope for cache result
	Scope CacheScope

	// CleanupFunc if not nil - called for cleanup fixture results
	// it called exactly once for every succesully call fixture
	CleanupFunc FixtureCleanupFunc
}

// T is subtype of testing.TB
type T interface {
	Cleanup(func())
	Fatalf(format string, args ...interface{})
	Name() string
}
