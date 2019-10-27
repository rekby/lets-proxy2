package cert_manager

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cache.Value -o ./value_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"github.com/gojuno/minimock/v3"
)

// ValueMock implements cache.Value
type ValueMock struct {
	t minimock.Tester

	funcDelete          func(ctx context.Context, key string) (err error)
	afterDeleteCounter  uint64
	beforeDeleteCounter uint64
	DeleteMock          mValueMockDelete

	funcGet          func(ctx context.Context, key string) (p1 interface{}, err error)
	afterGetCounter  uint64
	beforeGetCounter uint64
	GetMock          mValueMockGet

	funcPut          func(ctx context.Context, key string, value interface{}) (err error)
	afterPutCounter  uint64
	beforePutCounter uint64
	PutMock          mValueMockPut
}

// NewValueMock returns a mock for cache.Value
func NewValueMock(t minimock.Tester) *ValueMock {
	m := &ValueMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.DeleteMock = mValueMockDelete{mock: m}
	m.GetMock = mValueMockGet{mock: m}
	m.PutMock = mValueMockPut{mock: m}

	return m
}

type mValueMockDelete struct {
	mock               *ValueMock
	defaultExpectation *ValueMockDeleteExpectation
	expectations       []*ValueMockDeleteExpectation
}

// ValueMockDeleteExpectation specifies expectation struct of the Value.Delete
type ValueMockDeleteExpectation struct {
	mock    *ValueMock
	params  *ValueMockDeleteParams
	results *ValueMockDeleteResults
	Counter uint64
}

// ValueMockDeleteParams contains parameters of the Value.Delete
type ValueMockDeleteParams struct {
	ctx context.Context
	key string
}

// ValueMockDeleteResults contains results of the Value.Delete
type ValueMockDeleteResults struct {
	err error
}

// Expect sets up expected params for Value.Delete
func (m *mValueMockDelete) Expect(ctx context.Context, key string) *mValueMockDelete {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("ValueMock.Delete mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ValueMockDeleteExpectation{}
	}

	m.defaultExpectation.params = &ValueMockDeleteParams{ctx, key}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Value.Delete
func (m *mValueMockDelete) Return(err error) *ValueMock {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("ValueMock.Delete mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ValueMockDeleteExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ValueMockDeleteResults{err}
	return m.mock
}

//Set uses given function f to mock the Value.Delete method
func (m *mValueMockDelete) Set(f func(ctx context.Context, key string) (err error)) *ValueMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Value.Delete method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Value.Delete method")
	}

	m.mock.funcDelete = f
	return m.mock
}

// When sets expectation for the Value.Delete which will trigger the result defined by the following
// Then helper
func (m *mValueMockDelete) When(ctx context.Context, key string) *ValueMockDeleteExpectation {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("ValueMock.Delete mock is already set by Set")
	}

	expectation := &ValueMockDeleteExpectation{
		mock:   m.mock,
		params: &ValueMockDeleteParams{ctx, key},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Value.Delete return parameters for the expectation previously defined by the When method
func (e *ValueMockDeleteExpectation) Then(err error) *ValueMock {
	e.results = &ValueMockDeleteResults{err}
	return e.mock
}

// Delete implements cache.Value
func (m *ValueMock) Delete(ctx context.Context, key string) (err error) {
	atomic.AddUint64(&m.beforeDeleteCounter, 1)
	defer atomic.AddUint64(&m.afterDeleteCounter, 1)

	for _, e := range m.DeleteMock.expectations {
		if minimock.Equal(*e.params, ValueMockDeleteParams{ctx, key}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.DeleteMock.defaultExpectation != nil {
		atomic.AddUint64(&m.DeleteMock.defaultExpectation.Counter, 1)
		want := m.DeleteMock.defaultExpectation.params
		got := ValueMockDeleteParams{ctx, key}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ValueMock.Delete got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.DeleteMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ValueMock.Delete")
		}
		return (*results).err
	}
	if m.funcDelete != nil {
		return m.funcDelete(ctx, key)
	}
	m.t.Fatalf("Unexpected call to ValueMock.Delete. %v %v", ctx, key)
	return
}

// DeleteAfterCounter returns a count of finished ValueMock.Delete invocations
func (m *ValueMock) DeleteAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterDeleteCounter)
}

