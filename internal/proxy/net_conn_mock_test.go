package proxy

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i net.Conn -o ./internal/proxy/net_conn_mock_test.go

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

// ConnMock implements net.Conn
type ConnMock struct {
	t minimock.Tester

	funcClose          func() (err error)
	afterCloseCounter  uint64
	beforeCloseCounter uint64
	CloseMock          mConnMockClose

	funcLocalAddr          func() (a1 net.Addr)
	afterLocalAddrCounter  uint64
	beforeLocalAddrCounter uint64
	LocalAddrMock          mConnMockLocalAddr

	funcRead          func(b []byte) (n int, err error)
	afterReadCounter  uint64
	beforeReadCounter uint64
	ReadMock          mConnMockRead

	funcRemoteAddr          func() (a1 net.Addr)
	afterRemoteAddrCounter  uint64
	beforeRemoteAddrCounter uint64
	RemoteAddrMock          mConnMockRemoteAddr

	funcSetDeadline          func(t time.Time) (err error)
	afterSetDeadlineCounter  uint64
	beforeSetDeadlineCounter uint64
	SetDeadlineMock          mConnMockSetDeadline

	funcSetReadDeadline          func(t time.Time) (err error)
	afterSetReadDeadlineCounter  uint64
	beforeSetReadDeadlineCounter uint64
	SetReadDeadlineMock          mConnMockSetReadDeadline

	funcSetWriteDeadline          func(t time.Time) (err error)
	afterSetWriteDeadlineCounter  uint64
	beforeSetWriteDeadlineCounter uint64
	SetWriteDeadlineMock          mConnMockSetWriteDeadline

	funcWrite          func(b []byte) (n int, err error)
	afterWriteCounter  uint64
	beforeWriteCounter uint64
	WriteMock          mConnMockWrite
}

// NewConnMock returns a mock for net.Conn
func NewConnMock(t minimock.Tester) *ConnMock {
	m := &ConnMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.CloseMock = mConnMockClose{mock: m}
	m.LocalAddrMock = mConnMockLocalAddr{mock: m}
	m.ReadMock = mConnMockRead{mock: m}
	m.RemoteAddrMock = mConnMockRemoteAddr{mock: m}
	m.SetDeadlineMock = mConnMockSetDeadline{mock: m}
	m.SetReadDeadlineMock = mConnMockSetReadDeadline{mock: m}
	m.SetWriteDeadlineMock = mConnMockSetWriteDeadline{mock: m}
	m.WriteMock = mConnMockWrite{mock: m}

	return m
}

type mConnMockClose struct {
	mock               *ConnMock
	defaultExpectation *ConnMockCloseExpectation
	expectations       []*ConnMockCloseExpectation
}

// ConnMockCloseExpectation specifies expectation struct of the Conn.Close
type ConnMockCloseExpectation struct {
	mock *ConnMock

	results *ConnMockCloseResults
	Counter uint64
}

// ConnMockCloseResults contains results of the Conn.Close
type ConnMockCloseResults struct {
	err error
}

// Expect sets up expected params for Conn.Close
func (m *mConnMockClose) Expect() *mConnMockClose {
	if m.mock.funcClose != nil {
		m.mock.t.Fatalf("ConnMock.Close mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockCloseExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Conn.Close
func (m *mConnMockClose) Return(err error) *ConnMock {
	if m.mock.funcClose != nil {
		m.mock.t.Fatalf("ConnMock.Close mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockCloseExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockCloseResults{err}
	return m.mock
}

//Set uses given function f to mock the Conn.Close method
func (m *mConnMockClose) Set(f func() (err error)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.Close method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.Close method")
	}

	m.mock.funcClose = f
	return m.mock
}

// Close implements net.Conn
func (m *ConnMock) Close() (err error) {
	atomic.AddUint64(&m.beforeCloseCounter, 1)
	defer atomic.AddUint64(&m.afterCloseCounter, 1)

	if m.CloseMock.defaultExpectation != nil {
		atomic.AddUint64(&m.CloseMock.defaultExpectation.Counter, 1)

		results := m.CloseMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.Close")
		}
		return (*results).err
	}
	if m.funcClose != nil {
		return m.funcClose()
	}
	m.t.Fatalf("Unexpected call to ConnMock.Close.")
	return
}

// CloseAfterCounter returns a count of finished ConnMock.Close invocations
func (m *ConnMock) CloseAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterCloseCounter)
}

// CloseBeforeCounter returns a count of ConnMock.Close invocations
func (m *ConnMock) CloseBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeCloseCounter)
}

