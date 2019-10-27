package cert_manager

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cert_manager.AcmeClient -o ./acme_client_mock_test.go

import (
	"sync/atomic"
	"time"

	"context"

	"crypto/tls"

	"golang.org/x/crypto/acme"

	"github.com/gojuno/minimock"
)

// AcmeClientMock implements AcmeClient
type AcmeClientMock struct {
	t minimock.Tester

	funcAccept          func(ctx context.Context, chal *acme.Challenge) (cp1 *acme.Challenge, err error)
	afterAcceptCounter  uint64
	beforeAcceptCounter uint64
	AcceptMock          mAcmeClientMockAccept

	funcAuthorizeOrder          func(ctx context.Context, id []acme.AuthzID, opt ...acme.OrderOption) (op1 *acme.Order, err error)
	afterAuthorizeOrderCounter  uint64
	beforeAuthorizeOrderCounter uint64
	AuthorizeOrderMock          mAcmeClientMockAuthorizeOrder

	funcCreateOrderCert          func(ctx context.Context, url string, csr []byte, bundle bool) (der [][]byte, certURL string, err error)
	afterCreateOrderCertCounter  uint64
	beforeCreateOrderCertCounter uint64
	CreateOrderCertMock          mAcmeClientMockCreateOrderCert

	funcGetAuthorization          func(ctx context.Context, url string) (ap1 *acme.Authorization, err error)
	afterGetAuthorizationCounter  uint64
	beforeGetAuthorizationCounter uint64
	GetAuthorizationMock          mAcmeClientMockGetAuthorization

	funcHTTP01ChallengeResponse          func(token string) (s1 string, err error)
	afterHTTP01ChallengeResponseCounter  uint64
	beforeHTTP01ChallengeResponseCounter uint64
	HTTP01ChallengeResponseMock          mAcmeClientMockHTTP01ChallengeResponse

	funcRevokeAuthorization          func(ctx context.Context, url string) (err error)
	afterRevokeAuthorizationCounter  uint64
	beforeRevokeAuthorizationCounter uint64
	RevokeAuthorizationMock          mAcmeClientMockRevokeAuthorization

	funcTLSALPN01ChallengeCert          func(token string, domain string, opt ...acme.CertOption) (cert tls.Certificate, err error)
	afterTLSALPN01ChallengeCertCounter  uint64
	beforeTLSALPN01ChallengeCertCounter uint64
	TLSALPN01ChallengeCertMock          mAcmeClientMockTLSALPN01ChallengeCert

	funcWaitAuthorization          func(ctx context.Context, url string) (ap1 *acme.Authorization, err error)
	afterWaitAuthorizationCounter  uint64
	beforeWaitAuthorizationCounter uint64
	WaitAuthorizationMock          mAcmeClientMockWaitAuthorization

	funcWaitOrder          func(ctx context.Context, url string) (op1 *acme.Order, err error)
	afterWaitOrderCounter  uint64
	beforeWaitOrderCounter uint64
	WaitOrderMock          mAcmeClientMockWaitOrder
}

// NewAcmeClientMock returns a mock for AcmeClient
func NewAcmeClientMock(t minimock.Tester) *AcmeClientMock {
	m := &AcmeClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.AcceptMock = mAcmeClientMockAccept{mock: m}
	m.AuthorizeOrderMock = mAcmeClientMockAuthorizeOrder{mock: m}
	m.CreateOrderCertMock = mAcmeClientMockCreateOrderCert{mock: m}
	m.GetAuthorizationMock = mAcmeClientMockGetAuthorization{mock: m}
	m.HTTP01ChallengeResponseMock = mAcmeClientMockHTTP01ChallengeResponse{mock: m}
	m.RevokeAuthorizationMock = mAcmeClientMockRevokeAuthorization{mock: m}
	m.TLSALPN01ChallengeCertMock = mAcmeClientMockTLSALPN01ChallengeCert{mock: m}
	m.WaitAuthorizationMock = mAcmeClientMockWaitAuthorization{mock: m}
	m.WaitOrderMock = mAcmeClientMockWaitOrder{mock: m}

	return m
}

type mAcmeClientMockAccept struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockAcceptExpectation
	expectations       []*AcmeClientMockAcceptExpectation
}

// AcmeClientMockAcceptExpectation specifies expectation struct of the AcmeClient.Accept
type AcmeClientMockAcceptExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockAcceptParams
	results *AcmeClientMockAcceptResults
	Counter uint64
}

// AcmeClientMockAcceptParams contains parameters of the AcmeClient.Accept
type AcmeClientMockAcceptParams struct {
	ctx  context.Context
	chal *acme.Challenge
}

// AcmeClientMockAcceptResults contains results of the AcmeClient.Accept
type AcmeClientMockAcceptResults struct {
	cp1 *acme.Challenge
	err error
}

// Expect sets up expected params for AcmeClient.Accept
func (m *mAcmeClientMockAccept) Expect(ctx context.Context, chal *acme.Challenge) *mAcmeClientMockAccept {
	if m.mock.funcAccept != nil {
		m.mock.t.Fatalf("AcmeClientMock.Accept mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockAcceptExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockAcceptParams{ctx, chal}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.Accept
func (m *mAcmeClientMockAccept) Return(cp1 *acme.Challenge, err error) *AcmeClientMock {
	if m.mock.funcAccept != nil {
		m.mock.t.Fatalf("AcmeClientMock.Accept mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockAcceptExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockAcceptResults{cp1, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.Accept method
func (m *mAcmeClientMockAccept) Set(f func(ctx context.Context, chal *acme.Challenge) (cp1 *acme.Challenge, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.Accept method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.Accept method")
	}

	m.mock.funcAccept = f
	return m.mock
}

// When sets expectation for the AcmeClient.Accept which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockAccept) When(ctx context.Context, chal *acme.Challenge) *AcmeClientMockAcceptExpectation {
	if m.mock.funcAccept != nil {
		m.mock.t.Fatalf("AcmeClientMock.Accept mock is already set by Set")
	}

	expectation := &AcmeClientMockAcceptExpectation{
		mock:   m.mock,
		params: &AcmeClientMockAcceptParams{ctx, chal},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.Accept return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockAcceptExpectation) Then(cp1 *acme.Challenge, err error) *AcmeClientMock {
	e.results = &AcmeClientMockAcceptResults{cp1, err}
	return e.mock
}

// Accept implements AcmeClient
func (m *AcmeClientMock) Accept(ctx context.Context, chal *acme.Challenge) (cp1 *acme.Challenge, err error) {
	atomic.AddUint64(&m.beforeAcceptCounter, 1)
	defer atomic.AddUint64(&m.afterAcceptCounter, 1)

	for _, e := range m.AcceptMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockAcceptParams{ctx, chal}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.cp1, e.results.err
		}
	}

	if m.AcceptMock.defaultExpectation != nil {
		atomic.AddUint64(&m.AcceptMock.defaultExpectation.Counter, 1)
		want := m.AcceptMock.defaultExpectation.params
		got := AcmeClientMockAcceptParams{ctx, chal}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.Accept got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.AcceptMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.Accept")
		}
		return (*results).cp1, (*results).err
	}
	if m.funcAccept != nil {
		return m.funcAccept(ctx, chal)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.Accept. %v %v", ctx, chal)
	return
}

// AcceptAfterCounter returns a count of finished AcmeClientMock.Accept invocations
func (m *AcmeClientMock) AcceptAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterAcceptCounter)
}

// AcceptBeforeCounter returns a count of AcmeClientMock.Accept invocations
func (m *AcmeClientMock) AcceptBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeAcceptCounter)
}

