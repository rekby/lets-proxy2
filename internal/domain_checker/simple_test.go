//nolint:golint
package domain_checker

import (
	"regexp"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestTrue(t *testing.T) {
	var _ DomainChecker = True{}

	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)
	res, err := True{}.IsDomainAllowed(ctx, "")
	td.True(res)
	td.CmpNoError(err)
}

func TestFalse(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	var _ DomainChecker = False{}

	td := testdeep.NewT(t)
	res, err := False{}.IsDomainAllowed(ctx, "")
	td.False(res)
	td.CmpNoError(err)
}

func TestRegexp(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	var _ DomainChecker = &Regexp{}

	td := testdeep.NewT(t)
	res, err := NewRegexp(regexp.MustCompile(`\.ru$`)).IsDomainAllowed(ctx, "test.ru")
	td.True(res)
	td.CmpNoError(err)

	res, err = NewRegexp(regexp.MustCompile(`\.ru$`)).IsDomainAllowed(ctx, "test.com")
	td.False(res)
	td.CmpNoError(err)
}