// MinimockCloseDone returns true if the count of the Close invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockCloseDone() bool {
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
func (m *ConnMock) MinimockCloseInspect() {
	for _, e := range m.CloseMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ConnMock.Close")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CloseMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		m.t.Error("Expected call to ConnMock.Close")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcClose != nil && atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		m.t.Error("Expected call to ConnMock.Close")
	}
}

type mConnMockLocalAddr struct {
	mock               *ConnMock
	defaultExpectation *ConnMockLocalAddrExpectation
	expectations       []*ConnMockLocalAddrExpectation
}

// ConnMockLocalAddrExpectation specifies expectation struct of the Conn.LocalAddr
type ConnMockLocalAddrExpectation struct {
	mock *ConnMock

	results *ConnMockLocalAddrResults
	Counter uint64
}

// ConnMockLocalAddrResults contains results of the Conn.LocalAddr
type ConnMockLocalAddrResults struct {
	a1 net.Addr
}

// Expect sets up expected params for Conn.LocalAddr
func (m *mConnMockLocalAddr) Expect() *mConnMockLocalAddr {
	if m.mock.funcLocalAddr != nil {
		m.mock.t.Fatalf("ConnMock.LocalAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockLocalAddrExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Conn.LocalAddr
func (m *mConnMockLocalAddr) Return(a1 net.Addr) *ConnMock {
	if m.mock.funcLocalAddr != nil {
		m.mock.t.Fatalf("ConnMock.LocalAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockLocalAddrExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockLocalAddrResults{a1}
	return m.mock
}

//Set uses given function f to mock the Conn.LocalAddr method
func (m *mConnMockLocalAddr) Set(f func() (a1 net.Addr)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.LocalAddr method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.LocalAddr method")
	}

	m.mock.funcLocalAddr = f
	return m.mock
}

// LocalAddr implements net.Conn
func (m *ConnMock) LocalAddr() (a1 net.Addr) {
	atomic.AddUint64(&m.beforeLocalAddrCounter, 1)
	defer atomic.AddUint64(&m.afterLocalAddrCounter, 1)

	if m.LocalAddrMock.defaultExpectation != nil {
		atomic.AddUint64(&m.LocalAddrMock.defaultExpectation.Counter, 1)

		results := m.LocalAddrMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.LocalAddr")
		}
		return (*results).a1
	}
	if m.funcLocalAddr != nil {
		return m.funcLocalAddr()
	}
	m.t.Fatalf("Unexpected call to ConnMock.LocalAddr.")
	return
}

// LocalAddrAfterCounter returns a count of finished ConnMock.LocalAddr invocations
func (m *ConnMock) LocalAddrAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterLocalAddrCounter)
}

// LocalAddrBeforeCounter returns a count of ConnMock.LocalAddr invocations
func (m *ConnMock) LocalAddrBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeLocalAddrCounter)
}

// MinimockLocalAddrDone returns true if the count of the LocalAddr invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockLocalAddrDone() bool {
	for _, e := range m.LocalAddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LocalAddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterLocalAddrCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLocalAddr != nil && atomic.LoadUint64(&m.afterLocalAddrCounter) < 1 {
		return false
	}
	return true
}

// MinimockLocalAddrInspect logs each unmet expectation
func (m *ConnMock) MinimockLocalAddrInspect() {
	for _, e := range m.LocalAddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ConnMock.LocalAddr")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LocalAddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterLocalAddrCounter) < 1 {
		m.t.Error("Expected call to ConnMock.LocalAddr")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLocalAddr != nil && atomic.LoadUint64(&m.afterLocalAddrCounter) < 1 {
		m.t.Error("Expected call to ConnMock.LocalAddr")
	}
}

type mConnMockRead struct {
	mock               *ConnMock
	defaultExpectation *ConnMockReadExpectation
	expectations       []*ConnMockReadExpectation
}

// ConnMockReadExpectation specifies expectation struct of the Conn.Read
type ConnMockReadExpectation struct {
	mock    *ConnMock
	params  *ConnMockReadParams
	results *ConnMockReadResults
	Counter uint64
}

// ConnMockReadParams contains parameters of the Conn.Read
type ConnMockReadParams struct {
	b []byte
}

// ConnMockReadResults contains results of the Conn.Read
type ConnMockReadResults struct {
	n   int
	err error
}

