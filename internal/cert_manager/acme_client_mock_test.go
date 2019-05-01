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

	funcAuthorize          func(ctx context.Context, domain string) (ap1 *acme.Authorization, err error)
	afterAuthorizeCounter  uint64
	beforeAuthorizeCounter uint64
	AuthorizeMock          mAcmeClientMockAuthorize

	funcCreateCert          func(ctx context.Context, csr []byte, exp time.Duration, bundle bool) (der [][]byte, certURL string, err error)
	afterCreateCertCounter  uint64
	beforeCreateCertCounter uint64
	CreateCertMock          mAcmeClientMockCreateCert

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
}

// NewAcmeClientMock returns a mock for AcmeClient
func NewAcmeClientMock(t minimock.Tester) *AcmeClientMock {
	m := &AcmeClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.AcceptMock = mAcmeClientMockAccept{mock: m}
	m.AuthorizeMock = mAcmeClientMockAuthorize{mock: m}
	m.CreateCertMock = mAcmeClientMockCreateCert{mock: m}
	m.HTTP01ChallengeResponseMock = mAcmeClientMockHTTP01ChallengeResponse{mock: m}
	m.RevokeAuthorizationMock = mAcmeClientMockRevokeAuthorization{mock: m}
	m.TLSALPN01ChallengeCertMock = mAcmeClientMockTLSALPN01ChallengeCert{mock: m}
	m.WaitAuthorizationMock = mAcmeClientMockWaitAuthorization{mock: m}

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

type mAcmeClientMockAuthorize struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockAuthorizeExpectation
	expectations       []*AcmeClientMockAuthorizeExpectation
}

// AcmeClientMockAuthorizeExpectation specifies expectation struct of the AcmeClient.Authorize
type AcmeClientMockAuthorizeExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockAuthorizeParams
	results *AcmeClientMockAuthorizeResults
	Counter uint64
}

// AcmeClientMockAuthorizeParams contains parameters of the AcmeClient.Authorize
type AcmeClientMockAuthorizeParams struct {
	ctx    context.Context
	domain string
}

// AcmeClientMockAuthorizeResults contains results of the AcmeClient.Authorize
type AcmeClientMockAuthorizeResults struct {
	ap1 *acme.Authorization
	err error
}

// Expect sets up expected params for AcmeClient.Authorize
func (m *mAcmeClientMockAuthorize) Expect(ctx context.Context, domain string) *mAcmeClientMockAuthorize {
	if m.mock.funcAuthorize != nil {
		m.mock.t.Fatalf("AcmeClientMock.Authorize mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockAuthorizeExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockAuthorizeParams{ctx, domain}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.Authorize
func (m *mAcmeClientMockAuthorize) Return(ap1 *acme.Authorization, err error) *AcmeClientMock {
	if m.mock.funcAuthorize != nil {
		m.mock.t.Fatalf("AcmeClientMock.Authorize mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockAuthorizeExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockAuthorizeResults{ap1, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.Authorize method
func (m *mAcmeClientMockAuthorize) Set(f func(ctx context.Context, domain string) (ap1 *acme.Authorization, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.Authorize method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.Authorize method")
	}

	m.mock.funcAuthorize = f
	return m.mock
}

// When sets expectation for the AcmeClient.Authorize which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockAuthorize) When(ctx context.Context, domain string) *AcmeClientMockAuthorizeExpectation {
	if m.mock.funcAuthorize != nil {
		m.mock.t.Fatalf("AcmeClientMock.Authorize mock is already set by Set")
	}

	expectation := &AcmeClientMockAuthorizeExpectation{
		mock:   m.mock,
		params: &AcmeClientMockAuthorizeParams{ctx, domain},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.Authorize return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockAuthorizeExpectation) Then(ap1 *acme.Authorization, err error) *AcmeClientMock {
	e.results = &AcmeClientMockAuthorizeResults{ap1, err}
	return e.mock
}

// Authorize implements AcmeClient
func (m *AcmeClientMock) Authorize(ctx context.Context, domain string) (ap1 *acme.Authorization, err error) {
	atomic.AddUint64(&m.beforeAuthorizeCounter, 1)
	defer atomic.AddUint64(&m.afterAuthorizeCounter, 1)

	for _, e := range m.AuthorizeMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockAuthorizeParams{ctx, domain}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.ap1, e.results.err
		}
	}

	if m.AuthorizeMock.defaultExpectation != nil {
		atomic.AddUint64(&m.AuthorizeMock.defaultExpectation.Counter, 1)
		want := m.AuthorizeMock.defaultExpectation.params
		got := AcmeClientMockAuthorizeParams{ctx, domain}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.Authorize got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.AuthorizeMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.Authorize")
		}
		return (*results).ap1, (*results).err
	}
	if m.funcAuthorize != nil {
		return m.funcAuthorize(ctx, domain)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.Authorize. %v %v", ctx, domain)
	return
}

// AuthorizeAfterCounter returns a count of finished AcmeClientMock.Authorize invocations
func (m *AcmeClientMock) AuthorizeAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterAuthorizeCounter)
}

// AuthorizeBeforeCounter returns a count of AcmeClientMock.Authorize invocations
func (m *AcmeClientMock) AuthorizeBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeAuthorizeCounter)
}

