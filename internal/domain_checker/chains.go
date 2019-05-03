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
			return false, nil
		}
	}
	return true, nil
}

func NewAll(checkers ...DomainChecker) All {
	return All(checkers)
}

func NewNot(origin DomainChecker) Not {
	return Not{origin: origin}
}
