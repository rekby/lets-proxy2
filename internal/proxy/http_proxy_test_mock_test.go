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

	funcGetContext          func(req *http.Request) (c1 context.Context)
	afterGetContextCounter  uint64
	beforeGetContextCounter uint64
	GetContextMock          mHttpProxyTestMockGetContext

	funcGetDestination          func(ctx context.Context, remoteAddr string) (addr string, err error)
	afterGetDestinationCounter  uint64
	beforeGetDestinationCounter uint64
	GetDestinationMock          mHttpProxyTestMockGetDestination

	funcHandleHttpValidation          func(w http.ResponseWriter, r *http.Request) (b1 bool)
	afterHandleHttpValidationCounter  uint64
	beforeHandleHttpValidationCounter uint64
	HandleHttpValidationMock          mHttpProxyTestMockHandleHttpValidation
}

// NewHttpProxyTestMock returns a mock for HttpProxyTest
func NewHttpProxyTestMock(t minimock.Tester) *HttpProxyTestMock {
	m := &HttpProxyTestMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.GetContextMock = mHttpProxyTestMockGetContext{mock: m}
	m.GetDestinationMock = mHttpProxyTestMockGetDestination{mock: m}
	m.HandleHttpValidationMock = mHttpProxyTestMockHandleHttpValidation{mock: m}

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
	c1 context.Context
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
func (m *mHttpProxyTestMockGetContext) Return(c1 context.Context) *HttpProxyTestMock {
	if m.mock.funcGetContext != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.GetContext mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockGetContextExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &HttpProxyTestMockGetContextResults{c1}
	return m.mock
}

//Set uses given function f to mock the HttpProxyTest.GetContext method
func (m *mHttpProxyTestMockGetContext) Set(f func(req *http.Request) (c1 context.Context)) *HttpProxyTestMock {
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
func (e *HttpProxyTestMockGetContextExpectation) Then(c1 context.Context) *HttpProxyTestMock {
	e.results = &HttpProxyTestMockGetContextResults{c1}
	return e.mock
}

// GetContext implements HttpProxyTest
func (m *HttpProxyTestMock) GetContext(req *http.Request) (c1 context.Context) {
	atomic.AddUint64(&m.beforeGetContextCounter, 1)
	defer atomic.AddUint64(&m.afterGetContextCounter, 1)

	for _, e := range m.GetContextMock.expectations {
		if minimock.Equal(*e.params, HttpProxyTestMockGetContextParams{req}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.c1
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
		return (*results).c1
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

type mHttpProxyTestMockGetDestination struct {
	mock               *HttpProxyTestMock
	defaultExpectation *HttpProxyTestMockGetDestinationExpectation
	expectations       []*HttpProxyTestMockGetDestinationExpectation
}

// HttpProxyTestMockGetDestinationExpectation specifies expectation struct of the HttpProxyTest.GetDestination
type HttpProxyTestMockGetDestinationExpectation struct {
	mock    *HttpProxyTestMock
	params  *HttpProxyTestMockGetDestinationParams
	results *HttpProxyTestMockGetDestinationResults
	Counter uint64
}

// HttpProxyTestMockGetDestinationParams contains parameters of the HttpProxyTest.GetDestination
type HttpProxyTestMockGetDestinationParams struct {
	ctx        context.Context
	remoteAddr string
}

// HttpProxyTestMockGetDestinationResults contains results of the HttpProxyTest.GetDestination
type HttpProxyTestMockGetDestinationResults struct {
	addr string
	err  error
}

// Expect sets up expected params for HttpProxyTest.GetDestination
func (m *mHttpProxyTestMockGetDestination) Expect(ctx context.Context, remoteAddr string) *mHttpProxyTestMockGetDestination {
	if m.mock.funcGetDestination != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.GetDestination mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockGetDestinationExpectation{}
	}

	m.defaultExpectation.params = &HttpProxyTestMockGetDestinationParams{ctx, remoteAddr}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by HttpProxyTest.GetDestination
func (m *mHttpProxyTestMockGetDestination) Return(addr string, err error) *HttpProxyTestMock {
	if m.mock.funcGetDestination != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.GetDestination mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockGetDestinationExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &HttpProxyTestMockGetDestinationResults{addr, err}
	return m.mock
}

//Set uses given function f to mock the HttpProxyTest.GetDestination method
func (m *mHttpProxyTestMockGetDestination) Set(f func(ctx context.Context, remoteAddr string) (addr string, err error)) *HttpProxyTestMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the HttpProxyTest.GetDestination method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the HttpProxyTest.GetDestination method")
	}

	m.mock.funcGetDestination = f
	return m.mock
}

// When sets expectation for the HttpProxyTest.GetDestination which will trigger the result defined by the following
// Then helper
func (m *mHttpProxyTestMockGetDestination) When(ctx context.Context, remoteAddr string) *HttpProxyTestMockGetDestinationExpectation {
	if m.mock.funcGetDestination != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.GetDestination mock is already set by Set")
	}

	expectation := &HttpProxyTestMockGetDestinationExpectation{
		mock:   m.mock,
		params: &HttpProxyTestMockGetDestinationParams{ctx, remoteAddr},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up HttpProxyTest.GetDestination return parameters for the expectation previously defined by the When method
func (e *HttpProxyTestMockGetDestinationExpectation) Then(addr string, err error) *HttpProxyTestMock {
	e.results = &HttpProxyTestMockGetDestinationResults{addr, err}
	return e.mock
}

// GetDestination implements HttpProxyTest
func (m *HttpProxyTestMock) GetDestination(ctx context.Context, remoteAddr string) (addr string, err error) {
	atomic.AddUint64(&m.beforeGetDestinationCounter, 1)
	defer atomic.AddUint64(&m.afterGetDestinationCounter, 1)

	for _, e := range m.GetDestinationMock.expectations {
		if minimock.Equal(*e.params, HttpProxyTestMockGetDestinationParams{ctx, remoteAddr}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.addr, e.results.err
		}
	}

	if m.GetDestinationMock.defaultExpectation != nil {
		atomic.AddUint64(&m.GetDestinationMock.defaultExpectation.Counter, 1)
		want := m.GetDestinationMock.defaultExpectation.params
		got := HttpProxyTestMockGetDestinationParams{ctx, remoteAddr}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("HttpProxyTestMock.GetDestination got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetDestinationMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the HttpProxyTestMock.GetDestination")
		}
		return (*results).addr, (*results).err
	}
	if m.funcGetDestination != nil {
		return m.funcGetDestination(ctx, remoteAddr)
	}
	m.t.Fatalf("Unexpected call to HttpProxyTestMock.GetDestination. %v %v", ctx, remoteAddr)
	return
}

// GetDestinationAfterCounter returns a count of finished HttpProxyTestMock.GetDestination invocations
func (m *HttpProxyTestMock) GetDestinationAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterGetDestinationCounter)
}

// GetDestinationBeforeCounter returns a count of HttpProxyTestMock.GetDestination invocations
func (m *HttpProxyTestMock) GetDestinationBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeGetDestinationCounter)
}

