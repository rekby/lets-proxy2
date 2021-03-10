//nolint:golint
package domain_checker

import (
	"context"

	"github.com/rekby/lets-proxy2/internal/log"
)

type Any []DomainChecker

func (any Any) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	for _, checker := range any {
		res, err := checker.IsDomainAllowed(ctx, domain)
		if err != nil {
			return false, err
		}
		if res {
			log.DebugCtx(ctx, "Allowed by 'any' rule chain: some of subrules allow domain (details in debug log)")
			return true, nil
		}
	}
	log.InfoCtx(ctx, "Deny by 'any' rule chain: nothing of subrules allow domain.")
	return false, nil
}

func NewAny(checkers ...DomainChecker) Any {
	return Any(checkers)
}

type All []DomainChecker

func (all All) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	for _, checker := range all {
		res, err := checker.IsDomainAllowed(ctx, domain)
		if err != nil {
			return false, err
		}
		if !res {
			log.InfoCtx(ctx, "Deny by 'all' chain: any of subrules denied domain.")
			return false, nil
		}
	}
	log.DebugCtx(ctx, "Allow by 'all' chain: all of subrules allow domain.")
	return true, nil
}

func NewAll(checkers ...DomainChecker) All {
	return All(checkers)
}

func NewNot(origin DomainChecker) Not {
	return Not{origin: origin}
}
