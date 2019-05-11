//nolint:golint
package domain_checker

import "context"

type DomainChecker interface {
	// IsDomainAllowed called for check domain for allow certificate
	// It can call concurrency for many domains same time
	// guarantee about domain will correct domain name (as minimum for character set)
	IsDomainAllowed(ctx context.Context, domain string) (bool, error)
}
