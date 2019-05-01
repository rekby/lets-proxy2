//nolint:golint
package domain_checker

import "context"

type Any []DomainChecker

func (any Any) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	for _, checker := range any {
		res, err := checker.IsDomainAllowed(ctx, domain)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

func NewAny(slise []DomainChecker) Any {
	return Any(slise)
}

type All []DomainChecker

func (all All) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	for _, checker := range all {
		res, err := checker.IsDomainAllowed(ctx, domain)
		if err != nil {
			return false, err
		}
		if !res {
			return false, nil
		}
	}
	return true, nil
}

func NewAll(slise []DomainChecker) All {
	return All(slise)
}