// MinimockAuthorizeDone returns true if the count of the Authorize invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockAuthorizeDone() bool {
	for _, e := range m.AuthorizeMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AuthorizeMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAuthorizeCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAuthorize != nil && atomic.LoadUint64(&m.afterAuthorizeCounter) < 1 {
		return false
	}
	return true
}

// MinimockAuthorizeInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockAuthorizeInspect() {
	for _, e := range m.AuthorizeMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.Authorize with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AuthorizeMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterAuthorizeCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.Authorize with params: %#v", *m.AuthorizeMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAuthorize != nil && atomic.LoadUint64(&m.afterAuthorizeCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.Authorize")
	}
}

type mAcmeClientMockCreateCert struct {
	mock               *AcmeClientMock
	defaultExpectation *AcmeClientMockCreateCertExpectation
	expectations       []*AcmeClientMockCreateCertExpectation
}

// AcmeClientMockCreateCertExpectation specifies expectation struct of the AcmeClient.CreateCert
type AcmeClientMockCreateCertExpectation struct {
	mock    *AcmeClientMock
	params  *AcmeClientMockCreateCertParams
	results *AcmeClientMockCreateCertResults
	Counter uint64
}

// AcmeClientMockCreateCertParams contains parameters of the AcmeClient.CreateCert
type AcmeClientMockCreateCertParams struct {
	ctx    context.Context
	csr    []byte
	exp    time.Duration
	bundle bool
}

// AcmeClientMockCreateCertResults contains results of the AcmeClient.CreateCert
type AcmeClientMockCreateCertResults struct {
	der     [][]byte
	certURL string
	err     error
}

// Expect sets up expected params for AcmeClient.CreateCert
func (m *mAcmeClientMockCreateCert) Expect(ctx context.Context, csr []byte, exp time.Duration, bundle bool) *mAcmeClientMockCreateCert {
	if m.mock.funcCreateCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.CreateCert mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockCreateCertExpectation{}
	}

	m.defaultExpectation.params = &AcmeClientMockCreateCertParams{ctx, csr, exp, bundle}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by AcmeClient.CreateCert
func (m *mAcmeClientMockCreateCert) Return(der [][]byte, certURL string, err error) *AcmeClientMock {
	if m.mock.funcCreateCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.CreateCert mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &AcmeClientMockCreateCertExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &AcmeClientMockCreateCertResults{der, certURL, err}
	return m.mock
}

//Set uses given function f to mock the AcmeClient.CreateCert method
func (m *mAcmeClientMockCreateCert) Set(f func(ctx context.Context, csr []byte, exp time.Duration, bundle bool) (der [][]byte, certURL string, err error)) *AcmeClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the AcmeClient.CreateCert method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the AcmeClient.CreateCert method")
	}

	m.mock.funcCreateCert = f
	return m.mock
}