// Expect sets up expected params for Conn.Read
func (m *mConnMockRead) Expect(b []byte) *mConnMockRead {
	if m.mock.funcRead != nil {
		m.mock.t.Fatalf("ConnMock.Read mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockReadExpectation{}
	}

	m.defaultExpectation.params = &ConnMockReadParams{b}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Conn.Read
func (m *mConnMockRead) Return(n int, err error) *ConnMock {
	if m.mock.funcRead != nil {
		m.mock.t.Fatalf("ConnMock.Read mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockReadExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockReadResults{n, err}
	return m.mock
}

//Set uses given function f to mock the Conn.Read method
func (m *mConnMockRead) Set(f func(b []byte) (n int, err error)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.Read method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.Read method")
	}

	m.mock.funcRead = f
	return m.mock
}

// When sets expectation for the Conn.Read which will trigger the result defined by the following
// Then helper
func (m *mConnMockRead) When(b []byte) *ConnMockReadExpectation {
	if m.mock.funcRead != nil {
		m.mock.t.Fatalf("ConnMock.Read mock is already set by Set")
	}

	expectation := &ConnMockReadExpectation{
		mock:   m.mock,
		params: &ConnMockReadParams{b},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Conn.Read return parameters for the expectation previously defined by the When method
func (e *ConnMockReadExpectation) Then(n int, err error) *ConnMock {
	e.results = &ConnMockReadResults{n, err}
	return e.mock
}

// Read implements net.Conn
func (m *ConnMock) Read(b []byte) (n int, err error) {
	atomic.AddUint64(&m.beforeReadCounter, 1)
	defer atomic.AddUint64(&m.afterReadCounter, 1)

	for _, e := range m.ReadMock.expectations {
		if minimock.Equal(*e.params, ConnMockReadParams{b}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.n, e.results.err
		}
	}

	if m.ReadMock.defaultExpectation != nil {
		atomic.AddUint64(&m.ReadMock.defaultExpectation.Counter, 1)
		want := m.ReadMock.defaultExpectation.params
		got := ConnMockReadParams{b}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ConnMock.Read got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.ReadMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.Read")
		}
		return (*results).n, (*results).err
	}
	if m.funcRead != nil {
		return m.funcRead(b)
	}
	m.t.Fatalf("Unexpected call to ConnMock.Read. %v", b)
	return
}

// ReadAfterCounter returns a count of finished ConnMock.Read invocations
func (m *ConnMock) ReadAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterReadCounter)
}

// ReadBeforeCounter returns a count of ConnMock.Read invocations
func (m *ConnMock) ReadBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeReadCounter)
}

// MinimockReadDone returns true if the count of the Read invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockReadDone() bool {
	for _, e := range m.ReadMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ReadMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRead != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		return false
	}
	return true
}

// MinimockReadInspect logs each unmet expectation
func (m *ConnMock) MinimockReadInspect() {
	for _, e := range m.ReadMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ConnMock.Read with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ReadMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		m.t.Errorf("Expected call to ConnMock.Read with params: %#v", *m.ReadMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRead != nil && atomic.LoadUint64(&m.afterReadCounter) < 1 {
		m.t.Error("Expected call to ConnMock.Read")
	}
}

type mConnMockRemoteAddr struct {
	mock               *ConnMock
	defaultExpectation *ConnMockRemoteAddrExpectation
	expectations       []*ConnMockRemoteAddrExpectation
}

// ConnMockRemoteAddrExpectation specifies expectation struct of the Conn.RemoteAddr
type ConnMockRemoteAddrExpectation struct {
	mock *ConnMock

	results *ConnMockRemoteAddrResults
	Counter uint64
}

// ConnMockRemoteAddrResults contains results of the Conn.RemoteAddr
type ConnMockRemoteAddrResults struct {
	a1 net.Addr
}

// Expect sets up expected params for Conn.RemoteAddr
func (m *mConnMockRemoteAddr) Expect() *mConnMockRemoteAddr {
	if m.mock.funcRemoteAddr != nil {
		m.mock.t.Fatalf("ConnMock.RemoteAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockRemoteAddrExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Conn.RemoteAddr
func (m *mConnMockRemoteAddr) Return(a1 net.Addr) *ConnMock {
	if m.mock.funcRemoteAddr != nil {
		m.mock.t.Fatalf("ConnMock.RemoteAddr mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockRemoteAddrExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockRemoteAddrResults{a1}
	return m.mock
}

//Set uses given function f to mock the Conn.RemoteAddr method
func (m *mConnMockRemoteAddr) Set(f func() (a1 net.Addr)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.RemoteAddr method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.RemoteAddr method")
	}

	m.mock.funcRemoteAddr = f
	return m.mock
}

// RemoteAddr implements net.Conn
func (m *ConnMock) RemoteAddr() (a1 net.Addr) {
	atomic.AddUint64(&m.beforeRemoteAddrCounter, 1)
	defer atomic.AddUint64(&m.afterRemoteAddrCounter, 1)

	if m.RemoteAddrMock.defaultExpectation != nil {
		atomic.AddUint64(&m.RemoteAddrMock.defaultExpectation.Counter, 1)

		results := m.RemoteAddrMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.RemoteAddr")
		}
		return (*results).a1
	}
	if m.funcRemoteAddr != nil {
		return m.funcRemoteAddr()
	}
	m.t.Fatalf("Unexpected call to ConnMock.RemoteAddr.")
	return
}

// RemoteAddrAfterCounter returns a count of finished ConnMock.RemoteAddr invocations
func (m *ConnMock) RemoteAddrAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterRemoteAddrCounter)
}

// RemoteAddrBeforeCounter returns a count of ConnMock.RemoteAddr invocations
func (m *ConnMock) RemoteAddrBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeRemoteAddrCounter)
}

