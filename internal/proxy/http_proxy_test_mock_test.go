package proxy

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/proxy.HTTPProxyTest -o ./http_proxy_test_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"net/http"

	"github.com/gojuno/minimock"
)

// HTTPProxyTestMock implements HTTPProxyTest
type HTTPProxyTestMock struct {
	t minimock.Tester

	funcGetContext          func(req *http.Request) (c1 context.Context, err error)
	afterGetContextCounter  uint64
	beforeGetContextCounter uint64
	GetContextMock          mHTTPProxyTestMockGetContext

	funcHandleHTTPValidation          func(w http.ResponseWriter, r *http.Request) (b1 bool)
	afterHandleHTTPValidationCounter  uint64
	beforeHandleHTTPValidationCounter uint64
	HandleHTTPValidationMock          mHTTPProxyTestMockHandleHTTPValidation
}

// NewHTTPProxyTestMock returns a mock for HTTPProxyTest
func NewHTTPProxyTestMock(t minimock.Tester) *HTTPProxyTestMock {
	m := &HTTPProxyTestMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.GetContextMock = mHTTPProxyTestMockGetContext{mock: m}
	m.HandleHTTPValidationMock = mHTTPProxyTestMockHandleHTTPValidation{mock: m}

	return m
}

type mHTTPProxyTestMockGetContext struct {
	mock               *HTTPProxyTestMock
	defaultExpectation *HTTPProxyTestMockGetContextExpectation
	expectations       []*HTTPProxyTestMockGetContextExpectation
}

// HTTPProxyTestMockGetContextExpectation specifies expectation struct of the HTTPProxyTest.GetContext
type HTTPProxyTestMockGetContextExpectation struct {
	mock    *HTTPProxyTestMock
	params  *HTTPProxyTestMockGetContextParams
	results *HTTPProxyTestMockGetContextResults
	Counter uint64
}

// HTTPProxyTestMockGetContextParams contains parameters of the HTTPProxyTest.GetContext
type HTTPProxyTestMockGetContextParams struct {
	req *http.Request
}

// HTTPProxyTestMockGetContextResults contains results of the HTTPProxyTest.GetContext
type HTTPProxyTestMockGetContextResults struct {
	c1  context.Context
	err error
}

// Expect sets up expected params for HTTPProxyTest.GetContext
func (m *mHTTPProxyTestMockGetContext) Expect(req *http.Request) *mHTTPProxyTestMockGetContext {
	if m.mock.funcGetContext != nil {
		m.mock.t.Fatalf("HTTPProxyTestMock.GetContext mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HTTPProxyTestMockGetContextExpectation{}
	}

	m.defaultExpectation.params = &HTTPProxyTestMockGetContextParams{req}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by HTTPProxyTest.GetContext
func (m *mHTTPProxyTestMockGetContext) Return(c1 context.Context, err error) *HTTPProxyTestMock {
	if m.mock.funcGetContext != nil {
		m.mock.t.Fatalf("HTTPProxyTestMock.GetContext mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HTTPProxyTestMockGetContextExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &HTTPProxyTestMockGetContextResults{c1, err}
	return m.mock
}

//Set uses given function f to mock the HTTPProxyTest.GetContext method
func (m *mHTTPProxyTestMockGetContext) Set(f func(req *http.Request) (c1 context.Context, err error)) *HTTPProxyTestMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the HTTPProxyTest.GetContext method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the HTTPProxyTest.GetContext method")
	}

	m.mock.funcGetContext = f
	return m.mock
}

// When sets expectation for the HTTPProxyTest.GetContext which will trigger the result defined by the following
// Then helper
func (m *mHTTPProxyTestMockGetContext) When(req *http.Request) *HTTPProxyTestMockGetContextExpectation {
	if m.mock.funcGetContext != nil {
		m.mock.t.Fatalf("HTTPProxyTestMock.GetContext mock is already set by Set")
	}

	expectation := &HTTPProxyTestMockGetContextExpectation{
		mock:   m.mock,
		params: &HTTPProxyTestMockGetContextParams{req},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up HTTPProxyTest.GetContext return parameters for the expectation previously defined by the When method
func (e *HTTPProxyTestMockGetContextExpectation) Then(c1 context.Context, err error) *HTTPProxyTestMock {
	e.results = &HTTPProxyTestMockGetContextResults{c1, err}
	return e.mock
}

// GetContext implements HTTPProxyTest
func (m *HTTPProxyTestMock) GetContext(req *http.Request) (c1 context.Context, err error) {
	atomic.AddUint64(&m.beforeGetContextCounter, 1)
	defer atomic.AddUint64(&m.afterGetContextCounter, 1)

	for _, e := range m.GetContextMock.expectations {
		if minimock.Equal(*e.params, HTTPProxyTestMockGetContextParams{req}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.c1, e.results.err
		}
	}

	if m.GetContextMock.defaultExpectation != nil {
		atomic.AddUint64(&m.GetContextMock.defaultExpectation.Counter, 1)
		want := m.GetContextMock.defaultExpectation.params
		got := HTTPProxyTestMockGetContextParams{req}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("HTTPProxyTestMock.GetContext got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetContextMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the HTTPProxyTestMock.GetContext")
		}
		return (*results).c1, (*results).err
	}
	if m.funcGetContext != nil {
		return m.funcGetContext(req)
	}
	m.t.Fatalf("Unexpected call to HTTPProxyTestMock.GetContext. %v", req)
	return
}

// GetContextAfterCounter returns a count of finished HTTPProxyTestMock.GetContext invocations
func (m *HTTPProxyTestMock) GetContextAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterGetContextCounter)
}

// GetContextBeforeCounter returns a count of HTTPProxyTestMock.GetContext invocations
func (m *HTTPProxyTestMock) GetContextBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeGetContextCounter)
}

