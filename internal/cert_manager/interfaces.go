//nolint:golint
package cert_manager

import (
	"context"
	"crypto/tls"
	"time"

	"golang.org/x/crypto/acme"
)

type DomainChecker interface {
	// IsDomainAllowed called for check domain for allow certificate
	// It can call concurrency for many domains same time
	// guarantee about domain will correct domain name (as minimum for character set)
	IsDomainAllowed(ctx context.Context, domain string) (bool, error)
}

type AcmeClient interface {
	Accept(ctx context.Context, chal *acme.Challenge) (*acme.Challenge, error)
	Authorize(ctx context.Context, domain string) (*acme.Authorization, error)
	CreateCert(ctx context.Context, csr []byte, exp time.Duration, bundle bool) (der [][]byte, certURL string, err error)
	HTTP01ChallengeResponse(token string) (string, error)
	RevokeAuthorization(ctx context.Context, url string) error
	TLSALPN01ChallengeCert(token, domain string, opt ...acme.CertOption) (cert tls.Certificate, err error)
	WaitAuthorization(ctx context.Context, url string) (*acme.Authorization, error)
}

type managerDefaults struct{}

func (managerDefaults) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	return true, nil
}