// MinimockRemoteAddrDone returns true if the count of the RemoteAddr invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockRemoteAddrDone() bool {
	for _, e := range m.RemoteAddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoteAddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRemoteAddrCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemoteAddr != nil && atomic.LoadUint64(&m.afterRemoteAddrCounter) < 1 {
		return false
	}
	return true
}

// MinimockRemoteAddrInspect logs each unmet expectation
func (m *ConnMock) MinimockRemoteAddrInspect() {
	for _, e := range m.RemoteAddrMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ConnMock.RemoteAddr")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoteAddrMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRemoteAddrCounter) < 1 {
		m.t.Error("Expected call to ConnMock.RemoteAddr")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemoteAddr != nil && atomic.LoadUint64(&m.afterRemoteAddrCounter) < 1 {
		m.t.Error("Expected call to ConnMock.RemoteAddr")
	}
}

type mConnMockSetDeadline struct {
	mock               *ConnMock
	defaultExpectation *ConnMockSetDeadlineExpectation
	expectations       []*ConnMockSetDeadlineExpectation
}

// ConnMockSetDeadlineExpectation specifies expectation struct of the Conn.SetDeadline
type ConnMockSetDeadlineExpectation struct {
	mock    *ConnMock
	params  *ConnMockSetDeadlineParams
	results *ConnMockSetDeadlineResults
	Counter uint64
}

// ConnMockSetDeadlineParams contains parameters of the Conn.SetDeadline
type ConnMockSetDeadlineParams struct {
	t time.Time
}

// ConnMockSetDeadlineResults contains results of the Conn.SetDeadline
type ConnMockSetDeadlineResults struct {
	err error
}

// Expect sets up expected params for Conn.SetDeadline
func (m *mConnMockSetDeadline) Expect(t time.Time) *mConnMockSetDeadline {
	if m.mock.funcSetDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetDeadline mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockSetDeadlineExpectation{}
	}

	m.defaultExpectation.params = &ConnMockSetDeadlineParams{t}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Conn.SetDeadline
func (m *mConnMockSetDeadline) Return(err error) *ConnMock {
	if m.mock.funcSetDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetDeadline mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockSetDeadlineExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockSetDeadlineResults{err}
	return m.mock
}

//Set uses given function f to mock the Conn.SetDeadline method
func (m *mConnMockSetDeadline) Set(f func(t time.Time) (err error)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.SetDeadline method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.SetDeadline method")
	}

	m.mock.funcSetDeadline = f
	return m.mock
}

