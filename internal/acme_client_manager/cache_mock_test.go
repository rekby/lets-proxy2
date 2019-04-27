package acme_client_manager

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cache.Cache -o ./cache_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"github.com/gojuno/minimock"
)

// CacheMock implements cache.Cache
type CacheMock struct {
	t minimock.Tester

	funcDelete          func(ctx context.Context, key string) (err error)
	afterDeleteCounter  uint64
	beforeDeleteCounter uint64
	DeleteMock          mCacheMockDelete

	funcGet          func(ctx context.Context, key string) (ba1 []byte, err error)
	afterGetCounter  uint64
	beforeGetCounter uint64
	GetMock          mCacheMockGet

	funcPut          func(ctx context.Context, key string, data []byte) (err error)
	afterPutCounter  uint64
	beforePutCounter uint64
	PutMock          mCacheMockPut
}

// NewCacheMock returns a mock for cache.Cache
func NewCacheMock(t minimock.Tester) *CacheMock {
	m := &CacheMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.DeleteMock = mCacheMockDelete{mock: m}
	m.GetMock = mCacheMockGet{mock: m}
	m.PutMock = mCacheMockPut{mock: m}

	return m
}

type mCacheMockDelete struct {
	mock               *CacheMock
	defaultExpectation *CacheMockDeleteExpectation
	expectations       []*CacheMockDeleteExpectation
}

// CacheMockDeleteExpectation specifies expectation struct of the Cache.Delete
type CacheMockDeleteExpectation struct {
	mock    *CacheMock
	params  *CacheMockDeleteParams
	results *CacheMockDeleteResults
	Counter uint64
}

// CacheMockDeleteParams contains parameters of the Cache.Delete
type CacheMockDeleteParams struct {
	ctx context.Context
	key string
}

// CacheMockDeleteResults contains results of the Cache.Delete
type CacheMockDeleteResults struct {
	err error
}

// Expect sets up expected params for Cache.Delete
func (m *mCacheMockDelete) Expect(ctx context.Context, key string) *mCacheMockDelete {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("CacheMock.Delete mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CacheMockDeleteExpectation{}
	}

	m.defaultExpectation.params = &CacheMockDeleteParams{ctx, key}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Cache.Delete
func (m *mCacheMockDelete) Return(err error) *CacheMock {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("CacheMock.Delete mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CacheMockDeleteExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &CacheMockDeleteResults{err}
	return m.mock
}

//Set uses given function f to mock the Cache.Delete method
func (m *mCacheMockDelete) Set(f func(ctx context.Context, key string) (err error)) *CacheMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Cache.Delete method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Cache.Delete method")
	}

	m.mock.funcDelete = f
	return m.mock
}

// When sets expectation for the Cache.Delete which will trigger the result defined by the following
// Then helper
func (m *mCacheMockDelete) When(ctx context.Context, key string) *CacheMockDeleteExpectation {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("CacheMock.Delete mock is already set by Set")
	}

	expectation := &CacheMockDeleteExpectation{
		mock:   m.mock,
		params: &CacheMockDeleteParams{ctx, key},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Cache.Delete return parameters for the expectation previously defined by the When method
func (e *CacheMockDeleteExpectation) Then(err error) *CacheMock {
	e.results = &CacheMockDeleteResults{err}
	return e.mock
}

// Delete implements cache.Cache
func (m *CacheMock) Delete(ctx context.Context, key string) (err error) {
	atomic.AddUint64(&m.beforeDeleteCounter, 1)
	defer atomic.AddUint64(&m.afterDeleteCounter, 1)

	for _, e := range m.DeleteMock.expectations {
		if minimock.Equal(*e.params, CacheMockDeleteParams{ctx, key}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.DeleteMock.defaultExpectation != nil {
		atomic.AddUint64(&m.DeleteMock.defaultExpectation.Counter, 1)
		want := m.DeleteMock.defaultExpectation.params
		got := CacheMockDeleteParams{ctx, key}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("CacheMock.Delete got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.DeleteMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the CacheMock.Delete")
		}
		return (*results).err
	}
	if m.funcDelete != nil {
		return m.funcDelete(ctx, key)
	}
	m.t.Fatalf("Unexpected call to CacheMock.Delete. %v %v", ctx, key)
	return
}

// DeleteAfterCounter returns a count of finished CacheMock.Delete invocations
func (m *CacheMock) DeleteAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterDeleteCounter)
}

