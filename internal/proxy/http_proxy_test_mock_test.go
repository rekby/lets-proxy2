package proxy

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/proxy.HttpProxyTest -o ./http_proxy_test_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"net/http"

	"github.com/gojuno/minimock"
)

// HttpProxyTestMock implements HttpProxyTest
type HttpProxyTestMock struct {
	t minimock.Tester

	funcGetContext          func(req *http.Request) (c1 context.Context, err error)
	afterGetContextCounter  uint64
	beforeGetContextCounter uint64
	GetContextMock          mHttpProxyTestMockGetContext

	funcHandleHTTPValidation          func(w http.ResponseWriter, r *http.Request) (b1 bool)
	afterHandleHTTPValidationCounter  uint64
	beforeHandleHTTPValidationCounter uint64
	HandleHTTPValidationMock          mHttpProxyTestMockHandleHTTPValidation
}

// NewHttpProxyTestMock returns a mock for HttpProxyTest
func NewHttpProxyTestMock(t minimock.Tester) *HttpProxyTestMock {
	m := &HttpProxyTestMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.GetContextMock = mHttpProxyTestMockGetContext{mock: m}
	m.HandleHTTPValidationMock = mHttpProxyTestMockHandleHTTPValidation{mock: m}

	return m
}

type mHttpProxyTestMockGetContext struct {
	mock               *HttpProxyTestMock
	defaultExpectation *HttpProxyTestMockGetContextExpectation
	expectations       []*HttpProxyTestMockGetContextExpectation
}

// HttpProxyTestMockGetContextExpectation specifies expectation struct of the HttpProxyTest.GetContext
type HttpProxyTestMockGetContextExpectation struct {
	mock    *HttpProxyTestMock
	params  *HttpProxyTestMockGetContextParams
	results *HttpProxyTestMockGetContextResults
	Counter uint64
}

// HttpProxyTestMockGetContextParams contains parameters of the HttpProxyTest.GetContext
type HttpProxyTestMockGetContextParams struct {
	req *http.Request
}

// HttpProxyTestMockGetContextResults contains results of the HttpProxyTest.GetContext
type HttpProxyTestMockGetContextResults struct {
	c1  context.Context
	err error
}

// Expect sets up expected params for HttpProxyTest.GetContext
func (m *mHttpProxyTestMockGetContext) Expect(req *http.Request) *mHttpProxyTestMockGetContext {
	if m.mock.funcGetContext != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.GetContext mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockGetContextExpectation{}
	}

	m.defaultExpectation.params = &HttpProxyTestMockGetContextParams{req}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by HttpProxyTest.GetContext
func (m *mHttpProxyTestMockGetContext) Return(c1 context.Context, err error) *HttpProxyTestMock {
	if m.mock.funcGetContext != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.GetContext mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockGetContextExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &HttpProxyTestMockGetContextResults{c1, err}
	return m.mock
}

//Set uses given function f to mock the HttpProxyTest.GetContext method
func (m *mHttpProxyTestMockGetContext) Set(f func(req *http.Request) (c1 context.Context, err error)) *HttpProxyTestMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the HttpProxyTest.GetContext method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the HttpProxyTest.GetContext method")
	}

	m.mock.funcGetContext = f
	return m.mock
}

// When sets expectation for the HttpProxyTest.GetContext which will trigger the result defined by the following
// Then helper
func (m *mHttpProxyTestMockGetContext) When(req *http.Request) *HttpProxyTestMockGetContextExpectation {
	if m.mock.funcGetContext != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.GetContext mock is already set by Set")
	}

	expectation := &HttpProxyTestMockGetContextExpectation{
		mock:   m.mock,
		params: &HttpProxyTestMockGetContextParams{req},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up HttpProxyTest.GetContext return parameters for the expectation previously defined by the When method
func (e *HttpProxyTestMockGetContextExpectation) Then(c1 context.Context, err error) *HttpProxyTestMock {
	e.results = &HttpProxyTestMockGetContextResults{c1, err}
	return e.mock
}

// GetContext implements HttpProxyTest
func (m *HttpProxyTestMock) GetContext(req *http.Request) (c1 context.Context, err error) {
	atomic.AddUint64(&m.beforeGetContextCounter, 1)
	defer atomic.AddUint64(&m.afterGetContextCounter, 1)

	for _, e := range m.GetContextMock.expectations {
		if minimock.Equal(*e.params, HttpProxyTestMockGetContextParams{req}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.c1, e.results.err
		}
	}

	if m.GetContextMock.defaultExpectation != nil {
		atomic.AddUint64(&m.GetContextMock.defaultExpectation.Counter, 1)
		want := m.GetContextMock.defaultExpectation.params
		got := HttpProxyTestMockGetContextParams{req}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("HttpProxyTestMock.GetContext got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetContextMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the HttpProxyTestMock.GetContext")
		}
		return (*results).c1, (*results).err
	}
	if m.funcGetContext != nil {
		return m.funcGetContext(req)
	}
	m.t.Fatalf("Unexpected call to HttpProxyTestMock.GetContext. %v", req)
	return
}

