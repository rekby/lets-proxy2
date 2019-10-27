package domain_checker

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/domain_checker.Resolver -o ./resolver_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"net"

	"github.com/gojuno/minimock/v3"
)

// ResolverMock implements Resolver
type ResolverMock struct {
	t minimock.Tester

	funcLookupIPAddr          func(ctx context.Context, host string) (ia1 []net.IPAddr, err error)
	afterLookupIPAddrCounter  uint64
	beforeLookupIPAddrCounter uint64
	LookupIPAddrMock          mResolverMockLookupIPAddr
}

// NewResolverMock returns a mock for Resolver
func NewResolverMock(t minimock.Tester) *ResolverMock {
	m := &ResolverMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.LookupIPAddrMock = mResolverMockLookupIPAddr{mock: m}

	return m
}

type mResolverMockLookupIPAddr struct {
	mock               *ResolverMock
	defaultExpectation *ResolverMockLookupIPAddrExpectation
	expectations       []*ResolverMockLookupIPAddrExpectation
}

// ResolverMockLookupIPAddrExpectation specifies expectation struct of the Resolver.LookupIPAddr
type ResolverMockLookupIPAddrExpectation struct {
	mock    *ResolverMock
	params  *ResolverMockLookupIPAddrParams
	results *ResolverMockLookupIPAddrResults
	Counter uint64
}

// ResolverMockLookupIPAddrParams contains parameters of the Resolver.LookupIPAddr
type ResolverMockLookupIPAddrParams struct {
	ctx  context.Context
	host string
}

// ResolverMockLookupIPAddrResults contains results of the Resolver.LookupIPAddr
type ResolverMockLookupIPAddrResults struct {
	ia1 []net.IPAddr
	err error
}

// Expect sets up expected params for Resolver.LookupIPAddr
func (m *mResolverMockLookupIPAddr) Expect(ctx context.Context, host string) *mResolverMockLookupIPAddr {
	if m.mock.funcLookupIPAddr != nil {
		m.mock.t.Fatalf("ResolverMock.LookupIPAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ResolverMockLookupIPAddrExpectation{}
	}

	m.defaultExpectation.params = &ResolverMockLookupIPAddrParams{ctx, host}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Resolver.LookupIPAddr
func (m *mResolverMockLookupIPAddr) Return(ia1 []net.IPAddr, err error) *ResolverMock {
	if m.mock.funcLookupIPAddr != nil {
		m.mock.t.Fatalf("ResolverMock.LookupIPAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ResolverMockLookupIPAddrExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ResolverMockLookupIPAddrResults{ia1, err}
	return m.mock
}

//Set uses given function f to mock the Resolver.LookupIPAddr method
func (m *mResolverMockLookupIPAddr) Set(f func(ctx context.Context, host string) (ia1 []net.IPAddr, err error)) *ResolverMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Resolver.LookupIPAddr method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Resolver.LookupIPAddr method")
	}

	m.mock.funcLookupIPAddr = f
	return m.mock
}

// When sets expectation for the Resolver.LookupIPAddr which will trigger the result defined by the following
// Then helper
func (m *mResolverMockLookupIPAddr) When(ctx context.Context, host string) *ResolverMockLookupIPAddrExpectation {
	if m.mock.funcLookupIPAddr != nil {
		m.mock.t.Fatalf("ResolverMock.LookupIPAddr mock is already set by Set")
	}

	expectation := &ResolverMockLookupIPAddrExpectation{
		mock:   m.mock,
		params: &ResolverMockLookupIPAddrParams{ctx, host},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Resolver.LookupIPAddr return parameters for the expectation previously defined by the When method
func (e *ResolverMockLookupIPAddrExpectation) Then(ia1 []net.IPAddr, err error) *ResolverMock {
	e.results = &ResolverMockLookupIPAddrResults{ia1, err}
	return e.mock
}

// LookupIPAddr implements Resolver
func (m *ResolverMock) LookupIPAddr(ctx context.Context, host string) (ia1 []net.IPAddr, err error) {
	atomic.AddUint64(&m.beforeLookupIPAddrCounter, 1)
	defer atomic.AddUint64(&m.afterLookupIPAddrCounter, 1)

	for _, e := range m.LookupIPAddrMock.expectations {
		if minimock.Equal(*e.params, ResolverMockLookupIPAddrParams{ctx, host}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.ia1, e.results.err
		}
	}

	if m.LookupIPAddrMock.defaultExpectation != nil {
		atomic.AddUint64(&m.LookupIPAddrMock.defaultExpectation.Counter, 1)
		want := m.LookupIPAddrMock.defaultExpectation.params
		got := ResolverMockLookupIPAddrParams{ctx, host}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ResolverMock.LookupIPAddr got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.LookupIPAddrMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ResolverMock.LookupIPAddr")
		}
		return (*results).ia1, (*results).err
	}
	if m.funcLookupIPAddr != nil {
		return m.funcLookupIPAddr(ctx, host)
	}
	m.t.Fatalf("Unexpected call to ResolverMock.LookupIPAddr. %v %v", ctx, host)
	return
}

// LookupIPAddrAfterCounter returns a count of finished ResolverMock.LookupIPAddr invocations
func (m *ResolverMock) LookupIPAddrAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterLookupIPAddrCounter)
}

// LookupIPAddrBeforeCounter returns a count of ResolverMock.LookupIPAddr invocations
func (m *ResolverMock) LookupIPAddrBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeLookupIPAddrCounter)
}

// MinimockLookupIPAddrDone returns true if the count of the LookupIPAddr invocations corresponds
// the number of defined expectations
func (m *ResolverMock) MinimockLookupIPAddrDone() bool {
	for _, e := range m.LookupIPAddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LookupIPAddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterLookupIPAddrCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLookupIPAddr != nil && atomic.LoadUint64(&m.afterLookupIPAddrCounter) < 1 {
		return false
	}
	return true
}

// MinimockLookupIPAddrInspect logs each unmet expectation
func (m *ResolverMock) MinimockLookupIPAddrInspect() {
	for _, e := range m.LookupIPAddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ResolverMock.LookupIPAddr with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LookupIPAddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterLookupIPAddrCounter) < 1 {
		m.t.Errorf("Expected call to ResolverMock.LookupIPAddr with params: %#v", *m.LookupIPAddrMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLookupIPAddr != nil && atomic.LoadUint64(&m.afterLookupIPAddrCounter) < 1 {
		m.t.Error("Expected call to ResolverMock.LookupIPAddr")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ResolverMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockLookupIPAddrInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ResolverMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func (m *ResolverMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockLookupIPAddrDone()
}