// MinimockAcceptDone returns true if the count of the Accept invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockAcceptDone() bool {
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
func (m *AcmeClientMock) MinimockAcceptInspect() {
	for _, e := range m.AcceptMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.Accept with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AcceptMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAcceptCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.Accept with params: %#v", *m.AcceptMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAccept != nil && atomic.LoadUint64(&m.afterAcceptCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.Accept")
	}
}

type mAcmeClientMockAuthorizeOrder struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockAuthorizeOrderExpectation
	expectations       []*AcmeClientMockAuthorizeOrderExpectation
}

// AcmeClientMockAuthorizeOrderExpectation specifies expectation struct of the AcmeClient.AuthorizeOrder
type AcmeClientMockAuthorizeOrderExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockAuthorizeOrderParams
	results *AcmeClientMockAuthorizeOrderResults
	Counter uint64
}

// AcmeClientMockAuthorizeOrderParams contains parameters of the AcmeClient.AuthorizeOrder
type AcmeClientMockAuthorizeOrderParams struct {
	ctx context.Context
	id  []acme.AuthzID
	opt []acme.OrderOption
}

// AcmeClientMockAuthorizeOrderResults contains results of the AcmeClient.AuthorizeOrder
type AcmeClientMockAuthorizeOrderResults struct {
	op1 *acme.Order
	err error
}

// Expect sets up expected params for AcmeClient.AuthorizeOrder
func (m *mAcmeClientMockAuthorizeOrder) Expect(ctx context.Context, id []acme.AuthzID, opt ...acme.OrderOption) *mAcmeClientMockAuthorizeOrder {
	if m.mock.funcAuthorizeOrder != nil {
		m.mock.t.Fatalf("AcmeClientMock.AuthorizeOrder mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockAuthorizeOrderExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockAuthorizeOrderParams{ctx, id, opt}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.AuthorizeOrder
func (m *mAcmeClientMockAuthorizeOrder) Return(op1 *acme.Order, err error) *AcmeClientMock {
	if m.mock.funcAuthorizeOrder != nil {
		m.mock.t.Fatalf("AcmeClientMock.AuthorizeOrder mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockAuthorizeOrderExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockAuthorizeOrderResults{op1, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.AuthorizeOrder method
func (m *mAcmeClientMockAuthorizeOrder) Set(f func(ctx context.Context, id []acme.AuthzID, opt ...acme.OrderOption) (op1 *acme.Order, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.AuthorizeOrder method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.AuthorizeOrder method")
	}

	m.mock.funcAuthorizeOrder = f
	return m.mock
}

// When sets expectation for the AcmeClient.AuthorizeOrder which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockAuthorizeOrder) When(ctx context.Context, id []acme.AuthzID, opt ...acme.OrderOption) *AcmeClientMockAuthorizeOrderExpectation {
	if m.mock.funcAuthorizeOrder != nil {
		m.mock.t.Fatalf("AcmeClientMock.AuthorizeOrder mock is already set by Set")
	}

	expectation := &AcmeClientMockAuthorizeOrderExpectation{
		mock:   m.mock,
		params: &AcmeClientMockAuthorizeOrderParams{ctx, id, opt},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.AuthorizeOrder return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockAuthorizeOrderExpectation) Then(op1 *acme.Order, err error) *AcmeClientMock {
	e.results = &AcmeClientMockAuthorizeOrderResults{op1, err}
	return e.mock
}

// AuthorizeOrder implements AcmeClient
func (m *AcmeClientMock) AuthorizeOrder(ctx context.Context, id []acme.AuthzID, opt ...acme.OrderOption) (op1 *acme.Order, err error) {
	atomic.AddUint64(&m.beforeAuthorizeOrderCounter, 1)
	defer atomic.AddUint64(&m.afterAuthorizeOrderCounter, 1)

	for _, e := range m.AuthorizeOrderMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockAuthorizeOrderParams{ctx, id, opt}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.op1, e.results.err
		}
	}

	if m.AuthorizeOrderMock.defaultExpectation != nil {
		atomic.AddUint64(&m.AuthorizeOrderMock.defaultExpectation.Counter, 1)
		want := m.AuthorizeOrderMock.defaultExpectation.params
		got := AcmeClientMockAuthorizeOrderParams{ctx, id, opt}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.AuthorizeOrder got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.AuthorizeOrderMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.AuthorizeOrder")
		}
		return (*results).op1, (*results).err
	}
	if m.funcAuthorizeOrder != nil {
		return m.funcAuthorizeOrder(ctx, id, opt...)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.AuthorizeOrder. %v %v %v", ctx, id, opt)
	return
}

// AuthorizeOrderAfterCounter returns a count of finished AcmeClientMock.AuthorizeOrder invocations
func (m *AcmeClientMock) AuthorizeOrderAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterAuthorizeOrderCounter)
}

// AuthorizeOrderBeforeCounter returns a count of AcmeClientMock.AuthorizeOrder invocations
func (m *AcmeClientMock) AuthorizeOrderBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeAuthorizeOrderCounter)
}

