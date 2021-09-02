//nolint:golint
package acme_client_manager

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"math/big"
	"testing"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"

	"golang.org/x/crypto/acme"

	"github.com/rekby/lets-proxy2/internal/cache"

	"github.com/gojuno/minimock/v3"
	"github.com/rekby/lets-proxy2/internal/th"
)

const testACMEServer = "https://acme-server:4001/dir"

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cache.Bytes -o ./cache_bytes_mock_test.go
func TestClientManagerCreateNew(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()

	mc := minimock.NewController(e)
	defer mc.Finish()

	c := NewBytesMock(mc)

	var err error

	//register account
	manager := New(ctx, c)
	manager.httpClient = th.GetHttpClient()

	c.PutMock.Return(nil)
	c.GetMock.Return(nil, cache.ErrCacheMiss)
	manager.DirectoryURL = testACMEServer
	client, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.NotNil(client)

	client2, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.True(client == client2)
}

func TestClientManagerGetFromCache(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()
	ctx = zc.WithLogger(ctx, zap.NewNop().WithOptions(zap.Development()))

	mc := minimock.NewController(e)
	defer mc.Finish()

	c := NewBytesMock(mc)

	var err error

	manager := New(ctx, c)
	defer func() { _ = manager.Close() }()

	state := acmeManagerState{
		AcmeAccount: &acme.Account{},
		PrivateKey: &rsa.PrivateKey{
			D: big.NewInt(123),
		},
	}
	stateBytes, _ := json.Marshal(state)

	c.GetMock.Return(stateBytes, nil)
	client, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.NotNil(client)
	e.CmpDeeply(client.Key, state.PrivateKey)

	client2, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.True(client == client2)

	ctxCancelled, ctxCancelledCancel := context.WithCancel(ctx)
	ctxCancelledCancel()

	client3, err := manager.GetClient(ctxCancelled)
	e.CmpError(err)
	e.Nil(client3)
}
