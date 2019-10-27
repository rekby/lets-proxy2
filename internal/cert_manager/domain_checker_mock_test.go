package cert_manager

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cert_manager.DomainChecker -o ./domain_checker_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"github.com/gojuno/minimock/v3"
)

// DomainCheckerMock implements DomainChecker
type DomainCheckerMock struct {
	t minimock.Tester

	funcIsDomainAllowed          func(ctx context.Context, domain string) (b1 bool, err error)
	afterIsDomainAllowedCounter  uint64
	beforeIsDomainAllowedCounter uint64
	IsDomainAllowedMock          mDomainCheckerMockIsDomainAllowed
}

// NewDomainCheckerMock returns a mock for DomainChecker
func NewDomainCheckerMock(t minimock.Tester) *DomainCheckerMock {
	m := &DomainCheckerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.IsDomainAllowedMock = mDomainCheckerMockIsDomainAllowed{mock: m}

	return m
}

type mDomainCheckerMockIsDomainAllowed struct {
	mock               *DomainCheckerMock
	defaultExpectation *DomainCheckerMockIsDomainAllowedExpectation
	expectations       []*DomainCheckerMockIsDomainAllowedExpectation
}

// DomainCheckerMockIsDomainAllowedExpectation specifies expectation struct of the DomainChecker.IsDomainAllowed
type DomainCheckerMockIsDomainAllowedExpectation struct {
	mock    *DomainCheckerMock
	params  *DomainCheckerMockIsDomainAllowedParams
	results *DomainCheckerMockIsDomainAllowedResults
	Counter uint64
}

// DomainCheckerMockIsDomainAllowedParams contains parameters of the DomainChecker.IsDomainAllowed
type DomainCheckerMockIsDomainAllowedParams struct {
	ctx    context.Context
	domain string
}

// DomainCheckerMockIsDomainAllowedResults contains results of the DomainChecker.IsDomainAllowed
type DomainCheckerMockIsDomainAllowedResults struct {
	b1  bool
	err error
}

// Expect sets up expected params for DomainChecker.IsDomainAllowed
func (m *mDomainCheckerMockIsDomainAllowed) Expect(ctx context.Context, domain string) *mDomainCheckerMockIsDomainAllowed {
	if m.mock.funcIsDomainAllowed != nil {
		m.mock.t.Fatalf("DomainCheckerMock.IsDomainAllowed mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &DomainCheckerMockIsDomainAllowedExpectation{}
	}

	m.defaultExpectation.params = &DomainCheckerMockIsDomainAllowedParams{ctx, domain}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by DomainChecker.IsDomainAllowed
func (m *mDomainCheckerMockIsDomainAllowed) Return(b1 bool, err error) *DomainCheckerMock {
	if m.mock.funcIsDomainAllowed != nil {
		m.mock.t.Fatalf("DomainCheckerMock.IsDomainAllowed mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &DomainCheckerMockIsDomainAllowedExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &DomainCheckerMockIsDomainAllowedResults{b1, err}
	return m.mock
}

//Set uses given function f to mock the DomainChecker.IsDomainAllowed method
func (m *mDomainCheckerMockIsDomainAllowed) Set(f func(ctx context.Context, domain string) (b1 bool, err error)) *DomainCheckerMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the DomainChecker.IsDomainAllowed method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the DomainChecker.IsDomainAllowed method")
	}

	m.mock.funcIsDomainAllowed = f
	return m.mock
}

// When sets expectation for the DomainChecker.IsDomainAllowed which will trigger the result defined by the following
// Then helper
func (m *mDomainCheckerMockIsDomainAllowed) When(ctx context.Context, domain string) *DomainCheckerMockIsDomainAllowedExpectation {
	if m.mock.funcIsDomainAllowed != nil {
		m.mock.t.Fatalf("DomainCheckerMock.IsDomainAllowed mock is already set by Set")
	}

	expectation := &DomainCheckerMockIsDomainAllowedExpectation{
		mock:   m.mock,
		params: &DomainCheckerMockIsDomainAllowedParams{ctx, domain},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up DomainChecker.IsDomainAllowed return parameters for the expectation previously defined by the When method
func (e *DomainCheckerMockIsDomainAllowedExpectation) Then(b1 bool, err error) *DomainCheckerMock {
	e.results = &DomainCheckerMockIsDomainAllowedResults{b1, err}
	return e.mock
}

// IsDomainAllowed implements DomainChecker
func (m *DomainCheckerMock) IsDomainAllowed(ctx context.Context, domain string) (b1 bool, err error) {
	atomic.AddUint64(&m.beforeIsDomainAllowedCounter, 1)
	defer atomic.AddUint64(&m.afterIsDomainAllowedCounter, 1)

	for _, e := range m.IsDomainAllowedMock.expectations {
		if minimock.Equal(*e.params, DomainCheckerMockIsDomainAllowedParams{ctx, domain}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.b1, e.results.err
		}
	}

	if m.IsDomainAllowedMock.defaultExpectation != nil {
		atomic.AddUint64(&m.IsDomainAllowedMock.defaultExpectation.Counter, 1)
		want := m.IsDomainAllowedMock.defaultExpectation.params
		got := DomainCheckerMockIsDomainAllowedParams{ctx, domain}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("DomainCheckerMock.IsDomainAllowed got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.IsDomainAllowedMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the DomainCheckerMock.IsDomainAllowed")
		}
		return (*results).b1, (*results).err
	}
	if m.funcIsDomainAllowed != nil {
		return m.funcIsDomainAllowed(ctx, domain)
	}
	m.t.Fatalf("Unexpected call to DomainCheckerMock.IsDomainAllowed. %v %v", ctx, domain)
	return
}

// IsDomainAllowedAfterCounter returns a count of finished DomainCheckerMock.IsDomainAllowed invocations
func (m *DomainCheckerMock) IsDomainAllowedAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterIsDomainAllowedCounter)
}

// IsDomainAllowedBeforeCounter returns a count of DomainCheckerMock.IsDomainAllowed invocations
func (m *DomainCheckerMock) IsDomainAllowedBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeIsDomainAllowedCounter)
}

// MinimockIsDomainAllowedDone returns true if the count of the IsDomainAllowed invocations corresponds
// the number of defined expectations
func (m *DomainCheckerMock) MinimockIsDomainAllowedDone() bool {
	for _, e := range m.IsDomainAllowedMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.IsDomainAllowedMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterIsDomainAllowedCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcIsDomainAllowed != nil && atomic.LoadUint64(&m.afterIsDomainAllowedCounter) < 1 {
		return false
	}
	return true
}

// MinimockIsDomainAllowedInspect logs each unmet expectation
func (m *DomainCheckerMock) MinimockIsDomainAllowedInspect() {
	for _, e := range m.IsDomainAllowedMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to DomainCheckerMock.IsDomainAllowed with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.IsDomainAllowedMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterIsDomainAllowedCounter) < 1 {
		m.t.Errorf("Expected call to DomainCheckerMock.IsDomainAllowed with params: %#v", *m.IsDomainAllowedMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcIsDomainAllowed != nil && atomic.LoadUint64(&m.afterIsDomainAllowedCounter) < 1 {
		m.t.Error("Expected call to DomainCheckerMock.IsDomainAllowed")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *DomainCheckerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockIsDomainAllowedInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *DomainCheckerMock) MinimockWait(timeout time.Duration) {
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

func (m *DomainCheckerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockIsDomainAllowedDone()
}