// MinimockAuthorizeOrderDone returns true if the count of the AuthorizeOrder invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockAuthorizeOrderDone() bool {
	for _, e := range m.AuthorizeOrderMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AuthorizeOrderMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAuthorizeOrderCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAuthorizeOrder != nil && atomic.LoadUint64(&m.afterAuthorizeOrderCounter) < 1 {
		return false
	}
	return true
}

// MinimockAuthorizeOrderInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockAuthorizeOrderInspect() {
	for _, e := range m.AuthorizeOrderMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.AuthorizeOrder with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AuthorizeOrderMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAuthorizeOrderCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.AuthorizeOrder with params: %#v", *m.AuthorizeOrderMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAuthorizeOrder != nil && atomic.LoadUint64(&m.afterAuthorizeOrderCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.AuthorizeOrder")
	}
}

type mAcmeClientMockCreateOrderCert struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockCreateOrderCertExpectation
	expectations       []*AcmeClientMockCreateOrderCertExpectation
}

// AcmeClientMockCreateOrderCertExpectation specifies expectation struct of the AcmeClient.CreateOrderCert
type AcmeClientMockCreateOrderCertExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockCreateOrderCertParams
	results *AcmeClientMockCreateOrderCertResults
	Counter uint64
}

// AcmeClientMockCreateOrderCertParams contains parameters of the AcmeClient.CreateOrderCert
type AcmeClientMockCreateOrderCertParams struct {
	ctx    context.Context
	url    string
	csr    []byte
	bundle bool
}

// AcmeClientMockCreateOrderCertResults contains results of the AcmeClient.CreateOrderCert
type AcmeClientMockCreateOrderCertResults struct {
	der     [][]byte
	certURL string
	err     error
}

// Expect sets up expected params for AcmeClient.CreateOrderCert
func (m *mAcmeClientMockCreateOrderCert) Expect(ctx context.Context, url string, csr []byte, bundle bool) *mAcmeClientMockCreateOrderCert {
	if m.mock.funcCreateOrderCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.CreateOrderCert mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockCreateOrderCertExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockCreateOrderCertParams{ctx, url, csr, bundle}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.CreateOrderCert
func (m *mAcmeClientMockCreateOrderCert) Return(der [][]byte, certURL string, err error) *AcmeClientMock {
	if m.mock.funcCreateOrderCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.CreateOrderCert mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockCreateOrderCertExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockCreateOrderCertResults{der, certURL, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.CreateOrderCert method
func (m *mAcmeClientMockCreateOrderCert) Set(f func(ctx context.Context, url string, csr []byte, bundle bool) (der [][]byte, certURL string, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.CreateOrderCert method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.CreateOrderCert method")
	}

	m.mock.funcCreateOrderCert = f
	return m.mock
}

// When sets expectation for the AcmeClient.CreateOrderCert which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockCreateOrderCert) When(ctx context.Context, url string, csr []byte, bundle bool) *AcmeClientMockCreateOrderCertExpectation {
	if m.mock.funcCreateOrderCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.CreateOrderCert mock is already set by Set")
	}

	expectation := &AcmeClientMockCreateOrderCertExpectation{
		mock:   m.mock,
		params: &AcmeClientMockCreateOrderCertParams{ctx, url, csr, bundle},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.CreateOrderCert return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockCreateOrderCertExpectation) Then(der [][]byte, certURL string, err error) *AcmeClientMock {
	e.results = &AcmeClientMockCreateOrderCertResults{der, certURL, err}
	return e.mock
}

// CreateOrderCert implements AcmeClient
func (m *AcmeClientMock) CreateOrderCert(ctx context.Context, url string, csr []byte, bundle bool) (der [][]byte, certURL string, err error) {
	atomic.AddUint64(&m.beforeCreateOrderCertCounter, 1)
	defer atomic.AddUint64(&m.afterCreateOrderCertCounter, 1)

	for _, e := range m.CreateOrderCertMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockCreateOrderCertParams{ctx, url, csr, bundle}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.der, e.results.certURL, e.results.err
		}
	}

	if m.CreateOrderCertMock.defaultExpectation != nil {
		atomic.AddUint64(&m.CreateOrderCertMock.defaultExpectation.Counter, 1)
		want := m.CreateOrderCertMock.defaultExpectation.params
		got := AcmeClientMockCreateOrderCertParams{ctx, url, csr, bundle}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.CreateOrderCert got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.CreateOrderCertMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.CreateOrderCert")
		}
		return (*results).der, (*results).certURL, (*results).err
	}
	if m.funcCreateOrderCert != nil {
		return m.funcCreateOrderCert(ctx, url, csr, bundle)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.CreateOrderCert. %v %v %v %v", ctx, url, csr, bundle)
	return
}

// CreateOrderCertAfterCounter returns a count of finished AcmeClientMock.CreateOrderCert invocations
func (m *AcmeClientMock) CreateOrderCertAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterCreateOrderCertCounter)
}

// CreateOrderCertBeforeCounter returns a count of AcmeClientMock.CreateOrderCert invocations
func (m *AcmeClientMock) CreateOrderCertBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeCreateOrderCertCounter)
}

// MinimockCreateOrderCertDone returns true if the count of the CreateOrderCert invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockCreateOrderCertDone() bool {
	for _, e := range m.CreateOrderCertMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CreateOrderCertMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterCreateOrderCertCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreateOrderCert != nil && atomic.LoadUint64(&m.afterCreateOrderCertCounter) < 1 {
		return false
	}
	return true
}

// MinimockCreateOrderCertInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockCreateOrderCertInspect() {
	for _, e := range m.CreateOrderCertMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.CreateOrderCert with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CreateOrderCertMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterCreateOrderCertCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.CreateOrderCert with params: %#v", *m.CreateOrderCertMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreateOrderCert != nil && atomic.LoadUint64(&m.afterCreateOrderCertCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.CreateOrderCert")
	}
}

type mAcmeClientMockGetAuthorization struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockGetAuthorizationExpectation
	expectations       []*AcmeClientMockGetAuthorizationExpectation
}

// AcmeClientMockGetAuthorizationExpectation specifies expectation struct of the AcmeClient.GetAuthorization
type AcmeClientMockGetAuthorizationExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockGetAuthorizationParams
	results *AcmeClientMockGetAuthorizationResults
	Counter uint64
}

