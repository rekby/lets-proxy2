package proxy

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i net.Listener -o ./net_listener_mock_test.go

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

// ListenerMock implements net.Listener
type ListenerMock struct {
	t minimock.Tester

	funcAccept          func() (c1 net.Conn, err error)
	afterAcceptCounter  uint64
	beforeAcceptCounter uint64
	AcceptMock          mListenerMockAccept

	funcAddr          func() (a1 net.Addr)
	afterAddrCounter  uint64
	beforeAddrCounter uint64
	AddrMock          mListenerMockAddr

	funcClose          func() (err error)
	afterCloseCounter  uint64
	beforeCloseCounter uint64
	CloseMock          mListenerMockClose
}

// NewListenerMock returns a mock for net.Listener
func NewListenerMock(t minimock.Tester) *ListenerMock {
	m := &ListenerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.AcceptMock = mListenerMockAccept{mock: m}
	m.AddrMock = mListenerMockAddr{mock: m}
	m.CloseMock = mListenerMockClose{mock: m}

	return m
}

type mListenerMockAccept struct {
	mock               *ListenerMock
	defaultExpectation *ListenerMockAcceptExpectation
	expectations       []*ListenerMockAcceptExpectation
}

// ListenerMockAcceptExpectation specifies expectation struct of the Listener.Accept
type ListenerMockAcceptExpectation struct {
	mock *ListenerMock

	results *ListenerMockAcceptResults
	Counter uint64
}

// ListenerMockAcceptResults contains results of the Listener.Accept
type ListenerMockAcceptResults struct {
	c1  net.Conn
	err error
}

// Expect sets up expected params for Listener.Accept
func (m *mListenerMockAccept) Expect() *mListenerMockAccept {
	if m.mock.funcAccept != nil {
		m.mock.t.Fatalf("ListenerMock.Accept mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ListenerMockAcceptExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Listener.Accept
func (m *mListenerMockAccept) Return(c1 net.Conn, err error) *ListenerMock {
	if m.mock.funcAccept != nil {
		m.mock.t.Fatalf("ListenerMock.Accept mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ListenerMockAcceptExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ListenerMockAcceptResults{c1, err}
	return m.mock
}

//Set uses given function f to mock the Listener.Accept method
func (m *mListenerMockAccept) Set(f func() (c1 net.Conn, err error)) *ListenerMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Listener.Accept method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Listener.Accept method")
	}

	m.mock.funcAccept = f
	return m.mock
}

// Accept implements net.Listener
func (m *ListenerMock) Accept() (c1 net.Conn, err error) {
	atomic.AddUint64(&m.beforeAcceptCounter, 1)
	defer atomic.AddUint64(&m.afterAcceptCounter, 1)

	if m.AcceptMock.defaultExpectation != nil {
		atomic.AddUint64(&m.AcceptMock.defaultExpectation.Counter, 1)

		results := m.AcceptMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ListenerMock.Accept")
		}
		return (*results).c1, (*results).err
	}
	if m.funcAccept != nil {
		return m.funcAccept()
	}
	m.t.Fatalf("Unexpected call to ListenerMock.Accept.")
	return
}

// AcceptAfterCounter returns a count of finished ListenerMock.Accept invocations
func (m *ListenerMock) AcceptAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterAcceptCounter)
}

// AcceptBeforeCounter returns a count of ListenerMock.Accept invocations
func (m *ListenerMock) AcceptBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeAcceptCounter)
}

// MinimockAcceptDone returns true if the count of the Accept invocations corresponds
// the number of defined expectations
func (m *ListenerMock) MinimockAcceptDone() bool {
	for _, e := range m.AcceptMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AcceptMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAcceptCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAccept != nil && atomic.LoadUint64(&m.afterAcceptCounter) < 1 {
		return false
	}
	return true
}

// MinimockAcceptInspect logs each unmet expectation
func (m *ListenerMock) MinimockAcceptInspect() {
	for _, e := range m.AcceptMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ListenerMock.Accept")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AcceptMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAcceptCounter) < 1 {
		m.t.Error("Expected call to ListenerMock.Accept")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAccept != nil && atomic.LoadUint64(&m.afterAcceptCounter) < 1 {
		m.t.Error("Expected call to ListenerMock.Accept")
	}
}

type mListenerMockAddr struct {
	mock               *ListenerMock
	defaultExpectation *ListenerMockAddrExpectation
	expectations       []*ListenerMockAddrExpectation
}

// ListenerMockAddrExpectation specifies expectation struct of the Listener.Addr
type ListenerMockAddrExpectation struct {
	mock *ListenerMock

	results *ListenerMockAddrResults
	Counter uint64
}

// ListenerMockAddrResults contains results of the Listener.Addr
type ListenerMockAddrResults struct {
	a1 net.Addr
}