// GetContextAfterCounter returns a count of finished HttpProxyTestMock.GetContext invocations
func (m *HttpProxyTestMock) GetContextAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterGetContextCounter)
}

// GetContextBeforeCounter returns a count of HttpProxyTestMock.GetContext invocations
func (m *HttpProxyTestMock) GetContextBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeGetContextCounter)
}

// MinimockGetContextDone returns true if the count of the GetContext invocations corresponds
// the number of defined expectations
func (m *HttpProxyTestMock) MinimockGetContextDone() bool {
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
func (m *HttpProxyTestMock) MinimockGetContextInspect() {
	for _, e := range m.GetContextMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to HttpProxyTestMock.GetContext with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetContextMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetContextCounter) < 1 {
		m.t.Errorf("Expected call to HttpProxyTestMock.GetContext with params: %#v", *m.GetContextMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetContext != nil && atomic.LoadUint64(&m.afterGetContextCounter) < 1 {
		m.t.Error("Expected call to HttpProxyTestMock.GetContext")
	}
}

type mHttpProxyTestMockHandleHTTPValidation struct {
	mock               *HttpProxyTestMock
	defaultExpectation *HttpProxyTestMockHandleHTTPValidationExpectation
	expectations       []*HttpProxyTestMockHandleHTTPValidationExpectation
}

// HttpProxyTestMockHandleHTTPValidationExpectation specifies expectation struct of the HttpProxyTest.HandleHTTPValidation
type HttpProxyTestMockHandleHTTPValidationExpectation struct {
	mock    *HttpProxyTestMock
	params  *HttpProxyTestMockHandleHTTPValidationParams
	results *HttpProxyTestMockHandleHTTPValidationResults
	Counter uint64
}

// HttpProxyTestMockHandleHTTPValidationParams contains parameters of the HttpProxyTest.HandleHTTPValidation
type HttpProxyTestMockHandleHTTPValidationParams struct {
	w http.ResponseWriter
	r *http.Request
}

// HttpProxyTestMockHandleHTTPValidationResults contains results of the HttpProxyTest.HandleHTTPValidation
type HttpProxyTestMockHandleHTTPValidationResults struct {
	b1 bool
}

// Expect sets up expected params for HttpProxyTest.HandleHTTPValidation
func (m *mHttpProxyTestMockHandleHTTPValidation) Expect(w http.ResponseWriter, r *http.Request) *mHttpProxyTestMockHandleHTTPValidation {
	if m.mock.funcHandleHTTPValidation != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.HandleHTTPValidation mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockHandleHTTPValidationExpectation{}
	}

	m.defaultExpectation.params = &HttpProxyTestMockHandleHTTPValidationParams{w, r}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by HttpProxyTest.HandleHTTPValidation
func (m *mHttpProxyTestMockHandleHTTPValidation) Return(b1 bool) *HttpProxyTestMock {
	if m.mock.funcHandleHTTPValidation != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.HandleHTTPValidation mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockHandleHTTPValidationExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &HttpProxyTestMockHandleHTTPValidationResults{b1}
	return m.mock
}

//Set uses given function f to mock the HttpProxyTest.HandleHTTPValidation method
func (m *mHttpProxyTestMockHandleHTTPValidation) Set(f func(w http.ResponseWriter, r *http.Request) (b1 bool)) *HttpProxyTestMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the HttpProxyTest.HandleHTTPValidation method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the HttpProxyTest.HandleHTTPValidation method")
	}

	m.mock.funcHandleHTTPValidation = f
	return m.mock
}