// DeleteBeforeCounter returns a count of ValueMock.Delete invocations
func (m *ValueMock) DeleteBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeDeleteCounter)
}

// MinimockDeleteDone returns true if the count of the Delete invocations corresponds
// the number of defined expectations
func (m *ValueMock) MinimockDeleteDone() bool {
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
func (m *ValueMock) MinimockDeleteInspect() {
	for _, e := range m.DeleteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ValueMock.Delete with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		m.t.Errorf("Expected call to ValueMock.Delete with params: %#v", *m.DeleteMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDelete != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		m.t.Error("Expected call to ValueMock.Delete")
	}
}

type mValueMockGet struct {
	mock               *ValueMock
	defaultExpectation *ValueMockGetExpectation
	expectations       []*ValueMockGetExpectation
}

// ValueMockGetExpectation specifies expectation struct of the Value.Get
type ValueMockGetExpectation struct {
	mock    *ValueMock
	params  *ValueMockGetParams
	results *ValueMockGetResults
	Counter uint64
}

// ValueMockGetParams contains parameters of the Value.Get
type ValueMockGetParams struct {
	ctx context.Context
	key string
}

// ValueMockGetResults contains results of the Value.Get
type ValueMockGetResults struct {
	p1  interface{}
	err error
}

// Expect sets up expected params for Value.Get
func (m *mValueMockGet) Expect(ctx context.Context, key string) *mValueMockGet {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("ValueMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ValueMockGetExpectation{}
	}

	m.defaultExpectation.params = &ValueMockGetParams{ctx, key}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Value.Get
func (m *mValueMockGet) Return(p1 interface{}, err error) *ValueMock {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("ValueMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ValueMockGetExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ValueMockGetResults{p1, err}
	return m.mock
}

//Set uses given function f to mock the Value.Get method
func (m *mValueMockGet) Set(f func(ctx context.Context, key string) (p1 interface{}, err error)) *ValueMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Value.Get method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Value.Get method")
	}

	m.mock.funcGet = f
	return m.mock
}

// When sets expectation for the Value.Get which will trigger the result defined by the following
// Then helper
func (m *mValueMockGet) When(ctx context.Context, key string) *ValueMockGetExpectation {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("ValueMock.Get mock is already set by Set")
	}

	expectation := &ValueMockGetExpectation{
		mock:   m.mock,
		params: &ValueMockGetParams{ctx, key},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Value.Get return parameters for the expectation previously defined by the When method
func (e *ValueMockGetExpectation) Then(p1 interface{}, err error) *ValueMock {
	e.results = &ValueMockGetResults{p1, err}
	return e.mock
}

// Get implements cache.Value
func (m *ValueMock) Get(ctx context.Context, key string) (p1 interface{}, err error) {
	atomic.AddUint64(&m.beforeGetCounter, 1)
	defer atomic.AddUint64(&m.afterGetCounter, 1)

	for _, e := range m.GetMock.expectations {
		if minimock.Equal(*e.params, ValueMockGetParams{ctx, key}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.p1, e.results.err
		}
	}

	if m.GetMock.defaultExpectation != nil {
		atomic.AddUint64(&m.GetMock.defaultExpectation.Counter, 1)
		want := m.GetMock.defaultExpectation.params
		got := ValueMockGetParams{ctx, key}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ValueMock.Get got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ValueMock.Get")
		}
		return (*results).p1, (*results).err
	}
	if m.funcGet != nil {
		return m.funcGet(ctx, key)
	}
	m.t.Fatalf("Unexpected call to ValueMock.Get. %v %v", ctx, key)
	return
}

// GetAfterCounter returns a count of finished ValueMock.Get invocations
func (m *ValueMock) GetAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterGetCounter)
}

// GetBeforeCounter returns a count of ValueMock.Get invocations
func (m *ValueMock) GetBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeGetCounter)
}

