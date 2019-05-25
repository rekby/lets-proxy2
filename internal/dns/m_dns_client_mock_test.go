package dns

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/dns.mDNSClient -o ./m_dns_client_mock_test.go

import (
	"sync/atomic"
	"time"

	mdns "github.com/miekg/dns"

	"github.com/gojuno/minimock"
)

// MDNSClientMock implements mDNSClient
type MDNSClientMock struct {
	t minimock.Tester

	funcExchange          func(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error)
	afterExchangeCounter  uint64
	beforeExchangeCounter uint64
	ExchangeMock          mMDNSClientMockExchange
}

// NewMDNSClientMock returns a mock for mDNSClient
func NewMDNSClientMock(t minimock.Tester) *MDNSClientMock {
	m := &MDNSClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.ExchangeMock = mMDNSClientMockExchange{mock: m}

	return m
}

type mMDNSClientMockExchange struct {
	mock               *MDNSClientMock
	defaultExpectation *MDNSClientMockExchangeExpectation
	expectations       []*MDNSClientMockExchangeExpectation
}

// MDNSClientMockExchangeExpectation specifies expectation struct of the mDNSClient.Exchange
type MDNSClientMockExchangeExpectation struct {
	mock    *MDNSClientMock
	params  *MDNSClientMockExchangeParams
	results *MDNSClientMockExchangeResults
	Counter uint64
}

// MDNSClientMockExchangeParams contains parameters of the mDNSClient.Exchange
type MDNSClientMockExchangeParams struct {
	msg     *mdns.Msg
	address string
}

// MDNSClientMockExchangeResults contains results of the mDNSClient.Exchange
type MDNSClientMockExchangeResults struct {
	r   *mdns.Msg
	rtt time.Duration
	err error
}

// Expect sets up expected params for mDNSClient.Exchange
func (m *mMDNSClientMockExchange) Expect(msg *mdns.Msg, address string) *mMDNSClientMockExchange {
	if m.mock.funcExchange != nil {
		m.mock.t.Fatalf("MDNSClientMock.Exchange mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &MDNSClientMockExchangeExpectation{}
	}

	m.defaultExpectation.params = &MDNSClientMockExchangeParams{msg, address}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by mDNSClient.Exchange
func (m *mMDNSClientMockExchange) Return(r *mdns.Msg, rtt time.Duration, err error) *MDNSClientMock {
	if m.mock.funcExchange != nil {
		m.mock.t.Fatalf("MDNSClientMock.Exchange mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &MDNSClientMockExchangeExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &MDNSClientMockExchangeResults{r, rtt, err}
	return m.mock
}

//Set uses given function f to mock the mDNSClient.Exchange method
func (m *mMDNSClientMockExchange) Set(f func(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error)) *MDNSClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the mDNSClient.Exchange method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the mDNSClient.Exchange method")
	}

	m.mock.funcExchange = f
	return m.mock
}

// When sets expectation for the mDNSClient.Exchange which will trigger the result defined by the following
// Then helper
func (m *mMDNSClientMockExchange) When(msg *mdns.Msg, address string) *MDNSClientMockExchangeExpectation {
	if m.mock.funcExchange != nil {
		m.mock.t.Fatalf("MDNSClientMock.Exchange mock is already set by Set")
	}

	expectation := &MDNSClientMockExchangeExpectation{
		mock:   m.mock,
		params: &MDNSClientMockExchangeParams{msg, address},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up mDNSClient.Exchange return parameters for the expectation previously defined by the When method
func (e *MDNSClientMockExchangeExpectation) Then(r *mdns.Msg, rtt time.Duration, err error) *MDNSClientMock {
	e.results = &MDNSClientMockExchangeResults{r, rtt, err}
	return e.mock
}

// Exchange implements mDNSClient
func (m *MDNSClientMock) Exchange(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error) {
	atomic.AddUint64(&m.beforeExchangeCounter, 1)
	defer atomic.AddUint64(&m.afterExchangeCounter, 1)

	for _, e := range m.ExchangeMock.expectations {
		if minimock.Equal(*e.params, MDNSClientMockExchangeParams{msg, address}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.r, e.results.rtt, e.results.err
		}
	}

	if m.ExchangeMock.defaultExpectation != nil {
		atomic.AddUint64(&m.ExchangeMock.defaultExpectation.Counter, 1)
		want := m.ExchangeMock.defaultExpectation.params
		got := MDNSClientMockExchangeParams{msg, address}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("MDNSClientMock.Exchange got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.ExchangeMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the MDNSClientMock.Exchange")
		}
		return (*results).r, (*results).rtt, (*results).err
	}
	if m.funcExchange != nil {
		return m.funcExchange(msg, address)
	}
	m.t.Fatalf("Unexpected call to MDNSClientMock.Exchange. %v %v", msg, address)
	return
}

// ExchangeAfterCounter returns a count of finished MDNSClientMock.Exchange invocations
func (m *MDNSClientMock) ExchangeAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterExchangeCounter)
}

// ExchangeBeforeCounter returns a count of MDNSClientMock.Exchange invocations
func (m *MDNSClientMock) ExchangeBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeExchangeCounter)
}

// MinimockExchangeDone returns true if the count of the Exchange invocations corresponds
// the number of defined expectations
func (m *MDNSClientMock) MinimockExchangeDone() bool {
	for _, e := range m.ExchangeMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExchangeMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterExchangeCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExchange != nil && atomic.LoadUint64(&m.afterExchangeCounter) < 1 {
		return false
	}
	return true
}

// MinimockExchangeInspect logs each unmet expectation
func (m *MDNSClientMock) MinimockExchangeInspect() {
	for _, e := range m.ExchangeMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to MDNSClientMock.Exchange with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExchangeMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterExchangeCounter) < 1 {
		m.t.Errorf("Expected call to MDNSClientMock.Exchange with params: %#v", *m.ExchangeMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExchange != nil && atomic.LoadUint64(&m.afterExchangeCounter) < 1 {
		m.t.Error("Expected call to MDNSClientMock.Exchange")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *MDNSClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockExchangeInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *MDNSClientMock) MinimockWait(timeout time.Duration) {
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

func (m *MDNSClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockExchangeDone()
}