// Expect sets up expected params for Listener.Addr
func (m *mListenerMockAddr) Expect() *mListenerMockAddr {
	if m.mock.funcAddr != nil {
		m.mock.t.Fatalf("ListenerMock.Addr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ListenerMockAddrExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Listener.Addr
func (m *mListenerMockAddr) Return(a1 net.Addr) *ListenerMock {
	if m.mock.funcAddr != nil {
		m.mock.t.Fatalf("ListenerMock.Addr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ListenerMockAddrExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ListenerMockAddrResults{a1}
	return m.mock
}

//Set uses given function f to mock the Listener.Addr method
func (m *mListenerMockAddr) Set(f func() (a1 net.Addr)) *ListenerMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Listener.Addr method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Listener.Addr method")
	}

	m.mock.funcAddr = f
	return m.mock
}

// Addr implements net.Listener
func (m *ListenerMock) Addr() (a1 net.Addr) {
	atomic.AddUint64(&m.beforeAddrCounter, 1)
	defer atomic.AddUint64(&m.afterAddrCounter, 1)

	if m.AddrMock.defaultExpectation != nil {
		atomic.AddUint64(&m.AddrMock.defaultExpectation.Counter, 1)

		results := m.AddrMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ListenerMock.Addr")
		}
		return (*results).a1
	}
	if m.funcAddr != nil {
		return m.funcAddr()
	}
	m.t.Fatalf("Unexpected call to ListenerMock.Addr.")
	return
}

// AddrAfterCounter returns a count of finished ListenerMock.Addr invocations
func (m *ListenerMock) AddrAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterAddrCounter)
}

// AddrBeforeCounter returns a count of ListenerMock.Addr invocations
func (m *ListenerMock) AddrBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeAddrCounter)
}

// MinimockAddrDone returns true if the count of the Addr invocations corresponds
// the number of defined expectations
func (m *ListenerMock) MinimockAddrDone() bool {
	for _, e := range m.AddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAddrCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAddr != nil && atomic.LoadUint64(&m.afterAddrCounter) < 1 {
		return false
	}
	return true
}

// MinimockAddrInspect logs each unmet expectation
func (m *ListenerMock) MinimockAddrInspect() {
	for _, e := range m.AddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ListenerMock.Addr")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAddrCounter) < 1 {
		m.t.Error("Expected call to ListenerMock.Addr")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAddr != nil && atomic.LoadUint64(&m.afterAddrCounter) < 1 {
		m.t.Error("Expected call to ListenerMock.Addr")
	}
}

type mListenerMockClose struct {
	mock               *ListenerMock
	defaultExpectation *ListenerMockCloseExpectation
	expectations       []*ListenerMockCloseExpectation
}

// ListenerMockCloseExpectation specifies expectation struct of the Listener.Close
type ListenerMockCloseExpectation struct {
	mock *ListenerMock

	results *ListenerMockCloseResults
	Counter uint64
}

// ListenerMockCloseResults contains results of the Listener.Close
type ListenerMockCloseResults struct {
	err error
}

// Expect sets up expected params for Listener.Close
func (m *mListenerMockClose) Expect() *mListenerMockClose {
	if m.mock.funcClose != nil {
		m.mock.t.Fatalf("ListenerMock.Close mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ListenerMockCloseExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Listener.Close
func (m *mListenerMockClose) Return(err error) *ListenerMock {
	if m.mock.funcClose != nil {
		m.mock.t.Fatalf("ListenerMock.Close mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ListenerMockCloseExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ListenerMockCloseResults{err}
	return m.mock
}

//Set uses given function f to mock the Listener.Close method
func (m *mListenerMockClose) Set(f func() (err error)) *ListenerMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Listener.Close method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Listener.Close method")
	}

	m.mock.funcClose = f
	return m.mock
}

// Close implements net.Listener
func (m *ListenerMock) Close() (err error) {
	atomic.AddUint64(&m.beforeCloseCounter, 1)
	defer atomic.AddUint64(&m.afterCloseCounter, 1)

	if m.CloseMock.defaultExpectation != nil {
		atomic.AddUint64(&m.CloseMock.defaultExpectation.Counter, 1)

		results := m.CloseMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ListenerMock.Close")
		}
		return (*results).err
	}
	if m.funcClose != nil {
		return m.funcClose()
	}
	m.t.Fatalf("Unexpected call to ListenerMock.Close.")
	return
}

// CloseAfterCounter returns a count of finished ListenerMock.Close invocations
func (m *ListenerMock) CloseAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterCloseCounter)
}

// CloseBeforeCounter returns a count of ListenerMock.Close invocations
func (m *ListenerMock) CloseBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeCloseCounter)
}

// MinimockCloseDone returns true if the count of the Close invocations corresponds
// the number of defined expectations
func (m *ListenerMock) MinimockCloseDone() bool {
	for _, e := range m.CloseMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CloseMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcClose != nil && atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		return false
	}
	return true
}

// MinimockCloseInspect logs each unmet expectation
func (m *ListenerMock) MinimockCloseInspect() {
	for _, e := range m.CloseMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ListenerMock.Close")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CloseMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		m.t.Error("Expected call to ListenerMock.Close")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcClose != nil && atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		m.t.Error("Expected call to ListenerMock.Close")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ListenerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockAcceptInspect()

		m.MinimockAddrInspect()

		m.MinimockCloseInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ListenerMock) MinimockWait(timeout time.Duration) {
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

func (m *ListenerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAcceptDone() &&
		m.MinimockAddrDone() &&
		m.MinimockCloseDone()
}
