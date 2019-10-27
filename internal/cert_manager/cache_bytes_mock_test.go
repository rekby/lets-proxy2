package cert_manager

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cache.Bytes -o ./cache_bytes_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"github.com/gojuno/minimock/v3"
)

// BytesMock implements cache.Bytes
type BytesMock struct {
	t minimock.Tester

	funcDelete          func(ctx context.Context, key string) (err error)
	afterDeleteCounter  uint64
	beforeDeleteCounter uint64
	DeleteMock          mBytesMockDelete

	funcGet          func(ctx context.Context, key string) (ba1 []byte, err error)
	afterGetCounter  uint64
	beforeGetCounter uint64
	GetMock          mBytesMockGet

	funcPut          func(ctx context.Context, key string, data []byte) (err error)
	afterPutCounter  uint64
	beforePutCounter uint64
	PutMock          mBytesMockPut
}

// NewBytesMock returns a mock for cache.Bytes
func NewBytesMock(t minimock.Tester) *BytesMock {
	m := &BytesMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.DeleteMock = mBytesMockDelete{mock: m}
	m.GetMock = mBytesMockGet{mock: m}
	m.PutMock = mBytesMockPut{mock: m}

	return m
}

type mBytesMockDelete struct {
	mock               *BytesMock
	defaultExpectation *BytesMockDeleteExpectation
	expectations       []*BytesMockDeleteExpectation
}

// BytesMockDeleteExpectation specifies expectation struct of the Bytes.Delete
type BytesMockDeleteExpectation struct {
	mock    *BytesMock
	params  *BytesMockDeleteParams
	results *BytesMockDeleteResults
	Counter uint64
}

// BytesMockDeleteParams contains parameters of the Bytes.Delete
type BytesMockDeleteParams struct {
	ctx context.Context
	key string
}

// BytesMockDeleteResults contains results of the Bytes.Delete
type BytesMockDeleteResults struct {
	err error
}

// Expect sets up expected params for Bytes.Delete
func (m *mBytesMockDelete) Expect(ctx context.Context, key string) *mBytesMockDelete {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("BytesMock.Delete mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &BytesMockDeleteExpectation{}
	}

	m.defaultExpectation.params = &BytesMockDeleteParams{ctx, key}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Bytes.Delete
func (m *mBytesMockDelete) Return(err error) *BytesMock {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("BytesMock.Delete mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &BytesMockDeleteExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &BytesMockDeleteResults{err}
	return m.mock
}

//Set uses given function f to mock the Bytes.Delete method
func (m *mBytesMockDelete) Set(f func(ctx context.Context, key string) (err error)) *BytesMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Bytes.Delete method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Bytes.Delete method")
	}

	m.mock.funcDelete = f
	return m.mock
}

// When sets expectation for the Bytes.Delete which will trigger the result defined by the following
// Then helper
func (m *mBytesMockDelete) When(ctx context.Context, key string) *BytesMockDeleteExpectation {
	if m.mock.funcDelete != nil {
		m.mock.t.Fatalf("BytesMock.Delete mock is already set by Set")
	}

	expectation := &BytesMockDeleteExpectation{
		mock:   m.mock,
		params: &BytesMockDeleteParams{ctx, key},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Bytes.Delete return parameters for the expectation previously defined by the When method
func (e *BytesMockDeleteExpectation) Then(err error) *BytesMock {
	e.results = &BytesMockDeleteResults{err}
	return e.mock
}

// Delete implements cache.Bytes
func (m *BytesMock) Delete(ctx context.Context, key string) (err error) {
	atomic.AddUint64(&m.beforeDeleteCounter, 1)
	defer atomic.AddUint64(&m.afterDeleteCounter, 1)

	for _, e := range m.DeleteMock.expectations {
		if minimock.Equal(*e.params, BytesMockDeleteParams{ctx, key}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.DeleteMock.defaultExpectation != nil {
		atomic.AddUint64(&m.DeleteMock.defaultExpectation.Counter, 1)
		want := m.DeleteMock.defaultExpectation.params
		got := BytesMockDeleteParams{ctx, key}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("BytesMock.Delete got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.DeleteMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the BytesMock.Delete")
		}
		return (*results).err
	}
	if m.funcDelete != nil {
		return m.funcDelete(ctx, key)
	}
	m.t.Fatalf("Unexpected call to BytesMock.Delete. %v %v", ctx, key)
	return
}

// DeleteAfterCounter returns a count of finished BytesMock.Delete invocations
func (m *BytesMock) DeleteAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterDeleteCounter)
}

// DeleteBeforeCounter returns a count of BytesMock.Delete invocations
func (m *BytesMock) DeleteBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeDeleteCounter)
}