// AcmeClientMockGetAuthorizationParams contains parameters of the AcmeClient.GetAuthorization
type AcmeClientMockGetAuthorizationParams struct {
	ctx context.Context
	url string
}

// AcmeClientMockGetAuthorizationResults contains results of the AcmeClient.GetAuthorization
type AcmeClientMockGetAuthorizationResults struct {
	ap1 *acme.Authorization
	err error
}

// Expect sets up expected params for AcmeClient.GetAuthorization
func (m *mAcmeClientMockGetAuthorization) Expect(ctx context.Context, url string) *mAcmeClientMockGetAuthorization {
	if m.mock.funcGetAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.GetAuthorization mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockGetAuthorizationExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockGetAuthorizationParams{ctx, url}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.GetAuthorization
func (m *mAcmeClientMockGetAuthorization) Return(ap1 *acme.Authorization, err error) *AcmeClientMock {
	if m.mock.funcGetAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.GetAuthorization mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockGetAuthorizationExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockGetAuthorizationResults{ap1, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.GetAuthorization method
func (m *mAcmeClientMockGetAuthorization) Set(f func(ctx context.Context, url string) (ap1 *acme.Authorization, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.GetAuthorization method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.GetAuthorization method")
	}

	m.mock.funcGetAuthorization = f
	return m.mock
}

// When sets expectation for the AcmeClient.GetAuthorization which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockGetAuthorization) When(ctx context.Context, url string) *AcmeClientMockGetAuthorizationExpectation {
	if m.mock.funcGetAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.GetAuthorization mock is already set by Set")
	}

	expectation := &AcmeClientMockGetAuthorizationExpectation{
		mock:   m.mock,
		params: &AcmeClientMockGetAuthorizationParams{ctx, url},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.GetAuthorization return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockGetAuthorizationExpectation) Then(ap1 *acme.Authorization, err error) *AcmeClientMock {
	e.results = &AcmeClientMockGetAuthorizationResults{ap1, err}
	return e.mock
}

// GetAuthorization implements AcmeClient
func (m *AcmeClientMock) GetAuthorization(ctx context.Context, url string) (ap1 *acme.Authorization, err error) {
	atomic.AddUint64(&m.beforeGetAuthorizationCounter, 1)
	defer atomic.AddUint64(&m.afterGetAuthorizationCounter, 1)

	for _, e := range m.GetAuthorizationMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockGetAuthorizationParams{ctx, url}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.ap1, e.results.err
		}
	}

	if m.GetAuthorizationMock.defaultExpectation != nil {
		atomic.AddUint64(&m.GetAuthorizationMock.defaultExpectation.Counter, 1)
		want := m.GetAuthorizationMock.defaultExpectation.params
		got := AcmeClientMockGetAuthorizationParams{ctx, url}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.GetAuthorization got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetAuthorizationMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.GetAuthorization")
		}
		return (*results).ap1, (*results).err
	}
	if m.funcGetAuthorization != nil {
		return m.funcGetAuthorization(ctx, url)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.GetAuthorization. %v %v", ctx, url)
	return
}

// GetAuthorizationAfterCounter returns a count of finished AcmeClientMock.GetAuthorization invocations
func (m *AcmeClientMock) GetAuthorizationAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterGetAuthorizationCounter)
}

// GetAuthorizationBeforeCounter returns a count of AcmeClientMock.GetAuthorization invocations
func (m *AcmeClientMock) GetAuthorizationBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeGetAuthorizationCounter)
}

// MinimockGetAuthorizationDone returns true if the count of the GetAuthorization invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockGetAuthorizationDone() bool {
	for _, e := range m.GetAuthorizationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetAuthorizationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetAuthorizationCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetAuthorization != nil && atomic.LoadUint64(&m.afterGetAuthorizationCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetAuthorizationInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockGetAuthorizationInspect() {
	for _, e := range m.GetAuthorizationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.GetAuthorization with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetAuthorizationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterGetAuthorizationCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.GetAuthorization with params: %#v", *m.GetAuthorizationMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetAuthorization != nil && atomic.LoadUint64(&m.afterGetAuthorizationCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.GetAuthorization")
	}
}

type mAcmeClientMockHTTP01ChallengeResponse struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockHTTP01ChallengeResponseExpectation
	expectations       []*AcmeClientMockHTTP01ChallengeResponseExpectation
}

// AcmeClientMockHTTP01ChallengeResponseExpectation specifies expectation struct of the AcmeClient.HTTP01ChallengeResponse
type AcmeClientMockHTTP01ChallengeResponseExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockHTTP01ChallengeResponseParams
	results *AcmeClientMockHTTP01ChallengeResponseResults
	Counter uint64
}

// AcmeClientMockHTTP01ChallengeResponseParams contains parameters of the AcmeClient.HTTP01ChallengeResponse
type AcmeClientMockHTTP01ChallengeResponseParams struct {
	token string
}

// AcmeClientMockHTTP01ChallengeResponseResults contains results of the AcmeClient.HTTP01ChallengeResponse
type AcmeClientMockHTTP01ChallengeResponseResults struct {
	s1  string
	err error
}