// When sets expectation for the Conn.SetDeadline which will trigger the result defined by the following
// Then helper
func (m *mConnMockSetDeadline) When(t time.Time) *ConnMockSetDeadlineExpectation {
	if m.mock.funcSetDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetDeadline mock is already set by Set")
	}

	expectation := &ConnMockSetDeadlineExpectation{
		mock:   m.mock,
		params: &ConnMockSetDeadlineParams{t},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Conn.SetDeadline return parameters for the expectation previously defined by the When method
func (e *ConnMockSetDeadlineExpectation) Then(err error) *ConnMock {
	e.results = &ConnMockSetDeadlineResults{err}
	return e.mock
}

// SetDeadline implements net.Conn
func (m *ConnMock) SetDeadline(t time.Time) (err error) {
	atomic.AddUint64(&m.beforeSetDeadlineCounter, 1)
	defer atomic.AddUint64(&m.afterSetDeadlineCounter, 1)

	for _, e := range m.SetDeadlineMock.expectations {
		if minimock.Equal(*e.params, ConnMockSetDeadlineParams{t}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.SetDeadlineMock.defaultExpectation != nil {
		atomic.AddUint64(&m.SetDeadlineMock.defaultExpectation.Counter, 1)
		want := m.SetDeadlineMock.defaultExpectation.params
		got := ConnMockSetDeadlineParams{t}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ConnMock.SetDeadline got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.SetDeadlineMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.SetDeadline")
		}
		return (*results).err
	}
	if m.funcSetDeadline != nil {
		return m.funcSetDeadline(t)
	}
	m.t.Fatalf("Unexpected call to ConnMock.SetDeadline. %v", t)
	return
}

// SetDeadlineAfterCounter returns a count of finished ConnMock.SetDeadline invocations
func (m *ConnMock) SetDeadlineAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterSetDeadlineCounter)
}

// SetDeadlineBeforeCounter returns a count of ConnMock.SetDeadline invocations
func (m *ConnMock) SetDeadlineBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeSetDeadlineCounter)
}

// MinimockSetDeadlineDone returns true if the count of the SetDeadline invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockSetDeadlineDone() bool {
	for _, e := range m.SetDeadlineMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SetDeadlineMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterSetDeadlineCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSetDeadline != nil && atomic.LoadUint64(&m.afterSetDeadlineCounter) < 1 {
		return false
	}
	return true
}

// MinimockSetDeadlineInspect logs each unmet expectation
func (m *ConnMock) MinimockSetDeadlineInspect() {
	for _, e := range m.SetDeadlineMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ConnMock.SetDeadline with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SetDeadlineMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterSetDeadlineCounter) < 1 {
		m.t.Errorf("Expected call to ConnMock.SetDeadline with params: %#v", *m.SetDeadlineMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSetDeadline != nil && atomic.LoadUint64(&m.afterSetDeadlineCounter) < 1 {
		m.t.Error("Expected call to ConnMock.SetDeadline")
	}
}

type mConnMockSetReadDeadline struct {
	mock               *ConnMock
	defaultExpectation *ConnMockSetReadDeadlineExpectation
	expectations       []*ConnMockSetReadDeadlineExpectation
}

// ConnMockSetReadDeadlineExpectation specifies expectation struct of the Conn.SetReadDeadline
type ConnMockSetReadDeadlineExpectation struct {
	mock    *ConnMock
	params  *ConnMockSetReadDeadlineParams
	results *ConnMockSetReadDeadlineResults
	Counter uint64
}

// ConnMockSetReadDeadlineParams contains parameters of the Conn.SetReadDeadline
type ConnMockSetReadDeadlineParams struct {
	t time.Time
}

// ConnMockSetReadDeadlineResults contains results of the Conn.SetReadDeadline
type ConnMockSetReadDeadlineResults struct {
	err error
}

// Expect sets up expected params for Conn.SetReadDeadline
func (m *mConnMockSetReadDeadline) Expect(t time.Time) *mConnMockSetReadDeadline {
	if m.mock.funcSetReadDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetReadDeadline mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockSetReadDeadlineExpectation{}
	}

	m.defaultExpectation.params = &ConnMockSetReadDeadlineParams{t}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Conn.SetReadDeadline
func (m *mConnMockSetReadDeadline) Return(err error) *ConnMock {
	if m.mock.funcSetReadDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetReadDeadline mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockSetReadDeadlineExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockSetReadDeadlineResults{err}
	return m.mock
}

//Set uses given function f to mock the Conn.SetReadDeadline method
func (m *mConnMockSetReadDeadline) Set(f func(t time.Time) (err error)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.SetReadDeadline method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.SetReadDeadline method")
	}

	m.mock.funcSetReadDeadline = f
	return m.mock
}