// When sets expectation for the HttpProxyTest.HandleHTTPValidation which will trigger the result defined by the following
// Then helper
func (m *mHttpProxyTestMockHandleHTTPValidation) When(w http.ResponseWriter, r *http.Request) *HttpProxyTestMockHandleHTTPValidationExpectation {
	if m.mock.funcHandleHTTPValidation != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.HandleHTTPValidation mock is already set by Set")
	}

	expectation := &HttpProxyTestMockHandleHTTPValidationExpectation{
		mock:   m.mock,
		params: &HttpProxyTestMockHandleHTTPValidationParams{w, r},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up HttpProxyTest.HandleHTTPValidation return parameters for the expectation previously defined by the When method
func (e *HttpProxyTestMockHandleHTTPValidationExpectation) Then(b1 bool) *HttpProxyTestMock {
	e.results = &HttpProxyTestMockHandleHTTPValidationResults{b1}
	return e.mock
}

// HandleHTTPValidation implements HttpProxyTest
func (m *HttpProxyTestMock) HandleHTTPValidation(w http.ResponseWriter, r *http.Request) (b1 bool) {
	atomic.AddUint64(&m.beforeHandleHTTPValidationCounter, 1)
	defer atomic.AddUint64(&m.afterHandleHTTPValidationCounter, 1)

	for _, e := range m.HandleHTTPValidationMock.expectations {
		if minimock.Equal(*e.params, HttpProxyTestMockHandleHTTPValidationParams{w, r}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.b1
		}
	}

	if m.HandleHTTPValidationMock.defaultExpectation != nil {
		atomic.AddUint64(&m.HandleHTTPValidationMock.defaultExpectation.Counter, 1)
		want := m.HandleHTTPValidationMock.defaultExpectation.params
		got := HttpProxyTestMockHandleHTTPValidationParams{w, r}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("HttpProxyTestMock.HandleHTTPValidation got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.HandleHTTPValidationMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the HttpProxyTestMock.HandleHTTPValidation")
		}
		return (*results).b1
	}
	if m.funcHandleHTTPValidation != nil {
		return m.funcHandleHTTPValidation(w, r)
	}
	m.t.Fatalf("Unexpected call to HttpProxyTestMock.HandleHTTPValidation. %v %v", w, r)
	return
}

// HandleHTTPValidationAfterCounter returns a count of finished HttpProxyTestMock.HandleHTTPValidation invocations
func (m *HttpProxyTestMock) HandleHTTPValidationAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterHandleHTTPValidationCounter)
}

// HandleHTTPValidationBeforeCounter returns a count of HttpProxyTestMock.HandleHTTPValidation invocations
func (m *HttpProxyTestMock) HandleHTTPValidationBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeHandleHTTPValidationCounter)
}

// MinimockHandleHTTPValidationDone returns true if the count of the HandleHTTPValidation invocations corresponds
// the number of defined expectations
func (m *HttpProxyTestMock) MinimockHandleHTTPValidationDone() bool {
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
func (m *HttpProxyTestMock) MinimockHandleHTTPValidationInspect() {
	for _, e := range m.HandleHTTPValidationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to HttpProxyTestMock.HandleHTTPValidation with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HandleHTTPValidationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterHandleHTTPValidationCounter) < 1 {
		m.t.Errorf("Expected call to HttpProxyTestMock.HandleHTTPValidation with params: %#v", *m.HandleHTTPValidationMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHandleHTTPValidation != nil && atomic.LoadUint64(&m.afterHandleHTTPValidationCounter) < 1 {
		m.t.Error("Expected call to HttpProxyTestMock.HandleHTTPValidation")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *HttpProxyTestMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetContextInspect()

		m.MinimockHandleHTTPValidationInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *HttpProxyTestMock) MinimockWait(timeout time.Duration) {
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

func (m *HttpProxyTestMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetContextDone() &&
		m.MinimockHandleHTTPValidationDone()
}
