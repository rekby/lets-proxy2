//nolint:golint
package domain_checker

import (
	"context"
	"regexp"
)

type True struct{}

func (True) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	return true, nil
}

type False struct{}

func (False) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	return false, nil
}

type Not struct {
	origin DomainChecker
}

func (n Not) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	res, err := n.origin.IsDomainAllowed(ctx, domain)
	if err == nil {
		return !res, nil
	}
	return false, err
}
func NewNot(origin DomainChecker) Not {
	return Not{origin: origin}
}

type Regexp regexp.Regexp

func (r *Regexp) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	reg := (*regexp.Regexp)(r)
	return reg.MatchString(domain), nil
}

func NewRegexp(r *regexp.Regexp) *Regexp {
	return (*Regexp)(r)
}