// MinimockGetContextDone returns true if the count of the GetContext invocations corresponds
// the number of defined expectations
func (m *HTTPProxyTestMock) MinimockGetContextDone() bool {
	for _, e := range m.GetContextMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetContextMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetContextCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetContext != nil && atomic.LoadUint64(&m.afterGetContextCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetContextInspect logs each unmet expectation
func (m *HTTPProxyTestMock) MinimockGetContextInspect() {
	for _, e := range m.GetContextMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to HTTPProxyTestMock.GetContext with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetContextMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetContextCounter) < 1 {
		m.t.Errorf("Expected call to HTTPProxyTestMock.GetContext with params: %#v", *m.GetContextMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetContext != nil && atomic.LoadUint64(&m.afterGetContextCounter) < 1 {
		m.t.Error("Expected call to HTTPProxyTestMock.GetContext")
	}
}

type mHTTPProxyTestMockHandleHTTPValidation struct {
	mock               *HTTPProxyTestMock
	defaultExpectation *HTTPProxyTestMockHandleHTTPValidationExpectation
	expectations       []*HTTPProxyTestMockHandleHTTPValidationExpectation
}

// HTTPProxyTestMockHandleHTTPValidationExpectation specifies expectation struct of the HTTPProxyTest.HandleHTTPValidation
type HTTPProxyTestMockHandleHTTPValidationExpectation struct {
	mock    *HTTPProxyTestMock
	params  *HTTPProxyTestMockHandleHTTPValidationParams
	results *HTTPProxyTestMockHandleHTTPValidationResults
	Counter uint64
}

// HTTPProxyTestMockHandleHTTPValidationParams contains parameters of the HTTPProxyTest.HandleHTTPValidation
type HTTPProxyTestMockHandleHTTPValidationParams struct {
	w http.ResponseWriter
	r *http.Request
}

// HTTPProxyTestMockHandleHTTPValidationResults contains results of the HTTPProxyTest.HandleHTTPValidation
type HTTPProxyTestMockHandleHTTPValidationResults struct {
	b1 bool
}

// Expect sets up expected params for HTTPProxyTest.HandleHTTPValidation
func (m *mHTTPProxyTestMockHandleHTTPValidation) Expect(w http.ResponseWriter, r *http.Request) *mHTTPProxyTestMockHandleHTTPValidation {
	if m.mock.funcHandleHTTPValidation != nil {
		m.mock.t.Fatalf("HTTPProxyTestMock.HandleHTTPValidation mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HTTPProxyTestMockHandleHTTPValidationExpectation{}
	}

	m.defaultExpectation.params = &HTTPProxyTestMockHandleHTTPValidationParams{w, r}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by HTTPProxyTest.HandleHTTPValidation
func (m *mHTTPProxyTestMockHandleHTTPValidation) Return(b1 bool) *HTTPProxyTestMock {
	if m.mock.funcHandleHTTPValidation != nil {
		m.mock.t.Fatalf("HTTPProxyTestMock.HandleHTTPValidation mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HTTPProxyTestMockHandleHTTPValidationExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &HTTPProxyTestMockHandleHTTPValidationResults{b1}
	return m.mock
}

//Set uses given function f to mock the HTTPProxyTest.HandleHTTPValidation method
func (m *mHTTPProxyTestMockHandleHTTPValidation) Set(f func(w http.ResponseWriter, r *http.Request) (b1 bool)) *HTTPProxyTestMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the HTTPProxyTest.HandleHTTPValidation method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the HTTPProxyTest.HandleHTTPValidation method")
	}

	m.mock.funcHandleHTTPValidation = f
	return m.mock
}

// When sets expectation for the HTTPProxyTest.HandleHTTPValidation which will trigger the result defined by the following
// Then helper
func (m *mHTTPProxyTestMockHandleHTTPValidation) When(w http.ResponseWriter, r *http.Request) *HTTPProxyTestMockHandleHTTPValidationExpectation {
	if m.mock.funcHandleHTTPValidation != nil {
		m.mock.t.Fatalf("HTTPProxyTestMock.HandleHTTPValidation mock is already set by Set")
	}

	expectation := &HTTPProxyTestMockHandleHTTPValidationExpectation{
		mock:   m.mock,
		params: &HTTPProxyTestMockHandleHTTPValidationParams{w, r},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up HTTPProxyTest.HandleHTTPValidation return parameters for the expectation previously defined by the When method
func (e *HTTPProxyTestMockHandleHTTPValidationExpectation) Then(b1 bool) *HTTPProxyTestMock {
	e.results = &HTTPProxyTestMockHandleHTTPValidationResults{b1}
	return e.mock
}

// HandleHTTPValidation implements HTTPProxyTest
func (m *HTTPProxyTestMock) HandleHTTPValidation(w http.ResponseWriter, r *http.Request) (b1 bool) {
	atomic.AddUint64(&m.beforeHandleHTTPValidationCounter, 1)
	defer atomic.AddUint64(&m.afterHandleHTTPValidationCounter, 1)

	for _, e := range m.HandleHTTPValidationMock.expectations {
		if minimock.Equal(*e.params, HTTPProxyTestMockHandleHTTPValidationParams{w, r}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.b1
		}
	}

	if m.HandleHTTPValidationMock.defaultExpectation != nil {
		atomic.AddUint64(&m.HandleHTTPValidationMock.defaultExpectation.Counter, 1)
		want := m.HandleHTTPValidationMock.defaultExpectation.params
		got := HTTPProxyTestMockHandleHTTPValidationParams{w, r}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("HTTPProxyTestMock.HandleHTTPValidation got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.HandleHTTPValidationMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the HTTPProxyTestMock.HandleHTTPValidation")
		}
		return (*results).b1
	}
	if m.funcHandleHTTPValidation != nil {
		return m.funcHandleHTTPValidation(w, r)
	}
	m.t.Fatalf("Unexpected call to HTTPProxyTestMock.HandleHTTPValidation. %v %v", w, r)
	return
}

// HandleHTTPValidationAfterCounter returns a count of finished HTTPProxyTestMock.HandleHTTPValidation invocations
func (m *HTTPProxyTestMock) HandleHTTPValidationAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterHandleHTTPValidationCounter)
}

// HandleHTTPValidationBeforeCounter returns a count of HTTPProxyTestMock.HandleHTTPValidation invocations
func (m *HTTPProxyTestMock) HandleHTTPValidationBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeHandleHTTPValidationCounter)
}

// MinimockHandleHTTPValidationDone returns true if the count of the HandleHTTPValidation invocations corresponds
// the number of defined expectations
func (m *HTTPProxyTestMock) MinimockHandleHTTPValidationDone() bool {
	for _, e := range m.HandleHTTPValidationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HandleHTTPValidationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterHandleHTTPValidationCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHandleHTTPValidation != nil && atomic.LoadUint64(&m.afterHandleHTTPValidationCounter) < 1 {
		return false
	}
	return true
}

// MinimockHandleHTTPValidationInspect logs each unmet expectation
func (m *HTTPProxyTestMock) MinimockHandleHTTPValidationInspect() {
	for _, e := range m.HandleHTTPValidationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to HTTPProxyTestMock.HandleHTTPValidation with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HandleHTTPValidationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterHandleHTTPValidationCounter) < 1 {
		m.t.Errorf("Expected call to HTTPProxyTestMock.HandleHTTPValidation with params: %#v", *m.HandleHTTPValidationMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHandleHTTPValidation != nil && atomic.LoadUint64(&m.afterHandleHTTPValidationCounter) < 1 {
		m.t.Error("Expected call to HTTPProxyTestMock.HandleHTTPValidation")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *HTTPProxyTestMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetContextInspect()

		m.MinimockHandleHTTPValidationInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *HTTPProxyTestMock) MinimockWait(timeout time.Duration) {
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

func (m *HTTPProxyTestMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetContextDone() &&
		m.MinimockHandleHTTPValidationDone()
}