// Expect sets up expected params for AcmeClient.HTTP01ChallengeResponse
func (m *mAcmeClientMockHTTP01ChallengeResponse) Expect(token string) *mAcmeClientMockHTTP01ChallengeResponse {
	if m.mock.funcHTTP01ChallengeResponse != nil {
		m.mock.t.Fatalf("AcmeClientMock.HTTP01ChallengeResponse mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockHTTP01ChallengeResponseExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockHTTP01ChallengeResponseParams{token}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.HTTP01ChallengeResponse
func (m *mAcmeClientMockHTTP01ChallengeResponse) Return(s1 string, err error) *AcmeClientMock {
	if m.mock.funcHTTP01ChallengeResponse != nil {
		m.mock.t.Fatalf("AcmeClientMock.HTTP01ChallengeResponse mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockHTTP01ChallengeResponseExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockHTTP01ChallengeResponseResults{s1, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.HTTP01ChallengeResponse method
func (m *mAcmeClientMockHTTP01ChallengeResponse) Set(f func(token string) (s1 string, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.HTTP01ChallengeResponse method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.HTTP01ChallengeResponse method")
	}

	m.mock.funcHTTP01ChallengeResponse = f
	return m.mock
}

// When sets expectation for the AcmeClient.HTTP01ChallengeResponse which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockHTTP01ChallengeResponse) When(token string) *AcmeClientMockHTTP01ChallengeResponseExpectation {
	if m.mock.funcHTTP01ChallengeResponse != nil {
		m.mock.t.Fatalf("AcmeClientMock.HTTP01ChallengeResponse mock is already set by Set")
	}

	expectation := &AcmeClientMockHTTP01ChallengeResponseExpectation{
		mock:   m.mock,
		params: &AcmeClientMockHTTP01ChallengeResponseParams{token},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.HTTP01ChallengeResponse return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockHTTP01ChallengeResponseExpectation) Then(s1 string, err error) *AcmeClientMock {
	e.results = &AcmeClientMockHTTP01ChallengeResponseResults{s1, err}
	return e.mock
}

// HTTP01ChallengeResponse implements AcmeClient
func (m *AcmeClientMock) HTTP01ChallengeResponse(token string) (s1 string, err error) {
	atomic.AddUint64(&m.beforeHTTP01ChallengeResponseCounter, 1)
	defer atomic.AddUint64(&m.afterHTTP01ChallengeResponseCounter, 1)

	for _, e := range m.HTTP01ChallengeResponseMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockHTTP01ChallengeResponseParams{token}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.s1, e.results.err
		}
	}

	if m.HTTP01ChallengeResponseMock.defaultExpectation != nil {
		atomic.AddUint64(&m.HTTP01ChallengeResponseMock.defaultExpectation.Counter, 1)
		want := m.HTTP01ChallengeResponseMock.defaultExpectation.params
		got := AcmeClientMockHTTP01ChallengeResponseParams{token}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.HTTP01ChallengeResponse got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.HTTP01ChallengeResponseMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.HTTP01ChallengeResponse")
		}
		return (*results).s1, (*results).err
	}
	if m.funcHTTP01ChallengeResponse != nil {
		return m.funcHTTP01ChallengeResponse(token)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.HTTP01ChallengeResponse. %v", token)
	return
}

// HTTP01ChallengeResponseAfterCounter returns a count of finished AcmeClientMock.HTTP01ChallengeResponse invocations
func (m *AcmeClientMock) HTTP01ChallengeResponseAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterHTTP01ChallengeResponseCounter)
}

// HTTP01ChallengeResponseBeforeCounter returns a count of AcmeClientMock.HTTP01ChallengeResponse invocations
func (m *AcmeClientMock) HTTP01ChallengeResponseBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeHTTP01ChallengeResponseCounter)
}

// MinimockHTTP01ChallengeResponseDone returns true if the count of the HTTP01ChallengeResponse invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockHTTP01ChallengeResponseDone() bool {
	for _, e := range m.HTTP01ChallengeResponseMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HTTP01ChallengeResponseMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterHTTP01ChallengeResponseCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHTTP01ChallengeResponse != nil && atomic.LoadUint64(&m.afterHTTP01ChallengeResponseCounter) < 1 {
		return false
	}
	return true
}

// MinimockHTTP01ChallengeResponseInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockHTTP01ChallengeResponseInspect() {
	for _, e := range m.HTTP01ChallengeResponseMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.HTTP01ChallengeResponse with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.HTTP01ChallengeResponseMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterHTTP01ChallengeResponseCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.HTTP01ChallengeResponse with params: %#v", *m.HTTP01ChallengeResponseMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcHTTP01ChallengeResponse != nil && atomic.LoadUint64(&m.afterHTTP01ChallengeResponseCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.HTTP01ChallengeResponse")
	}
}

type mAcmeClientMockRevokeAuthorization struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockRevokeAuthorizationExpectation
	expectations       []*AcmeClientMockRevokeAuthorizationExpectation
}

// AcmeClientMockRevokeAuthorizationExpectation specifies expectation struct of the AcmeClient.RevokeAuthorization
type AcmeClientMockRevokeAuthorizationExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockRevokeAuthorizationParams
	results *AcmeClientMockRevokeAuthorizationResults
	Counter uint64
}

// AcmeClientMockRevokeAuthorizationParams contains parameters of the AcmeClient.RevokeAuthorization
type AcmeClientMockRevokeAuthorizationParams struct {
	ctx context.Context
	url string
}

// AcmeClientMockRevokeAuthorizationResults contains results of the AcmeClient.RevokeAuthorization
type AcmeClientMockRevokeAuthorizationResults struct {
	err error
}

// Expect sets up expected params for AcmeClient.RevokeAuthorization
func (m *mAcmeClientMockRevokeAuthorization) Expect(ctx context.Context, url string) *mAcmeClientMockRevokeAuthorization {
	if m.mock.funcRevokeAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.RevokeAuthorization mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockRevokeAuthorizationExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockRevokeAuthorizationParams{ctx, url}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.RevokeAuthorization
func (m *mAcmeClientMockRevokeAuthorization) Return(err error) *AcmeClientMock {
	if m.mock.funcRevokeAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.RevokeAuthorization mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockRevokeAuthorizationExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockRevokeAuthorizationResults{err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.RevokeAuthorization method
func (m *mAcmeClientMockRevokeAuthorization) Set(f func(ctx context.Context, url string) (err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.RevokeAuthorization method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.RevokeAuthorization method")
	}

	m.mock.funcRevokeAuthorization = f
	return m.mock
}

// When sets expectation for the AcmeClient.RevokeAuthorization which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockRevokeAuthorization) When(ctx context.Context, url string) *AcmeClientMockRevokeAuthorizationExpectation {
	if m.mock.funcRevokeAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.RevokeAuthorization mock is already set by Set")
	}

	expectation := &AcmeClientMockRevokeAuthorizationExpectation{
		mock:   m.mock,
		params: &AcmeClientMockRevokeAuthorizationParams{ctx, url},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.RevokeAuthorization return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockRevokeAuthorizationExpectation) Then(err error) *AcmeClientMock {
	e.results = &AcmeClientMockRevokeAuthorizationResults{err}
	return e.mock
}

// RevokeAuthorization implements AcmeClient
func (m *AcmeClientMock) RevokeAuthorization(ctx context.Context, url string) (err error) {
	atomic.AddUint64(&m.beforeRevokeAuthorizationCounter, 1)
	defer atomic.AddUint64(&m.afterRevokeAuthorizationCounter, 1)

	for _, e := range m.RevokeAuthorizationMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockRevokeAuthorizationParams{ctx, url}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.RevokeAuthorizationMock.defaultExpectation != nil {
		atomic.AddUint64(&m.RevokeAuthorizationMock.defaultExpectation.Counter, 1)
		want := m.RevokeAuthorizationMock.defaultExpectation.params
		got := AcmeClientMockRevokeAuthorizationParams{ctx, url}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.RevokeAuthorization got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.RevokeAuthorizationMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.RevokeAuthorization")
		}
		return (*results).err
	}
	if m.funcRevokeAuthorization != nil {
		return m.funcRevokeAuthorization(ctx, url)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.RevokeAuthorization. %v %v", ctx, url)
	return
}

// RevokeAuthorizationAfterCounter returns a count of finished AcmeClientMock.RevokeAuthorization invocations
func (m *AcmeClientMock) RevokeAuthorizationAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterRevokeAuthorizationCounter)
}

