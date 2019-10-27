package dns

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/dns.ResolverInterface -o ./resolver_interface_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"net"

	"github.com/gojuno/minimock/v3"
)

// ResolverInterfaceMock implements ResolverInterface
type ResolverInterfaceMock struct {
	t minimock.Tester

	funcLookupIPAddr          func(ctx context.Context, host string) (ia1 []net.IPAddr, err error)
	afterLookupIPAddrCounter  uint64
	beforeLookupIPAddrCounter uint64
	LookupIPAddrMock          mResolverInterfaceMockLookupIPAddr
}

// NewResolverInterfaceMock returns a mock for ResolverInterface
func NewResolverInterfaceMock(t minimock.Tester) *ResolverInterfaceMock {
	m := &ResolverInterfaceMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.LookupIPAddrMock = mResolverInterfaceMockLookupIPAddr{mock: m}

	return m
}

type mResolverInterfaceMockLookupIPAddr struct {
	mock               *ResolverInterfaceMock
	defaultExpectation *ResolverInterfaceMockLookupIPAddrExpectation
	expectations       []*ResolverInterfaceMockLookupIPAddrExpectation
}

// ResolverInterfaceMockLookupIPAddrExpectation specifies expectation struct of the ResolverInterface.LookupIPAddr
type ResolverInterfaceMockLookupIPAddrExpectation struct {
	mock    *ResolverInterfaceMock
	params  *ResolverInterfaceMockLookupIPAddrParams
	results *ResolverInterfaceMockLookupIPAddrResults
	Counter uint64
}

// ResolverInterfaceMockLookupIPAddrParams contains parameters of the ResolverInterface.LookupIPAddr
type ResolverInterfaceMockLookupIPAddrParams struct {
	ctx  context.Context
	host string
}

// ResolverInterfaceMockLookupIPAddrResults contains results of the ResolverInterface.LookupIPAddr
type ResolverInterfaceMockLookupIPAddrResults struct {
	ia1 []net.IPAddr
	err error
}

// Expect sets up expected params for ResolverInterface.LookupIPAddr
func (m *mResolverInterfaceMockLookupIPAddr) Expect(ctx context.Context, host string) *mResolverInterfaceMockLookupIPAddr {
	if m.mock.funcLookupIPAddr != nil {
		m.mock.t.Fatalf("ResolverInterfaceMock.LookupIPAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ResolverInterfaceMockLookupIPAddrExpectation{}
	}

	m.defaultExpectation.params = &ResolverInterfaceMockLookupIPAddrParams{ctx, host}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by ResolverInterface.LookupIPAddr
func (m *mResolverInterfaceMockLookupIPAddr) Return(ia1 []net.IPAddr, err error) *ResolverInterfaceMock {
	if m.mock.funcLookupIPAddr != nil {
		m.mock.t.Fatalf("ResolverInterfaceMock.LookupIPAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ResolverInterfaceMockLookupIPAddrExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ResolverInterfaceMockLookupIPAddrResults{ia1, err}
	return m.mock
}

//Set uses given function f to mock the ResolverInterface.LookupIPAddr method
func (m *mResolverInterfaceMockLookupIPAddr) Set(f func(ctx context.Context, host string) (ia1 []net.IPAddr, err error)) *ResolverInterfaceMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the ResolverInterface.LookupIPAddr method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the ResolverInterface.LookupIPAddr method")
	}

	m.mock.funcLookupIPAddr = f
	return m.mock
}

// When sets expectation for the ResolverInterface.LookupIPAddr which will trigger the result defined by the following
// Then helper
func (m *mResolverInterfaceMockLookupIPAddr) When(ctx context.Context, host string) *ResolverInterfaceMockLookupIPAddrExpectation {
	if m.mock.funcLookupIPAddr != nil {
		m.mock.t.Fatalf("ResolverInterfaceMock.LookupIPAddr mock is already set by Set")
	}

	expectation := &ResolverInterfaceMockLookupIPAddrExpectation{
		mock:   m.mock,
		params: &ResolverInterfaceMockLookupIPAddrParams{ctx, host},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up ResolverInterface.LookupIPAddr return parameters for the expectation previously defined by the When method
func (e *ResolverInterfaceMockLookupIPAddrExpectation) Then(ia1 []net.IPAddr, err error) *ResolverInterfaceMock {
	e.results = &ResolverInterfaceMockLookupIPAddrResults{ia1, err}
	return e.mock
}

// LookupIPAddr implements ResolverInterface
func (m *ResolverInterfaceMock) LookupIPAddr(ctx context.Context, host string) (ia1 []net.IPAddr, err error) {
	atomic.AddUint64(&m.beforeLookupIPAddrCounter, 1)
	defer atomic.AddUint64(&m.afterLookupIPAddrCounter, 1)

	for _, e := range m.LookupIPAddrMock.expectations {
		if minimock.Equal(*e.params, ResolverInterfaceMockLookupIPAddrParams{ctx, host}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.ia1, e.results.err
		}
	}

	if m.LookupIPAddrMock.defaultExpectation != nil {
		atomic.AddUint64(&m.LookupIPAddrMock.defaultExpectation.Counter, 1)
		want := m.LookupIPAddrMock.defaultExpectation.params
		got := ResolverInterfaceMockLookupIPAddrParams{ctx, host}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ResolverInterfaceMock.LookupIPAddr got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.LookupIPAddrMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ResolverInterfaceMock.LookupIPAddr")
		}
		return (*results).ia1, (*results).err
	}
	if m.funcLookupIPAddr != nil {
		return m.funcLookupIPAddr(ctx, host)
	}
	m.t.Fatalf("Unexpected call to ResolverInterfaceMock.LookupIPAddr. %v %v", ctx, host)
	return
}

// LookupIPAddrAfterCounter returns a count of finished ResolverInterfaceMock.LookupIPAddr invocations
func (m *ResolverInterfaceMock) LookupIPAddrAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterLookupIPAddrCounter)
}

// LookupIPAddrBeforeCounter returns a count of ResolverInterfaceMock.LookupIPAddr invocations
func (m *ResolverInterfaceMock) LookupIPAddrBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeLookupIPAddrCounter)
}

// MinimockLookupIPAddrDone returns true if the count of the LookupIPAddr invocations corresponds
// the number of defined expectations
func (m *ResolverInterfaceMock) MinimockLookupIPAddrDone() bool {
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
func (m *ResolverInterfaceMock) MinimockLookupIPAddrInspect() {
	for _, e := range m.LookupIPAddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ResolverInterfaceMock.LookupIPAddr with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LookupIPAddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterLookupIPAddrCounter) < 1 {
		m.t.Errorf("Expected call to ResolverInterfaceMock.LookupIPAddr with params: %#v", *m.LookupIPAddrMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLookupIPAddr != nil && atomic.LoadUint64(&m.afterLookupIPAddrCounter) < 1 {
		m.t.Error("Expected call to ResolverInterfaceMock.LookupIPAddr")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ResolverInterfaceMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockLookupIPAddrInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ResolverInterfaceMock) MinimockWait(timeout time.Duration) {
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

func (m *ResolverInterfaceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockLookupIPAddrDone()
}