// MinimockDeleteDone returns true if the count of the Delete invocations corresponds
// the number of defined expectations
func (m *BytesMock) MinimockDeleteDone() bool {
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
func (m *BytesMock) MinimockDeleteInspect() {
	for _, e := range m.DeleteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to BytesMock.Delete with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		m.t.Errorf("Expected call to BytesMock.Delete with params: %#v", *m.DeleteMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDelete != nil && atomic.LoadUint64(&m.afterDeleteCounter) < 1 {
		m.t.Error("Expected call to BytesMock.Delete")
	}
}

type mBytesMockGet struct {
	mock               *BytesMock
	defaultExpectation *BytesMockGetExpectation
	expectations       []*BytesMockGetExpectation
}

// BytesMockGetExpectation specifies expectation struct of the Bytes.Get
type BytesMockGetExpectation struct {
	mock    *BytesMock
	params  *BytesMockGetParams
	results *BytesMockGetResults
	Counter uint64
}

// BytesMockGetParams contains parameters of the Bytes.Get
type BytesMockGetParams struct {
	ctx context.Context
	key string
}

// BytesMockGetResults contains results of the Bytes.Get
type BytesMockGetResults struct {
	ba1 []byte
	err error
}

// Expect sets up expected params for Bytes.Get
func (m *mBytesMockGet) Expect(ctx context.Context, key string) *mBytesMockGet {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("BytesMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &BytesMockGetExpectation{}
	}

	m.defaultExpectation.params = &BytesMockGetParams{ctx, key}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Bytes.Get
func (m *mBytesMockGet) Return(ba1 []byte, err error) *BytesMock {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("BytesMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &BytesMockGetExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &BytesMockGetResults{ba1, err}
	return m.mock
}

//Set uses given function f to mock the Bytes.Get method
func (m *mBytesMockGet) Set(f func(ctx context.Context, key string) (ba1 []byte, err error)) *BytesMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Bytes.Get method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Bytes.Get method")
	}

	m.mock.funcGet = f
	return m.mock
}

// When sets expectation for the Bytes.Get which will trigger the result defined by the following
// Then helper
func (m *mBytesMockGet) When(ctx context.Context, key string) *BytesMockGetExpectation {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("BytesMock.Get mock is already set by Set")
	}

	expectation := &BytesMockGetExpectation{
		mock:   m.mock,
		params: &BytesMockGetParams{ctx, key},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Bytes.Get return parameters for the expectation previously defined by the When method
func (e *BytesMockGetExpectation) Then(ba1 []byte, err error) *BytesMock {
	e.results = &BytesMockGetResults{ba1, err}
	return e.mock
}

// Get implements cache.Bytes
func (m *BytesMock) Get(ctx context.Context, key string) (ba1 []byte, err error) {
	atomic.AddUint64(&m.beforeGetCounter, 1)
	defer atomic.AddUint64(&m.afterGetCounter, 1)

	for _, e := range m.GetMock.expectations {
		if minimock.Equal(*e.params, BytesMockGetParams{ctx, key}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.ba1, e.results.err
		}
	}

	if m.GetMock.defaultExpectation != nil {
		atomic.AddUint64(&m.GetMock.defaultExpectation.Counter, 1)
		want := m.GetMock.defaultExpectation.params
		got := BytesMockGetParams{ctx, key}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("BytesMock.Get got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the BytesMock.Get")
		}
		return (*results).ba1, (*results).err
	}
	if m.funcGet != nil {
		return m.funcGet(ctx, key)
	}
	m.t.Fatalf("Unexpected call to BytesMock.Get. %v %v", ctx, key)
	return
}

// GetAfterCounter returns a count of finished BytesMock.Get invocations
func (m *BytesMock) GetAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterGetCounter)
}

// GetBeforeCounter returns a count of BytesMock.Get invocations
func (m *BytesMock) GetBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeGetCounter)
}