// MinimockGetDestinationDone returns true if the count of the GetDestination invocations corresponds
// the number of defined expectations
func (m *HttpProxyTestMock) MinimockGetDestinationDone() bool {
	for _, e := range m.GetDestinationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetDestinationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetDestinationCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetDestination != nil && atomic.LoadUint64(&m.afterGetDestinationCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetDestinationInspect logs each unmet expectation
func (m *HttpProxyTestMock) MinimockGetDestinationInspect() {
	for _, e := range m.GetDestinationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to HttpProxyTestMock.GetDestination with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetDestinationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetDestinationCounter) < 1 {
		m.t.Errorf("Expected call to HttpProxyTestMock.GetDestination with params: %#v", *m.GetDestinationMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetDestination != nil && atomic.LoadUint64(&m.afterGetDestinationCounter) < 1 {
		m.t.Error("Expected call to HttpProxyTestMock.GetDestination")
	}
}

type mHttpProxyTestMockHandleHttpValidation struct {
	mock               *HttpProxyTestMock
	defaultExpectation *HttpProxyTestMockHandleHttpValidationExpectation
	expectations       []*HttpProxyTestMockHandleHttpValidationExpectation
}

// HttpProxyTestMockHandleHttpValidationExpectation specifies expectation struct of the HttpProxyTest.HandleHttpValidation
type HttpProxyTestMockHandleHttpValidationExpectation struct {
	mock    *HttpProxyTestMock
	params  *HttpProxyTestMockHandleHttpValidationParams
	results *HttpProxyTestMockHandleHttpValidationResults
	Counter uint64
}

// HttpProxyTestMockHandleHttpValidationParams contains parameters of the HttpProxyTest.HandleHttpValidation
type HttpProxyTestMockHandleHttpValidationParams struct {
	w http.ResponseWriter
	r *http.Request
}

// HttpProxyTestMockHandleHttpValidationResults contains results of the HttpProxyTest.HandleHttpValidation
type HttpProxyTestMockHandleHttpValidationResults struct {
	b1 bool
}

// Expect sets up expected params for HttpProxyTest.HandleHttpValidation
func (m *mHttpProxyTestMockHandleHttpValidation) Expect(w http.ResponseWriter, r *http.Request) *mHttpProxyTestMockHandleHttpValidation {
	if m.mock.funcHandleHttpValidation != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.HandleHttpValidation mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockHandleHttpValidationExpectation{}
	}

	m.defaultExpectation.params = &HttpProxyTestMockHandleHttpValidationParams{w, r}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by HttpProxyTest.HandleHttpValidation
func (m *mHttpProxyTestMockHandleHttpValidation) Return(b1 bool) *HttpProxyTestMock {
	if m.mock.funcHandleHttpValidation != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.HandleHttpValidation mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &HttpProxyTestMockHandleHttpValidationExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &HttpProxyTestMockHandleHttpValidationResults{b1}
	return m.mock
}

//Set uses given function f to mock the HttpProxyTest.HandleHttpValidation method
func (m *mHttpProxyTestMockHandleHttpValidation) Set(f func(w http.ResponseWriter, r *http.Request) (b1 bool)) *HttpProxyTestMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the HttpProxyTest.HandleHttpValidation method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the HttpProxyTest.HandleHttpValidation method")
	}

	m.mock.funcHandleHttpValidation = f
	return m.mock
}

// When sets expectation for the HttpProxyTest.HandleHttpValidation which will trigger the result defined by the following
// Then helper
func (m *mHttpProxyTestMockHandleHttpValidation) When(w http.ResponseWriter, r *http.Request) *HttpProxyTestMockHandleHttpValidationExpectation {
	if m.mock.funcHandleHttpValidation != nil {
		m.mock.t.Fatalf("HttpProxyTestMock.HandleHttpValidation mock is already set by Set")
	}

	expectation := &HttpProxyTestMockHandleHttpValidationExpectation{
		mock:   m.mock,
		params: &HttpProxyTestMockHandleHttpValidationParams{w, r},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up HttpProxyTest.HandleHttpValidation return parameters for the expectation previously defined by the When method
func (e *HttpProxyTestMockHandleHttpValidationExpectation) Then(b1 bool) *HttpProxyTestMock {
	e.results = &HttpProxyTestMockHandleHttpValidationResults{b1}
	return e.mock
}

// HandleHttpValidation implements HttpProxyTest
func (m *HttpProxyTestMock) HandleHttpValidation(w http.ResponseWriter, r *http.Request) (b1 bool) {
	atomic.AddUint64(&m.beforeHandleHttpValidationCounter, 1)
	defer atomic.AddUint64(&m.afterHandleHttpValidationCounter, 1)

	for _, e := range m.HandleHttpValidationMock.expectations {
		if minimock.Equal(*e.params, HttpProxyTestMockHandleHttpValidationParams{w, r}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.b1
		}
	}

	if m.HandleHttpValidationMock.defaultExpectation != nil {
		atomic.AddUint64(&m.HandleHttpValidationMock.defaultExpectation.Counter, 1)
		want := m.HandleHttpValidationMock.defaultExpectation.params
		got := HttpProxyTestMockHandleHttpValidationParams{w, r}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("HttpProxyTestMock.HandleHttpValidation got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.HandleHttpValidationMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the HttpProxyTestMock.HandleHttpValidation")
		}
		return (*results).b1
	}
	if m.funcHandleHttpValidation != nil {
		return m.funcHandleHttpValidation(w, r)
	}
	m.t.Fatalf("Unexpected call to HttpProxyTestMock.HandleHttpValidation. %v %v", w, r)
	return
}

// HandleHttpValidationAfterCounter returns a count of finished HttpProxyTestMock.HandleHttpValidation invocations
func (m *HttpProxyTestMock) HandleHttpValidationAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterHandleHttpValidationCounter)
}

