package domain_checker

import (
	"errors"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/gojuno/minimock"

	"github.com/maxatome/go-testdeep"
)

func TestAny(t *testing.T) {
	var _ DomainChecker = Any{}

	ctx, cancel := th.TestContext()
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)
	defer mc.Finish()

	m1 := NewDomainCheckerMock(mc)
	m2 := NewDomainCheckerMock(mc)
	any := NewAny([]DomainChecker{m1, m2})
	var res bool
	var err error

	m1.IsDomainAllowedMock.Expect(ctx, "aaa").Return(true, nil)
	m1.IsDomainAllowedMock.Expect(ctx, "aaa").Return(true, nil)
	res, err = any.IsDomainAllowed(ctx, "aaa")
	td.True(res)
	td.CmpNoError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "bbb").Return(false, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "bbb").Return(true, nil)
	res, err = any.IsDomainAllowed(ctx, "bbb")
	td.True(res)
	td.CmpNoError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "www").Return(false, errors.New("test"))
	res, err = any.IsDomainAllowed(ctx, "www")
	td.False(res)
	td.CmpError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "fadf").Return(false, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "fadf").Return(true, errors.New("test"))
	res, err = any.IsDomainAllowed(ctx, "fadf")
	td.False(res)
	td.CmpError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "edc").Return(false, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "edc").Return(false, errors.New("test"))
	res, err = any.IsDomainAllowed(ctx, "edc")
	td.False(res)
	td.CmpError(err)

}

func TestAll(t *testing.T) {
	var _ DomainChecker = Any{}

	ctx, cancel := th.TestContext()
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)
	defer mc.Finish()

	m1 := NewDomainCheckerMock(mc)
	m2 := NewDomainCheckerMock(mc)
	any := NewAll([]DomainChecker{m1, m2})
	var res bool
	var err error

	m1.IsDomainAllowedMock.Expect(ctx, "aaa").Return(true, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "aaa").Return(true, nil)
	res, err = any.IsDomainAllowed(ctx, "aaa")
	td.True(res)
	td.CmpNoError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "bbb").Return(false, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "bbb").Return(true, nil)
	res, err = any.IsDomainAllowed(ctx, "bbb")
	td.False(res)
	td.CmpNoError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "www").Return(false, errors.New("test"))
	res, err = any.IsDomainAllowed(ctx, "www")
	td.False(res)
	td.CmpError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "fadf").Return(false, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "fadf").Return(true, errors.New("test"))
	res, err = any.IsDomainAllowed(ctx, "fadf")
	td.False(res)
	td.CmpNoError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "edc").Return(false, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "edc").Return(false, errors.New("test"))
	res, err = any.IsDomainAllowed(ctx, "edc")
	td.False(res)
	td.CmpNoError(err)

}