// When sets expectation for the Conn.SetReadDeadline which will trigger the result defined by the following
// Then helper
func (m *mConnMockSetReadDeadline) When(t time.Time) *ConnMockSetReadDeadlineExpectation {
	if m.mock.funcSetReadDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetReadDeadline mock is already set by Set")
	}

	expectation := &ConnMockSetReadDeadlineExpectation{
		mock:   m.mock,
		params: &ConnMockSetReadDeadlineParams{t},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Conn.SetReadDeadline return parameters for the expectation previously defined by the When method
func (e *ConnMockSetReadDeadlineExpectation) Then(err error) *ConnMock {
	e.results = &ConnMockSetReadDeadlineResults{err}
	return e.mock
}

// SetReadDeadline implements net.Conn
func (m *ConnMock) SetReadDeadline(t time.Time) (err error) {
	atomic.AddUint64(&m.beforeSetReadDeadlineCounter, 1)
	defer atomic.AddUint64(&m.afterSetReadDeadlineCounter, 1)

	for _, e := range m.SetReadDeadlineMock.expectations {
		if minimock.Equal(*e.params, ConnMockSetReadDeadlineParams{t}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.SetReadDeadlineMock.defaultExpectation != nil {
		atomic.AddUint64(&m.SetReadDeadlineMock.defaultExpectation.Counter, 1)
		want := m.SetReadDeadlineMock.defaultExpectation.params
		got := ConnMockSetReadDeadlineParams{t}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ConnMock.SetReadDeadline got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.SetReadDeadlineMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.SetReadDeadline")
		}
		return (*results).err
	}
	if m.funcSetReadDeadline != nil {
		return m.funcSetReadDeadline(t)
	}
	m.t.Fatalf("Unexpected call to ConnMock.SetReadDeadline. %v", t)
	return
}

// SetReadDeadlineAfterCounter returns a count of finished ConnMock.SetReadDeadline invocations
func (m *ConnMock) SetReadDeadlineAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterSetReadDeadlineCounter)
}

// SetReadDeadlineBeforeCounter returns a count of ConnMock.SetReadDeadline invocations
func (m *ConnMock) SetReadDeadlineBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeSetReadDeadlineCounter)
}

// MinimockSetReadDeadlineDone returns true if the count of the SetReadDeadline invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockSetReadDeadlineDone() bool {
	for _, e := range m.SetReadDeadlineMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SetReadDeadlineMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterSetReadDeadlineCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSetReadDeadline != nil && atomic.LoadUint64(&m.afterSetReadDeadlineCounter) < 1 {
		return false
	}
	return true
}

// MinimockSetReadDeadlineInspect logs each unmet expectation
func (m *ConnMock) MinimockSetReadDeadlineInspect() {
	for _, e := range m.SetReadDeadlineMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ConnMock.SetReadDeadline with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SetReadDeadlineMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterSetReadDeadlineCounter) < 1 {
		m.t.Errorf("Expected call to ConnMock.SetReadDeadline with params: %#v", *m.SetReadDeadlineMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSetReadDeadline != nil && atomic.LoadUint64(&m.afterSetReadDeadlineCounter) < 1 {
		m.t.Error("Expected call to ConnMock.SetReadDeadline")
	}
}

type mConnMockSetWriteDeadline struct {
	mock               *ConnMock
	defaultExpectation *ConnMockSetWriteDeadlineExpectation
	expectations       []*ConnMockSetWriteDeadlineExpectation
}

// ConnMockSetWriteDeadlineExpectation specifies expectation struct of the Conn.SetWriteDeadline
type ConnMockSetWriteDeadlineExpectation struct {
	mock    *ConnMock
	params  *ConnMockSetWriteDeadlineParams
	results *ConnMockSetWriteDeadlineResults
	Counter uint64
}

// ConnMockSetWriteDeadlineParams contains parameters of the Conn.SetWriteDeadline
type ConnMockSetWriteDeadlineParams struct {
	t time.Time
}

// ConnMockSetWriteDeadlineResults contains results of the Conn.SetWriteDeadline
type ConnMockSetWriteDeadlineResults struct {
	err error
}

// Expect sets up expected params for Conn.SetWriteDeadline
func (m *mConnMockSetWriteDeadline) Expect(t time.Time) *mConnMockSetWriteDeadline {
	if m.mock.funcSetWriteDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetWriteDeadline mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockSetWriteDeadlineExpectation{}
	}

	m.defaultExpectation.params = &ConnMockSetWriteDeadlineParams{t}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Conn.SetWriteDeadline
func (m *mConnMockSetWriteDeadline) Return(err error) *ConnMock {
	if m.mock.funcSetWriteDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetWriteDeadline mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockSetWriteDeadlineExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockSetWriteDeadlineResults{err}
	return m.mock
}

//Set uses given function f to mock the Conn.SetWriteDeadline method
func (m *mConnMockSetWriteDeadline) Set(f func(t time.Time) (err error)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.SetWriteDeadline method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.SetWriteDeadline method")
	}

	m.mock.funcSetWriteDeadline = f
	return m.mock
}