// HandleHttpValidationBeforeCounter returns a count of HttpProxyTestMock.HandleHttpValidation invocations
func (m *HttpProxyTestMock) HandleHttpValidationBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeHandleHttpValidationCounter)
}

// MinimockHandleHttpValidationDone returns true if the count of the HandleHttpValidation invocations corresponds
// the number of defined expectations
func (m *HttpProxyTestMock) MinimockHandleHttpValidationDone() bool {
	for _, e := range m.HandleHttpValidationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HandleHttpValidationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterHandleHttpValidationCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHandleHttpValidation != nil && atomic.LoadUint64(&m.afterHandleHttpValidationCounter) < 1 {
		return false
	}
	return true
}

// MinimockHandleHttpValidationInspect logs each unmet expectation
func (m *HttpProxyTestMock) MinimockHandleHttpValidationInspect() {
	for _, e := range m.HandleHttpValidationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to HttpProxyTestMock.HandleHttpValidation with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HandleHttpValidationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterHandleHttpValidationCounter) < 1 {
		m.t.Errorf("Expected call to HttpProxyTestMock.HandleHttpValidation with params: %#v", *m.HandleHttpValidationMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHandleHttpValidation != nil && atomic.LoadUint64(&m.afterHandleHttpValidationCounter) < 1 {
		m.t.Error("Expected call to HttpProxyTestMock.HandleHttpValidation")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *HttpProxyTestMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetContextInspect()

		m.MinimockGetDestinationInspect()

		m.MinimockHandleHttpValidationInspect()
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
		m.MinimockGetDestinationDone() &&
		m.MinimockHandleHttpValidationDone()
}