// DeleteBeforeCounter returns a count of CacheMock.Delete invocations
func (m *CacheMock) DeleteBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeDeleteCounter)
}

// MinimockDeleteDone returns true if the count of the Delete invocations corresponds
// the number of defined expectations
func (m *CacheMock) MinimockDeleteDone() bool {
	for _, e := range m.DeleteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDelete != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		return false
	}
	return true
}

// MinimockDeleteInspect logs each unmet expectation
func (m *CacheMock) MinimockDeleteInspect() {
	for _, e := range m.DeleteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CacheMock.Delete with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		m.t.Errorf("Expected call to CacheMock.Delete with params: %#v", *m.DeleteMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDelete != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		m.t.Error("Expected call to CacheMock.Delete")
	}
}

type mCacheMockGet struct {
	mock               *CacheMock
	defaultExpectation *CacheMockGetExpectation
	expectations       []*CacheMockGetExpectation
}

// CacheMockGetExpectation specifies expectation struct of the Cache.Get
type CacheMockGetExpectation struct {
	mock    *CacheMock
	params  *CacheMockGetParams
	results *CacheMockGetResults
	Counter uint64
}

// CacheMockGetParams contains parameters of the Cache.Get
type CacheMockGetParams struct {
	ctx context.Context
	key string
}

// CacheMockGetResults contains results of the Cache.Get
type CacheMockGetResults struct {
	ba1 []byte
	err error
}

// Expect sets up expected params for Cache.Get
func (m *mCacheMockGet) Expect(ctx context.Context, key string) *mCacheMockGet {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("CacheMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CacheMockGetExpectation{}
	}

	m.defaultExpectation.params = &CacheMockGetParams{ctx, key}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Cache.Get
func (m *mCacheMockGet) Return(ba1 []byte, err error) *CacheMock {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("CacheMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CacheMockGetExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &CacheMockGetResults{ba1, err}
	return m.mock
}

//Set uses given function f to mock the Cache.Get method
func (m *mCacheMockGet) Set(f func(ctx context.Context, key string) (ba1 []byte, err error)) *CacheMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Cache.Get method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Cache.Get method")
	}

	m.mock.funcGet = f
	return m.mock
}

// When sets expectation for the Cache.Get which will trigger the result defined by the following
// Then helper
func (m *mCacheMockGet) When(ctx context.Context, key string) *CacheMockGetExpectation {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("CacheMock.Get mock is already set by Set")
	}

	expectation := &CacheMockGetExpectation{
		mock:   m.mock,
		params: &CacheMockGetParams{ctx, key},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Cache.Get return parameters for the expectation previously defined by the When method
func (e *CacheMockGetExpectation) Then(ba1 []byte, err error) *CacheMock {
	e.results = &CacheMockGetResults{ba1, err}
	return e.mock
}

// Get implements cache.Cache
func (m *CacheMock) Get(ctx context.Context, key string) (ba1 []byte, err error) {
	atomic.AddUint64(&m.beforeGetCounter, 1)
	defer atomic.AddUint64(&m.afterGetCounter, 1)

	for _, e := range m.GetMock.expectations {
		if minimock.Equal(*e.params, CacheMockGetParams{ctx, key}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.ba1, e.results.err
		}
	}

	if m.GetMock.defaultExpectation != nil {
		atomic.AddUint64(&m.GetMock.defaultExpectation.Counter, 1)
		want := m.GetMock.defaultExpectation.params
		got := CacheMockGetParams{ctx, key}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("CacheMock.Get got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the CacheMock.Get")
		}
		return (*results).ba1, (*results).err
	}
	if m.funcGet != nil {
		return m.funcGet(ctx, key)
	}
	m.t.Fatalf("Unexpected call to CacheMock.Get. %v %v", ctx, key)
	return
}

// GetAfterCounter returns a count of finished CacheMock.Get invocations
func (m *CacheMock) GetAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterGetCounter)
}

// GetBeforeCounter returns a count of CacheMock.Get invocations
func (m *CacheMock) GetBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeGetCounter)
}