// RevokeAuthorizationBeforeCounter returns a count of AcmeClientMock.RevokeAuthorization invocations
func (m *AcmeClientMock) RevokeAuthorizationBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeRevokeAuthorizationCounter)
}

// MinimockRevokeAuthorizationDone returns true if the count of the RevokeAuthorization invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockRevokeAuthorizationDone() bool {
	for _, e := range m.RevokeAuthorizationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RevokeAuthorizationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRevokeAuthorizationCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRevokeAuthorization != nil && atomic.LoadUint64(&m.afterRevokeAuthorizationCounter) < 1 {
		return false
	}
	return true
}

// MinimockRevokeAuthorizationInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockRevokeAuthorizationInspect() {
	for _, e := range m.RevokeAuthorizationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.RevokeAuthorization with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RevokeAuthorizationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRevokeAuthorizationCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.RevokeAuthorization with params: %#v", *m.RevokeAuthorizationMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRevokeAuthorization != nil && atomic.LoadUint64(&m.afterRevokeAuthorizationCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.RevokeAuthorization")
	}
}

type mAcmeClientMockTLSALPN01ChallengeCert struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockTLSALPN01ChallengeCertExpectation
	expectations       []*AcmeClientMockTLSALPN01ChallengeCertExpectation
}

// AcmeClientMockTLSALPN01ChallengeCertExpectation specifies expectation struct of the AcmeClient.TLSALPN01ChallengeCert
type AcmeClientMockTLSALPN01ChallengeCertExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockTLSALPN01ChallengeCertParams
	results *AcmeClientMockTLSALPN01ChallengeCertResults
	Counter uint64
}

// AcmeClientMockTLSALPN01ChallengeCertParams contains parameters of the AcmeClient.TLSALPN01ChallengeCert
type AcmeClientMockTLSALPN01ChallengeCertParams struct {
	token  string
	domain string
	opt    []acme.CertOption
}

// AcmeClientMockTLSALPN01ChallengeCertResults contains results of the AcmeClient.TLSALPN01ChallengeCert
type AcmeClientMockTLSALPN01ChallengeCertResults struct {
	cert tls.Certificate
	err  error
}

// Expect sets up expected params for AcmeClient.TLSALPN01ChallengeCert
func (m *mAcmeClientMockTLSALPN01ChallengeCert) Expect(token string, domain string, opt ...acme.CertOption) *mAcmeClientMockTLSALPN01ChallengeCert {
	if m.mock.funcTLSALPN01ChallengeCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.TLSALPN01ChallengeCert mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockTLSALPN01ChallengeCertExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockTLSALPN01ChallengeCertParams{token, domain, opt}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.TLSALPN01ChallengeCert
func (m *mAcmeClientMockTLSALPN01ChallengeCert) Return(cert tls.Certificate, err error) *AcmeClientMock {
	if m.mock.funcTLSALPN01ChallengeCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.TLSALPN01ChallengeCert mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockTLSALPN01ChallengeCertExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockTLSALPN01ChallengeCertResults{cert, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.TLSALPN01ChallengeCert method
func (m *mAcmeClientMockTLSALPN01ChallengeCert) Set(f func(token string, domain string, opt ...acme.CertOption) (cert tls.Certificate, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.TLSALPN01ChallengeCert method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.TLSALPN01ChallengeCert method")
	}

	m.mock.funcTLSALPN01ChallengeCert = f
	return m.mock
}

// When sets expectation for the AcmeClient.TLSALPN01ChallengeCert which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockTLSALPN01ChallengeCert) When(token string, domain string, opt ...acme.CertOption) *AcmeClientMockTLSALPN01ChallengeCertExpectation {
	if m.mock.funcTLSALPN01ChallengeCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.TLSALPN01ChallengeCert mock is already set by Set")
	}

	expectation := &AcmeClientMockTLSALPN01ChallengeCertExpectation{
		mock:   m.mock,
		params: &AcmeClientMockTLSALPN01ChallengeCertParams{token, domain, opt},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.TLSALPN01ChallengeCert return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockTLSALPN01ChallengeCertExpectation) Then(cert tls.Certificate, err error) *AcmeClientMock {
	e.results = &AcmeClientMockTLSALPN01ChallengeCertResults{cert, err}
	return e.mock
}

// TLSALPN01ChallengeCert implements AcmeClient
func (m *AcmeClientMock) TLSALPN01ChallengeCert(token string, domain string, opt ...acme.CertOption) (cert tls.Certificate, err error) {
	atomic.AddUint64(&m.beforeTLSALPN01ChallengeCertCounter, 1)
	defer atomic.AddUint64(&m.afterTLSALPN01ChallengeCertCounter, 1)

	for _, e := range m.TLSALPN01ChallengeCertMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockTLSALPN01ChallengeCertParams{token, domain, opt}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.cert, e.results.err
		}
	}

	if m.TLSALPN01ChallengeCertMock.defaultExpectation != nil {
		atomic.AddUint64(&m.TLSALPN01ChallengeCertMock.defaultExpectation.Counter, 1)
		want := m.TLSALPN01ChallengeCertMock.defaultExpectation.params
		got := AcmeClientMockTLSALPN01ChallengeCertParams{token, domain, opt}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.TLSALPN01ChallengeCert got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.TLSALPN01ChallengeCertMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.TLSALPN01ChallengeCert")
		}
		return (*results).cert, (*results).err
	}
	if m.funcTLSALPN01ChallengeCert != nil {
		return m.funcTLSALPN01ChallengeCert(token, domain, opt...)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.TLSALPN01ChallengeCert. %v %v %v", token, domain, opt)
	return
}