// MinimockGetDone returns true if the count of the Get invocations corresponds
// the number of defined expectations
func (m *BytesMock) MinimockGetDone() bool {
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
func (m *BytesMock) MinimockGetInspect() {
	for _, e := range m.GetMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to BytesMock.Get with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Errorf("Expected call to BytesMock.Get with params: %#v", *m.GetMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGet != nil && atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Error("Expected call to BytesMock.Get")
	}
}

type mBytesMockPut struct {
	mock               *BytesMock
	defaultExpectation *BytesMockPutExpectation
	expectations       []*BytesMockPutExpectation
}

// BytesMockPutExpectation specifies expectation struct of the Bytes.Put
type BytesMockPutExpectation struct {
	mock    *BytesMock
	params  *BytesMockPutParams
	results *BytesMockPutResults
	Counter uint64
}

// BytesMockPutParams contains parameters of the Bytes.Put
type BytesMockPutParams struct {
	ctx  context.Context
	key  string
	data []byte
}

// BytesMockPutResults contains results of the Bytes.Put
type BytesMockPutResults struct {
	err error
}

// Expect sets up expected params for Bytes.Put
func (m *mBytesMockPut) Expect(ctx context.Context, key string, data []byte) *mBytesMockPut {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("BytesMock.Put mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &BytesMockPutExpectation{}
	}

	m.defaultExpectation.params = &BytesMockPutParams{ctx, key, data}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Bytes.Put
func (m *mBytesMockPut) Return(err error) *BytesMock {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("BytesMock.Put mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &BytesMockPutExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &BytesMockPutResults{err}
	return m.mock
}

//Set uses given function f to mock the Bytes.Put method
func (m *mBytesMockPut) Set(f func(ctx context.Context, key string, data []byte) (err error)) *BytesMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Bytes.Put method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Bytes.Put method")
	}

	m.mock.funcPut = f
	return m.mock
}

// When sets expectation for the Bytes.Put which will trigger the result defined by the following
// Then helper
func (m *mBytesMockPut) When(ctx context.Context, key string, data []byte) *BytesMockPutExpectation {
	if m.mock.funcPut != nil {
		m.mock.t.Fatalf("BytesMock.Put mock is already set by Set")
	}

	expectation := &BytesMockPutExpectation{
		mock:   m.mock,
		params: &BytesMockPutParams{ctx, key, data},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Bytes.Put return parameters for the expectation previously defined by the When method
func (e *BytesMockPutExpectation) Then(err error) *BytesMock {
	e.results = &BytesMockPutResults{err}
	return e.mock
}

// Put implements cache.Bytes
func (m *BytesMock) Put(ctx context.Context, key string, data []byte) (err error) {
	atomic.AddUint64(&m.beforePutCounter, 1)
	defer atomic.AddUint64(&m.afterPutCounter, 1)

	for _, e := range m.PutMock.expectations {
		if minimock.Equal(*e.params, BytesMockPutParams{ctx, key, data}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.PutMock.defaultExpectation != nil {
		atomic.AddUint64(&m.PutMock.defaultExpectation.Counter, 1)
		want := m.PutMock.defaultExpectation.params
		got := BytesMockPutParams{ctx, key, data}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("BytesMock.Put got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.PutMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the BytesMock.Put")
		}
		return (*results).err
	}
	if m.funcPut != nil {
		return m.funcPut(ctx, key, data)
	}
	m.t.Fatalf("Unexpected call to BytesMock.Put. %v %v %v", ctx, key, data)
	return
}

// PutAfterCounter returns a count of finished BytesMock.Put invocations
func (m *BytesMock) PutAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterPutCounter)
}

// PutBeforeCounter returns a count of BytesMock.Put invocations
func (m *BytesMock) PutBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforePutCounter)
}

// MinimockPutDone returns true if the count of the Put invocations corresponds
// the number of defined expectations
func (m *BytesMock) MinimockPutDone() bool {
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
func (m *BytesMock) MinimockPutInspect() {
	for _, e := range m.PutMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to BytesMock.Put with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PutMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		m.t.Errorf("Expected call to BytesMock.Put with params: %#v", *m.PutMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPut != nil && atomic.LoadUint64(&m.afterPutCounter) < 1 {
		m.t.Error("Expected call to BytesMock.Put")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *BytesMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDeleteInspect()

		m.MinimockGetInspect()

		m.MinimockPutInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *BytesMock) MinimockWait(timeout time.Duration) {
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

func (m *BytesMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDeleteDone() &&
		m.MinimockGetDone() &&
		m.MinimockPutDone()
}