// MinimockGetDone returns true if the count of the Get invocations corresponds
// the number of defined expectations
func (m *CacheMock) MinimockGetDone() bool {
	for _, e := range m.GetMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGet != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetInspect logs each unmet expectation
func (m *CacheMock) MinimockGetInspect() {
	for _, e := range m.GetMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CacheMock.Get with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Errorf("Expected call to CacheMock.Get with params: %#v", *m.GetMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGet != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Error("Expected call to CacheMock.Get")
	}
}

type mCacheMockPut struct {
	mock               *CacheMock
	defaultExpectation *CacheMockPutExpectation
	expectations       []*CacheMockPutExpectation
}

// CacheMockPutExpectation specifies expectation struct of the Cache.Put
type CacheMockPutExpectation struct {
	mock    *CacheMock
	params  *CacheMockPutParams
	results *CacheMockPutResults
	Counter uint64
}

// CacheMockPutParams contains parameters of the Cache.Put
type CacheMockPutParams struct {
	ctx  context.Context
	key  string
	data []byte
}

// CacheMockPutResults contains results of the Cache.Put
type CacheMockPutResults struct {
	err error
}

// Expect sets up expected params for Cache.Put
func (m *mCacheMockPut) Expect(ctx context.Context, key string, data []byte) *mCacheMockPut {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("CacheMock.Put mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CacheMockPutExpectation{}
	}

	m.defaultExpectation.params = &CacheMockPutParams{ctx, key, data}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Cache.Put
func (m *mCacheMockPut) Return(err error) *CacheMock {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("CacheMock.Put mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CacheMockPutExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &CacheMockPutResults{err}
	return m.mock
}

//Set uses given function f to mock the Cache.Put method
func (m *mCacheMockPut) Set(f func(ctx context.Context, key string, data []byte) (err error)) *CacheMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Cache.Put method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Cache.Put method")
	}

	m.mock.funcPut = f
	return m.mock
}

// When sets expectation for the Cache.Put which will trigger the result defined by the following
// Then helper
func (m *mCacheMockPut) When(ctx context.Context, key string, data []byte) *CacheMockPutExpectation {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("CacheMock.Put mock is already set by Set")
	}

	expectation := &CacheMockPutExpectation{
		mock:   m.mock,
		params: &CacheMockPutParams{ctx, key, data},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Cache.Put return parameters for the expectation previously defined by the When method
func (e *CacheMockPutExpectation) Then(err error) *CacheMock {
	e.results = &CacheMockPutResults{err}
	return e.mock
}

// Put implements cache.Cache
func (m *CacheMock) Put(ctx context.Context, key string, data []byte) (err error) {
	atomic.AddUint64(&m.beforePutCounter, 1)
	defer atomic.AddUint64(&m.afterPutCounter, 1)

	for _, e := range m.PutMock.expectations {
		if minimock.Equal(*e.params, CacheMockPutParams{ctx, key, data}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.PutMock.defaultExpectation != nil {
		atomic.AddUint64(&m.PutMock.defaultExpectation.Counter, 1)
		want := m.PutMock.defaultExpectation.params
		got := CacheMockPutParams{ctx, key, data}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("CacheMock.Put got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.PutMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the CacheMock.Put")
		}
		return (*results).err
	}
	if m.funcPut != nil {
		return m.funcPut(ctx, key, data)
	}
	m.t.Fatalf("Unexpected call to CacheMock.Put. %v %v %v", ctx, key, data)
	return
}

// PutAfterCounter returns a count of finished CacheMock.Put invocations
func (m *CacheMock) PutAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterPutCounter)
}

// PutBeforeCounter returns a count of CacheMock.Put invocations
func (m *CacheMock) PutBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforePutCounter)
}

// MinimockPutDone returns true if the count of the Put invocations corresponds
// the number of defined expectations
func (m *CacheMock) MinimockPutDone() bool {
	for _, e := range m.PutMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PutMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPut != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		return false
	}
	return true
}

// MinimockPutInspect logs each unmet expectation
func (m *CacheMock) MinimockPutInspect() {
	for _, e := range m.PutMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CacheMock.Put with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PutMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		m.t.Errorf("Expected call to CacheMock.Put with params: %#v", *m.PutMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPut != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		m.t.Error("Expected call to CacheMock.Put")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CacheMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDeleteInspect()

		m.MinimockGetInspect()

		m.MinimockPutInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CacheMock) MinimockWait(timeout time.Duration) {
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

func (m *CacheMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDeleteDone() &&
		m.MinimockGetDone() &&
		m.MinimockPutDone()
}
