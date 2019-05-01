package domain_checker

import (
	"errors"
	"regexp"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestTrue(t *testing.T) {
	var _ DomainChecker = True{}

	td := testdeep.NewT(t)
	res, err := True{}.IsDomainAllowed(nil, "")
	td.True(res)
	td.CmpNoError(err)
}

func TestFalse(t *testing.T) {
	var _ DomainChecker = False{}

	td := testdeep.NewT(t)
	res, err := False{}.IsDomainAllowed(nil, "")
	td.False(res)
	td.CmpNoError(err)
}

func TestRegexp(t *testing.T) {
	var _ DomainChecker = &Regexp{}

	td := testdeep.NewT(t)
	res, err := NewRegexp(regexp.MustCompile("\\.ru$")).IsDomainAllowed(nil, "test.ru")
	td.True(res)
	td.CmpNoError(err)

	res, err = NewRegexp(regexp.MustCompile("\\.ru$")).IsDomainAllowed(nil, "test.com")
	td.False(res)
	td.CmpNoError(err)
}

func TestNewNot(t *testing.T) {
	var _ DomainChecker = Not{}

	ctx, cancel := th.TestContext()
	defer cancel()

	td := testdeep.NewT(t)

	m := NewDomainCheckerMock(td)
	defer m.MinimockFinish()

	not := NewNot(m)

	m.IsDomainAllowedMock.Expect(ctx, "asd").Return(true, nil)
	res, err := not.IsDomainAllowed(ctx, "asd")
	td.False(res)
	td.CmpNoError(err)

	m.IsDomainAllowedMock.Expect(ctx, "asd2").Return(false, nil)
	res, err = not.IsDomainAllowed(ctx, "asd2")
	td.True(res)
	td.CmpNoError(err)

	m.IsDomainAllowedMock.Expect(ctx, "qqq").Return(true, errors.New("test"))
	res, err = not.IsDomainAllowed(ctx, "qqq")
	td.False(res)
	td.CmpError(err)

	m.IsDomainAllowedMock.Expect(ctx, "kkk").Return(false, errors.New("test"))
	res, err = not.IsDomainAllowed(ctx, "kkk")
	td.False(res)
	td.CmpError(err)

}