// TLSALPN01ChallengeCertAfterCounter returns a count of finished AcmeClientMock.TLSALPN01ChallengeCert invocations
func (m *AcmeClientMock) TLSALPN01ChallengeCertAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterTLSALPN01ChallengeCertCounter)
}

// TLSALPN01ChallengeCertBeforeCounter returns a count of AcmeClientMock.TLSALPN01ChallengeCert invocations
func (m *AcmeClientMock) TLSALPN01ChallengeCertBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeTLSALPN01ChallengeCertCounter)
}

// MinimockTLSALPN01ChallengeCertDone returns true if the count of the TLSALPN01ChallengeCert invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockTLSALPN01ChallengeCertDone() bool {
	for _, e := range m.TLSALPN01ChallengeCertMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.TLSALPN01ChallengeCertMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterTLSALPN01ChallengeCertCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcTLSALPN01ChallengeCert != nil && atomic.LoadUint64(&m.afterTLSALPN01ChallengeCertCounter) < 1 {
		return false
	}
	return true
}

// MinimockTLSALPN01ChallengeCertInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockTLSALPN01ChallengeCertInspect() {
	for _, e := range m.TLSALPN01ChallengeCertMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.TLSALPN01ChallengeCert with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.TLSALPN01ChallengeCertMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterTLSALPN01ChallengeCertCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.TLSALPN01ChallengeCert with params: %#v", *m.TLSALPN01ChallengeCertMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcTLSALPN01ChallengeCert != nil && atomic.LoadUint64(&m.afterTLSALPN01ChallengeCertCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.TLSALPN01ChallengeCert")
	}
}

type mAcmeClientMockWaitAuthorization struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockWaitAuthorizationExpectation
	expectations       []*AcmeClientMockWaitAuthorizationExpectation
}

// AcmeClientMockWaitAuthorizationExpectation specifies expectation struct of the AcmeClient.WaitAuthorization
type AcmeClientMockWaitAuthorizationExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockWaitAuthorizationParams
	results *AcmeClientMockWaitAuthorizationResults
	Counter uint64
}

// AcmeClientMockWaitAuthorizationParams contains parameters of the AcmeClient.WaitAuthorization
type AcmeClientMockWaitAuthorizationParams struct {
	ctx context.Context
	url string
}

// AcmeClientMockWaitAuthorizationResults contains results of the AcmeClient.WaitAuthorization
type AcmeClientMockWaitAuthorizationResults struct {
	ap1 *acme.Authorization
	err error
}

// Expect sets up expected params for AcmeClient.WaitAuthorization
func (m *mAcmeClientMockWaitAuthorization) Expect(ctx context.Context, url string) *mAcmeClientMockWaitAuthorization {
	if m.mock.funcWaitAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.WaitAuthorization mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockWaitAuthorizationExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockWaitAuthorizationParams{ctx, url}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.WaitAuthorization
func (m *mAcmeClientMockWaitAuthorization) Return(ap1 *acme.Authorization, err error) *AcmeClientMock {
	if m.mock.funcWaitAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.WaitAuthorization mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockWaitAuthorizationExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockWaitAuthorizationResults{ap1, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.WaitAuthorization method
func (m *mAcmeClientMockWaitAuthorization) Set(f func(ctx context.Context, url string) (ap1 *acme.Authorization, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.WaitAuthorization method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.WaitAuthorization method")
	}

	m.mock.funcWaitAuthorization = f
	return m.mock
}

// When sets expectation for the AcmeClient.WaitAuthorization which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockWaitAuthorization) When(ctx context.Context, url string) *AcmeClientMockWaitAuthorizationExpectation {
	if m.mock.funcWaitAuthorization != nil {
		m.mock.t.Fatalf("AcmeClientMock.WaitAuthorization mock is already set by Set")
	}

	expectation := &AcmeClientMockWaitAuthorizationExpectation{
		mock:   m.mock,
		params: &AcmeClientMockWaitAuthorizationParams{ctx, url},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.WaitAuthorization return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockWaitAuthorizationExpectation) Then(ap1 *acme.Authorization, err error) *AcmeClientMock {
	e.results = &AcmeClientMockWaitAuthorizationResults{ap1, err}
	return e.mock
}

// WaitAuthorization implements AcmeClient
func (m *AcmeClientMock) WaitAuthorization(ctx context.Context, url string) (ap1 *acme.Authorization, err error) {
	atomic.AddUint64(&m.beforeWaitAuthorizationCounter, 1)
	defer atomic.AddUint64(&m.afterWaitAuthorizationCounter, 1)

	for _, e := range m.WaitAuthorizationMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockWaitAuthorizationParams{ctx, url}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.ap1, e.results.err
		}
	}

	if m.WaitAuthorizationMock.defaultExpectation != nil {
		atomic.AddUint64(&m.WaitAuthorizationMock.defaultExpectation.Counter, 1)
		want := m.WaitAuthorizationMock.defaultExpectation.params
		got := AcmeClientMockWaitAuthorizationParams{ctx, url}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.WaitAuthorization got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.WaitAuthorizationMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.WaitAuthorization")
		}
		return (*results).ap1, (*results).err
	}
	if m.funcWaitAuthorization != nil {
		return m.funcWaitAuthorization(ctx, url)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.WaitAuthorization. %v %v", ctx, url)
	return
}

// WaitAuthorizationAfterCounter returns a count of finished AcmeClientMock.WaitAuthorization invocations
func (m *AcmeClientMock) WaitAuthorizationAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterWaitAuthorizationCounter)
}

// WaitAuthorizationBeforeCounter returns a count of AcmeClientMock.WaitAuthorization invocations
func (m *AcmeClientMock) WaitAuthorizationBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeWaitAuthorizationCounter)
}