// When sets expectation for the AcmeClient.CreateCert which will trigger the result defined by the following
// Then helper
func (m *mAcmeClientMockCreateCert) When(ctx context.Context, csr []byte, exp time.Duration, bundle bool) *AcmeClientMockCreateCertExpectation {
	if m.mock.funcCreateCert != nil {
		m.mock.t.Fatalf("AcmeClientMock.CreateCert mock is already set by Set")
	}

	expectation := &AcmeClientMockCreateCertExpectation{
		mock:   m.mock,
		params: &AcmeClientMockCreateCertParams{ctx, csr, exp, bundle},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up AcmeClient.CreateCert return parameters for the expectation previously defined by the When method
func (e *AcmeClientMockCreateCertExpectation) Then(der [][]byte, certURL string, err error) *AcmeClientMock {
	e.results = &AcmeClientMockCreateCertResults{der, certURL, err}
	return e.mock
}

// CreateCert implements AcmeClient
func (m *AcmeClientMock) CreateCert(ctx context.Context, csr []byte, exp time.Duration, bundle bool) (der [][]byte, certURL string, err error) {
	atomic.AddUint64(&m.beforeCreateCertCounter, 1)
	defer atomic.AddUint64(&m.afterCreateCertCounter, 1)

	for _, e := range m.CreateCertMock.expectations {
		if minimock.Equal(*e.params, AcmeClientMockCreateCertParams{ctx, csr, exp, bundle}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.der, e.results.certURL, e.results.err
		}
	}

	if m.CreateCertMock.defaultExpectation != nil {
		atomic.AddUint64(&m.CreateCertMock.defaultExpectation.Counter, 1)
		want := m.CreateCertMock.defaultExpectation.params
		got := AcmeClientMockCreateCertParams{ctx, csr, exp, bundle}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("AcmeClientMock.CreateCert got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.CreateCertMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the AcmeClientMock.CreateCert")
		}
		return (*results).der, (*results).certURL, (*results).err
	}
	if m.funcCreateCert != nil {
		return m.funcCreateCert(ctx, csr, exp, bundle)
	}
	m.t.Fatalf("Unexpected call to AcmeClientMock.CreateCert. %v %v %v %v", ctx, csr, exp, bundle)
	return
}

// CreateCertAfterCounter returns a count of finished AcmeClientMock.CreateCert invocations
func (m *AcmeClientMock) CreateCertAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterCreateCertCounter)
}

// CreateCertBeforeCounter returns a count of AcmeClientMock.CreateCert invocations
func (m *AcmeClientMock) CreateCertBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeCreateCertCounter)
}

// MinimockCreateCertDone returns true if the count of the CreateCert invocations corresponds
// the number of defined expectations
func (m *AcmeClientMock) MinimockCreateCertDone() bool {
	for _, e := range m.CreateCertMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CreateCertMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterCreateCertCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreateCert != nil && atomic.LoadUint64(&m.afterCreateCertCounter) < 1 {
		return false
	}
	return true
}

// MinimockCreateCertInspect logs each unmet expectation
func (m *AcmeClientMock) MinimockCreateCertInspect() {
	for _, e := range m.CreateCertMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AcmeClientMock.CreateCert with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CreateCertMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterCreateCertCounter) < 1 {
		m.t.Errorf("Expected call to AcmeClientMock.CreateCert with params: %#v", *m.CreateCertMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreateCert != nil && atomic.LoadUint64(&m.afterCreateCertCounter) < 1 {
		m.t.Error("Expected call to AcmeClientMock.CreateCert")
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

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AcmeClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockAcceptInspect()

		m.MinimockAuthorizeInspect()

		m.MinimockCreateCertInspect()

		m.MinimockHTTP01ChallengeResponseInspect()

		m.MinimockRevokeAuthorizationInspect()

		m.MinimockTLSALPN01ChallengeCertInspect()

		m.MinimockWaitAuthorizationInspect()
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
		m.MinimockAuthorizeDone() &&
		m.MinimockCreateCertDone() &&
		m.MinimockHTTP01ChallengeResponseDone() &&
		m.MinimockRevokeAuthorizationDone() &&
		m.MinimockTLSALPN01ChallengeCertDone() &&
		m.MinimockWaitAuthorizationDone()
}