// When sets expectation for the Conn.SetWriteDeadline which will trigger the result defined by the following
// Then helper
func (m *mConnMockSetWriteDeadline) When(t time.Time) *ConnMockSetWriteDeadlineExpectation {
	if m.mock.funcSetWriteDeadline != nil {
		m.mock.t.Fatalf("ConnMock.SetWriteDeadline mock is already set by Set")
	}

	expectation := &ConnMockSetWriteDeadlineExpectation{
		mock:   m.mock,
		params: &ConnMockSetWriteDeadlineParams{t},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Conn.SetWriteDeadline return parameters for the expectation previously defined by the When method
func (e *ConnMockSetWriteDeadlineExpectation) Then(err error) *ConnMock {
	e.results = &ConnMockSetWriteDeadlineResults{err}
	return e.mock
}

// SetWriteDeadline implements net.Conn
func (m *ConnMock) SetWriteDeadline(t time.Time) (err error) {
	atomic.AddUint64(&m.beforeSetWriteDeadlineCounter, 1)
	defer atomic.AddUint64(&m.afterSetWriteDeadlineCounter, 1)

	for _, e := range m.SetWriteDeadlineMock.expectations {
		if minimock.Equal(*e.params, ConnMockSetWriteDeadlineParams{t}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.SetWriteDeadlineMock.defaultExpectation != nil {
		atomic.AddUint64(&m.SetWriteDeadlineMock.defaultExpectation.Counter, 1)
		want := m.SetWriteDeadlineMock.defaultExpectation.params
		got := ConnMockSetWriteDeadlineParams{t}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ConnMock.SetWriteDeadline got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.SetWriteDeadlineMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.SetWriteDeadline")
		}
		return (*results).err
	}
	if m.funcSetWriteDeadline != nil {
		return m.funcSetWriteDeadline(t)
	}
	m.t.Fatalf("Unexpected call to ConnMock.SetWriteDeadline. %v", t)
	return
}

// SetWriteDeadlineAfterCounter returns a count of finished ConnMock.SetWriteDeadline invocations
func (m *ConnMock) SetWriteDeadlineAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterSetWriteDeadlineCounter)
}

// SetWriteDeadlineBeforeCounter returns a count of ConnMock.SetWriteDeadline invocations
func (m *ConnMock) SetWriteDeadlineBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeSetWriteDeadlineCounter)
}

// MinimockSetWriteDeadlineDone returns true if the count of the SetWriteDeadline invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockSetWriteDeadlineDone() bool {
	for _, e := range m.SetWriteDeadlineMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SetWriteDeadlineMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterSetWriteDeadlineCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSetWriteDeadline != nil && atomic.LoadUint64(&m.afterSetWriteDeadlineCounter) < 1 {
		return false
	}
	return true
}

