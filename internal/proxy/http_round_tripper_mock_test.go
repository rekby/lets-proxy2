package proxy

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i net/http.RoundTripper -o ./http_round_tripper_mock_test.go

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock/v3"
)

// RoundTripperMock implements http.RoundTripper
type RoundTripperMock struct {
	t minimock.Tester

	funcRoundTrip          func(rp1 *http.Request) (rp2 *http.Response, err error)
	afterRoundTripCounter  uint64
	beforeRoundTripCounter uint64
	RoundTripMock          mRoundTripperMockRoundTrip
}

// NewRoundTripperMock returns a mock for http.RoundTripper
func NewRoundTripperMock(t minimock.Tester) *RoundTripperMock {
	m := &RoundTripperMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.RoundTripMock = mRoundTripperMockRoundTrip{mock: m}

	return m
}

type mRoundTripperMockRoundTrip struct {
	mock               *RoundTripperMock
	defaultExpectation *RoundTripperMockRoundTripExpectation
	expectations       []*RoundTripperMockRoundTripExpectation
}

// RoundTripperMockRoundTripExpectation specifies expectation struct of the RoundTripper.RoundTrip
type RoundTripperMockRoundTripExpectation struct {
	mock    *RoundTripperMock
	params  *RoundTripperMockRoundTripParams
	results *RoundTripperMockRoundTripResults
	Counter uint64
}

// RoundTripperMockRoundTripParams contains parameters of the RoundTripper.RoundTrip
type RoundTripperMockRoundTripParams struct {
	rp1 *http.Request
}

// RoundTripperMockRoundTripResults contains results of the RoundTripper.RoundTrip
type RoundTripperMockRoundTripResults struct {
	rp2 *http.Response
	err error
}

// Expect sets up expected params for RoundTripper.RoundTrip
func (m *mRoundTripperMockRoundTrip) Expect(rp1 *http.Request) *mRoundTripperMockRoundTrip {
	if m.mock.funcRoundTrip != nil {
		m.mock.t.Fatalf("RoundTripperMock.RoundTrip mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &RoundTripperMockRoundTripExpectation{}
	}

	m.defaultExpectation.params = &RoundTripperMockRoundTripParams{rp1}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by RoundTripper.RoundTrip
func (m *mRoundTripperMockRoundTrip) Return(rp2 *http.Response, err error) *RoundTripperMock {
	if m.mock.funcRoundTrip != nil {
		m.mock.t.Fatalf("RoundTripperMock.RoundTrip mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &RoundTripperMockRoundTripExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &RoundTripperMockRoundTripResults{rp2, err}
	return m.mock
}

//Set uses given function f to mock the RoundTripper.RoundTrip method
func (m *mRoundTripperMockRoundTrip) Set(f func(rp1 *http.Request) (rp2 *http.Response, err error)) *RoundTripperMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the RoundTripper.RoundTrip method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the RoundTripper.RoundTrip method")
	}

	m.mock.funcRoundTrip = f
	return m.mock
}

// When sets expectation for the RoundTripper.RoundTrip which will trigger the result defined by the following
// Then helper
func (m *mRoundTripperMockRoundTrip) When(rp1 *http.Request) *RoundTripperMockRoundTripExpectation {
	if m.mock.funcRoundTrip != nil {
		m.mock.t.Fatalf("RoundTripperMock.RoundTrip mock is already set by Set")
	}

	expectation := &RoundTripperMockRoundTripExpectation{
		mock:   m.mock,
		params: &RoundTripperMockRoundTripParams{rp1},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up RoundTripper.RoundTrip return parameters for the expectation previously defined by the When method
func (e *RoundTripperMockRoundTripExpectation) Then(rp2 *http.Response, err error) *RoundTripperMock {
	e.results = &RoundTripperMockRoundTripResults{rp2, err}
	return e.mock
}

// RoundTrip implements http.RoundTripper
func (m *RoundTripperMock) RoundTrip(rp1 *http.Request) (rp2 *http.Response, err error) {
	atomic.AddUint64(&m.beforeRoundTripCounter, 1)
	defer atomic.AddUint64(&m.afterRoundTripCounter, 1)

	for _, e := range m.RoundTripMock.expectations {
		if minimock.Equal(*e.params, RoundTripperMockRoundTripParams{rp1}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.rp2, e.results.err
		}
	}

	if m.RoundTripMock.defaultExpectation != nil {
		atomic.AddUint64(&m.RoundTripMock.defaultExpectation.Counter, 1)
		want := m.RoundTripMock.defaultExpectation.params
		got := RoundTripperMockRoundTripParams{rp1}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("RoundTripperMock.RoundTrip got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.RoundTripMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the RoundTripperMock.RoundTrip")
		}
		return (*results).rp2, (*results).err
	}
	if m.funcRoundTrip != nil {
		return m.funcRoundTrip(rp1)
	}
	m.t.Fatalf("Unexpected call to RoundTripperMock.RoundTrip. %v", rp1)
	return
}

// RoundTripAfterCounter returns a count of finished RoundTripperMock.RoundTrip invocations
func (m *RoundTripperMock) RoundTripAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterRoundTripCounter)
}

// RoundTripBeforeCounter returns a count of RoundTripperMock.RoundTrip invocations
func (m *RoundTripperMock) RoundTripBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeRoundTripCounter)
}

// MinimockRoundTripDone returns true if the count of the RoundTrip invocations corresponds
// the number of defined expectations
func (m *RoundTripperMock) MinimockRoundTripDone() bool {
	for _, e := range m.RoundTripMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RoundTripMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRoundTripCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRoundTrip != nil && atomic.LoadUint64(&m.afterRoundTripCounter) < 1 {
		return false
	}
	return true
}

// MinimockRoundTripInspect logs each unmet expectation
func (m *RoundTripperMock) MinimockRoundTripInspect() {
	for _, e := range m.RoundTripMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RoundTripperMock.RoundTrip with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RoundTripMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRoundTripCounter) < 1 {
		m.t.Errorf("Expected call to RoundTripperMock.RoundTrip with params: %#v", *m.RoundTripMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRoundTrip != nil && atomic.LoadUint64(&m.afterRoundTripCounter) < 1 {
		m.t.Error("Expected call to RoundTripperMock.RoundTrip")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *RoundTripperMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockRoundTripInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *RoundTripperMock) MinimockWait(timeout time.Duration) {
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

func (m *RoundTripperMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockRoundTripDone()
}
