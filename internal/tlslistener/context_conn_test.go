package tlslistener

import (
	"errors"
	"net"
	"testing"

	"github.com/rekby/lets-proxy2/internal/cert_manager"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

var (
	_ net.Conn                = ContextConnextion{}
	_ cert_manager.GetContext = ContextConnextion{}
)

func TestContextConnextion_Close(t *testing.T) {
	var c ContextConnextion
	var connMock *ConnMock
	var testErr error
	td := testdeep.NewT(t)
	ctx, flush := th.TestContext(t)
	defer flush()

	testErr = errors.New("asd")
	connMock = NewConnMock(td)
	connMock.CloseMock.Expect().Return(testErr)
	c = ContextConnextion{Conn: connMock, Context: ctx}
	td.CmpDeeply(c.Close(), testErr)

	testErr = errors.New("asd2")
	connMock = NewConnMock(td)
	connMock.CloseMock.Expect().Return(nil)
	c = ContextConnextion{Conn: connMock, Context: ctx, CloseFunc: func() error {
		return testErr
	}}
	td.CmpDeeply(c.Close(), testErr)
}

func TestFinalizeContextConnection(t *testing.T) {
	var c ContextConnextion
	var connMock *ConnMock
	td := testdeep.NewT(t)
	ctx, flush := th.TestContext(t)
	defer flush()

	connMock = NewConnMock(td)
	defer func() { _ = connMock.Close() }()

	connMock.CloseMock.Expect().Return(nil)

	closeHandlerCalledWithError := false

	c = ContextConnextion{
		Conn:    connMock,
		Context: ctx,
		connCloseHandler: func(err error) {
			if err != nil {
				closeHandlerCalledWithError = true
			}
		},
	}

	finalizeContextConnection(&c)
	td.True(closeHandlerCalledWithError)
}
