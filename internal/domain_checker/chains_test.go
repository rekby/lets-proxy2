//nolint:golint
package domain_checker

import (
	"errors"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/gojuno/minimock/v3"

	"github.com/maxatome/go-testdeep"
)

func TestAny(t *testing.T) {
	var _ DomainChecker = Any{}

	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)
	defer mc.Finish()

	m1 := NewDomainCheckerMock(mc)
	m2 := NewDomainCheckerMock(mc)
	any := NewAny(m1, m2)
	var res bool
	var err error

	res, err = NewAny().IsDomainAllowed(ctx, "bbb")
	td.False(res)
	td.CmpNoError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "aaa").Return(true, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "aaa").Return(true, nil)
	res, err = any.IsDomainAllowed(ctx, "aaa")
	td.True(res)
	td.CmpNoError(err)

	m1.IsDomainAllowedMock.Expect(ctx, "aaa").Return(false, nil)
	m2.IsDomainAllowedMock.Expect(ctx, "aaa").Return(false, nil)
	res, err = any.IsDomainAllowed(ctx, "aaa")
	td.False(res)
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

	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)
	defer mc.Finish()

	m1 := NewDomainCheckerMock(mc)
	m2 := NewDomainCheckerMock(mc)
	any := NewAll(m1, m2)
	var res bool
	var err error

	res, err = NewAll().IsDomainAllowed(ctx, "aaa")
	td.True(res)
	td.CmpNoError(err)

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

func TestNot(t *testing.T) {
	var _ DomainChecker = Not{}

	ctx, cancel := th.TestContext(t)
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