// MinimockWaitAuthorizationDone returns true if the count of the WaitAuthorization invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockWaitAuthorizationDone() bool {
	for _, e := range m.WaitAuthorizationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WaitAuthorizationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWaitAuthorizationCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWaitAuthorization != nil && atomic.LoadUint64(&m.afterWaitAuthorizationCounter) < 1 {
		return false
	}
	return true
}

// MinimockWaitAuthorizationInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockWaitAuthorizationInspect() {
	for _, e := range m.WaitAuthorizationMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.WaitAuthorization with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WaitAuthorizationMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWaitAuthorizationCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.WaitAuthorization with params: %#v", *m.WaitAuthorizationMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWaitAuthorization != nil && atomic.LoadUint64(&m.afterWaitAuthorizationCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.WaitAuthorization")
	}
}

type mAcmeClientMockWaitOrder struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockWaitOrderExpectation
	expectations       []*AcmeClientMockWaitOrderExpectation
}

// AcmeClientMockWaitOrderExpectation specifies expectation struct of the AcmeClient.WaitOrder
type AcmeClientMockWaitOrderExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockWaitOrderParams
	results *AcmeClientMockWaitOrderResults
	Counter uint64
}

// AcmeClientMockWaitOrderParams contains parameters of the AcmeClient.WaitOrder
type AcmeClientMockWaitOrderParams struct {
	ctx context.Context
	url string
}

// AcmeClientMockWaitOrderResults contains results of the AcmeClient.WaitOrder
type AcmeClientMockWaitOrderResults struct {
	op1 *acme.Order
	err error
}

// Expect sets up expected params for AcmeClient.WaitOrder
func (m *mAcmeClientMockWaitOrder) Expect(ctx context.Context, url string) *mAcmeClientMockWaitOrder {
	if m.mock.funcWaitOrder != nil {
		m.mock.t.Fatalf("AcmeClientMock.WaitOrder mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockWaitOrderExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockWaitOrderParams{ctx, url}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.WaitOrder
func (m *mAcmeClientMockWaitOrder) Return(op1 *acme.Order, err error) *AcmeClientMock {
	if m.mock.funcWaitOrder != nil {
		m.mock.t.Fatalf("AcmeClientMock.WaitOrder mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockWaitOrderExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockWaitOrderResults{op1, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.WaitOrder method
func (m *mAcmeClientMockWaitOrder) Set(f func(ctx context.Context, url string) (op1 *acme.Order, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.WaitOrder method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.WaitOrder method")
	}

	m.mock.funcWaitOrder = f
	return m.mock
}

// When sets expectation for the AcmeClient.WaitOrder which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockWaitOrder) When(ctx context.Context, url string) *AcmeClientMockWaitOrderExpectation {
	if m.mock.funcWaitOrder != nil {
		m.mock.t.Fatalf("AcmeClientMock.WaitOrder mock is already set by Set")
	}

	expectation := &AcmeClientMockWaitOrderExpectation{
		mock:   m.mock,
		params: &AcmeClientMockWaitOrderParams{ctx, url},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.WaitOrder return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockWaitOrderExpectation) Then(op1 *acme.Order, err error) *AcmeClientMock {
	e.results = &AcmeClientMockWaitOrderResults{op1, err}
	return e.mock
}

// WaitOrder implements AcmeClient
func (m *AcmeClientMock) WaitOrder(ctx context.Context, url string) (op1 *acme.Order, err error) {
	atomic.AddUint64(&m.beforeWaitOrderCounter, 1)
	defer atomic.AddUint64(&m.afterWaitOrderCounter, 1)

	for _, e := range m.WaitOrderMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockWaitOrderParams{ctx, url}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.op1, e.results.err
		}
	}

	if m.WaitOrderMock.defaultExpectation != nil {
		atomic.AddUint64(&m.WaitOrderMock.defaultExpectation.Counter, 1)
		want := m.WaitOrderMock.defaultExpectation.params
		got := AcmeClientMockWaitOrderParams{ctx, url}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.WaitOrder got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.WaitOrderMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.WaitOrder")
		}
		return (*results).op1, (*results).err
	}
	if m.funcWaitOrder != nil {
		return m.funcWaitOrder(ctx, url)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.WaitOrder. %v %v", ctx, url)
	return
}

// WaitOrderAfterCounter returns a count of finished AcmeClientMock.WaitOrder invocations
func (m *AcmeClientMock) WaitOrderAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterWaitOrderCounter)
}

// WaitOrderBeforeCounter returns a count of AcmeClientMock.WaitOrder invocations
func (m *AcmeClientMock) WaitOrderBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeWaitOrderCounter)
}

// MinimockWaitOrderDone returns true if the count of the WaitOrder invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockWaitOrderDone() bool {
	for _, e := range m.WaitOrderMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WaitOrderMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWaitOrderCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWaitOrder != nil && atomic.LoadUint64(&m.afterWaitOrderCounter) < 1 {
		return false
	}
	return true
}

// MinimockWaitOrderInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockWaitOrderInspect() {
	for _, e := range m.WaitOrderMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.WaitOrder with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WaitOrderMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterWaitOrderCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.WaitOrder with params: %#v", *m.WaitOrderMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWaitOrder != nil && atomic.LoadUint64(&m.afterWaitOrderCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.WaitOrder")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AcmeClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockAcceptInspect()

		m.MinimockAuthorizeOrderInspect()

		m.MinimockCreateOrderCertInspect()

		m.MinimockGetAuthorizationInspect()

		m.MinimockHTTP01ChallengeResponseInspect()

		m.MinimockRevokeAuthorizationInspect()

		m.MinimockTLSALPN01ChallengeCertInspect()

		m.MinimockWaitAuthorizationInspect()

		m.MinimockWaitOrderInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *AcmeClientMock) MinimockWait(timeout time.Duration) {
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

func (m *AcmeClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAcceptDone() &&
		m.MinimockAuthorizeOrderDone() &&
		m.MinimockCreateOrderCertDone() &&
		m.MinimockGetAuthorizationDone() &&
		m.MinimockHTTP01ChallengeResponseDone() &&
		m.MinimockRevokeAuthorizationDone() &&
		m.MinimockTLSALPN01ChallengeCertDone() &&
		m.MinimockWaitAuthorizationDone() &&
		m.MinimockWaitOrderDone()
}