// MinimockSetWriteDeadlineInspect logs each unmet expectation
func (m *ConnMock) MinimockSetWriteDeadlineInspect() {
	for _, e := range m.SetWriteDeadlineMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ConnMock.SetWriteDeadline with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SetWriteDeadlineMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterSetWriteDeadlineCounter) < 1 {
		m.t.Errorf("Expected call to ConnMock.SetWriteDeadline with params: %#v", *m.SetWriteDeadlineMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSetWriteDeadline != nil && atomic.LoadUint64(&m.afterSetWriteDeadlineCounter) < 1 {
		m.t.Error("Expected call to ConnMock.SetWriteDeadline")
	}
}

type mConnMockWrite struct {
	mock               *ConnMock
	defaultExpectation *ConnMockWriteExpectation
	expectations       []*ConnMockWriteExpectation
}

// ConnMockWriteExpectation specifies expectation struct of the Conn.Write
type ConnMockWriteExpectation struct {
	mock    *ConnMock
	params  *ConnMockWriteParams
	results *ConnMockWriteResults
	Counter uint64
}

// ConnMockWriteParams contains parameters of the Conn.Write
type ConnMockWriteParams struct {
	b []byte
}

// ConnMockWriteResults contains results of the Conn.Write
type ConnMockWriteResults struct {
	n   int
	err error
}

// Expect sets up expected params for Conn.Write
func (m *mConnMockWrite) Expect(b []byte) *mConnMockWrite {
	if m.mock.funcWrite != nil {
		m.mock.t.Fatalf("ConnMock.Write mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockWriteExpectation{}
	}

	m.defaultExpectation.params = &ConnMockWriteParams{b}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Conn.Write
func (m *mConnMockWrite) Return(n int, err error) *ConnMock {
	if m.mock.funcWrite != nil {
		m.mock.t.Fatalf("ConnMock.Write mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ConnMockWriteExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ConnMockWriteResults{n, err}
	return m.mock
}

//Set uses given function f to mock the Conn.Write method
func (m *mConnMockWrite) Set(f func(b []byte) (n int, err error)) *ConnMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Conn.Write method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Conn.Write method")
	}

	m.mock.funcWrite = f
	return m.mock
}

// When sets expectation for the Conn.Write which will trigger the result defined by the following
// Then helper
func (m *mConnMockWrite) When(b []byte) *ConnMockWriteExpectation {
	if m.mock.funcWrite != nil {
		m.mock.t.Fatalf("ConnMock.Write mock is already set by Set")
	}

	expectation := &ConnMockWriteExpectation{
		mock:   m.mock,
		params: &ConnMockWriteParams{b},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Conn.Write return parameters for the expectation previously defined by the When method
func (e *ConnMockWriteExpectation) Then(n int, err error) *ConnMock {
	e.results = &ConnMockWriteResults{n, err}
	return e.mock
}

// Write implements net.Conn
func (m *ConnMock) Write(b []byte) (n int, err error) {
	atomic.AddUint64(&m.beforeWriteCounter, 1)
	defer atomic.AddUint64(&m.afterWriteCounter, 1)

	for _, e := range m.WriteMock.expectations {
		if minimock.Equal(*e.params, ConnMockWriteParams{b}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.n, e.results.err
		}
	}

	if m.WriteMock.defaultExpectation != nil {
		atomic.AddUint64(&m.WriteMock.defaultExpectation.Counter, 1)
		want := m.WriteMock.defaultExpectation.params
		got := ConnMockWriteParams{b}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ConnMock.Write got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.WriteMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ConnMock.Write")
		}
		return (*results).n, (*results).err
	}
	if m.funcWrite != nil {
		return m.funcWrite(b)
	}
	m.t.Fatalf("Unexpected call to ConnMock.Write. %v", b)
	return
}

// WriteAfterCounter returns a count of finished ConnMock.Write invocations
func (m *ConnMock) WriteAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterWriteCounter)
}

// WriteBeforeCounter returns a count of ConnMock.Write invocations
func (m *ConnMock) WriteBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeWriteCounter)
}

// MinimockWriteDone returns true if the count of the Write invocations corresponds
// the number of defined expectations
func (m *ConnMock) MinimockWriteDone() bool {
	for _, e := range m.WriteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WriteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWrite != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		return false
	}
	return true
}

// MinimockWriteInspect logs each unmet expectation
func (m *ConnMock) MinimockWriteInspect() {
	for _, e := range m.WriteMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ConnMock.Write with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WriteMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		m.t.Errorf("Expected call to ConnMock.Write with params: %#v", *m.WriteMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWrite != nil && atomic.LoadUint64(&m.afterWriteCounter) < 1 {
		m.t.Error("Expected call to ConnMock.Write")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ConnMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockCloseInspect()

		m.MinimockLocalAddrInspect()

		m.MinimockReadInspect()

		m.MinimockRemoteAddrInspect()

		m.MinimockSetDeadlineInspect()

		m.MinimockSetReadDeadlineInspect()

		m.MinimockSetWriteDeadlineInspect()

		m.MinimockWriteInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ConnMock) MinimockWait(timeout time.Duration) {
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

func (m *ConnMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCloseDone() &&
		m.MinimockLocalAddrDone() &&
		m.MinimockReadDone() &&
		m.MinimockRemoteAddrDone() &&
		m.MinimockSetDeadlineDone() &&
		m.MinimockSetReadDeadlineDone() &&
		m.MinimockSetWriteDeadlineDone() &&
		m.MinimockWriteDone()
}