// MinimockGetDone returns true if the count of the Get invocations corresponds
// the number of defined expectations
func (m *ValueMock) MinimockGetDone() bool {
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
func (m *ValueMock) MinimockGetInspect() {
	for _, e := range m.GetMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ValueMock.Get with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Errorf("Expected call to ValueMock.Get with params: %#v", *m.GetMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGet != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Error("Expected call to ValueMock.Get")
	}
}

type mValueMockPut struct {
	mock               *ValueMock
	defaultExpectation *ValueMockPutExpectation
	expectations       []*ValueMockPutExpectation
}

// ValueMockPutExpectation specifies expectation struct of the Value.Put
type ValueMockPutExpectation struct {
	mock    *ValueMock
	params  *ValueMockPutParams
	results *ValueMockPutResults
	Counter uint64
}

// ValueMockPutParams contains parameters of the Value.Put
type ValueMockPutParams struct {
	ctx   context.Context
	key   string
	value interface{}
}

// ValueMockPutResults contains results of the Value.Put
type ValueMockPutResults struct {
	err error
}

// Expect sets up expected params for Value.Put
func (m *mValueMockPut) Expect(ctx context.Context, key string, value interface{}) *mValueMockPut {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("ValueMock.Put mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ValueMockPutExpectation{}
	}

	m.defaultExpectation.params = &ValueMockPutParams{ctx, key, value}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Value.Put
func (m *mValueMockPut) Return(err error) *ValueMock {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("ValueMock.Put mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ValueMockPutExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ValueMockPutResults{err}
	return m.mock
}

//Set uses given function f to mock the Value.Put method
func (m *mValueMockPut) Set(f func(ctx context.Context, key string, value interface{}) (err error)) *ValueMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Value.Put method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Value.Put method")
	}

	m.mock.funcPut = f
	return m.mock
}

// When sets expectation for the Value.Put which will trigger the result defined by the following
// Then helper
func (m *mValueMockPut) When(ctx context.Context, key string, value interface{}) *ValueMockPutExpectation {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("ValueMock.Put mock is already set by Set")
	}

	expectation := &ValueMockPutExpectation{
		mock:   m.mock,
		params: &ValueMockPutParams{ctx, key, value},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Value.Put return parameters for the expectation previously defined by the When method
func (e *ValueMockPutExpectation) Then(err error) *ValueMock {
	e.results = &ValueMockPutResults{err}
	return e.mock
}

// Put implements cache.Value
func (m *ValueMock) Put(ctx context.Context, key string, value interface{}) (err error) {
	atomic.AddUint64(&m.beforePutCounter, 1)
	defer atomic.AddUint64(&m.afterPutCounter, 1)

	for _, e := range m.PutMock.expectations {
		if minimock.Equal(*e.params, ValueMockPutParams{ctx, key, value}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.PutMock.defaultExpectation != nil {
		atomic.AddUint64(&m.PutMock.defaultExpectation.Counter, 1)
		want := m.PutMock.defaultExpectation.params
		got := ValueMockPutParams{ctx, key, value}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ValueMock.Put got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.PutMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ValueMock.Put")
		}
		return (*results).err
	}
	if m.funcPut != nil {
		return m.funcPut(ctx, key, value)
	}
	m.t.Fatalf("Unexpected call to ValueMock.Put. %v %v %v", ctx, key, value)
	return
}

// PutAfterCounter returns a count of finished ValueMock.Put invocations
func (m *ValueMock) PutAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterPutCounter)
}

// PutBeforeCounter returns a count of ValueMock.Put invocations
func (m *ValueMock) PutBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforePutCounter)
}

// MinimockPutDone returns true if the count of the Put invocations corresponds
// the number of defined expectations
func (m *ValueMock) MinimockPutDone() bool {
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
func (m *ValueMock) MinimockPutInspect() {
	for _, e := range m.PutMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ValueMock.Put with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PutMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		m.t.Errorf("Expected call to ValueMock.Put with params: %#v", *m.PutMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPut != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		m.t.Error("Expected call to ValueMock.Put")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ValueMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDeleteInspect()

		m.MinimockGetInspect()

		m.MinimockPutInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ValueMock) MinimockWait(timeout time.Duration) {
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

func (m *ValueMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDeleteDone() &&
		m.MinimockGetDone() &&
		m.MinimockPutDone()
}
