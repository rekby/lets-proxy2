package acme_client_manager

import (
	"crypto/rsa"
	"encoding/json"
	"math/big"
	"testing"

	"golang.org/x/crypto/acme"

	"github.com/maxatome/go-testdeep"

	"github.com/rekby/lets-proxy2/internal/cache"

	"github.com/gojuno/minimock"
	"github.com/rekby/lets-proxy2/internal/th"
)

const testACMEServer = "http://localhost:4000/directory"

func TestClientManagerCreateNew(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	td := testdeep.NewT(t)

	mc := minimock.NewController(td)
	defer mc.Finish()

	c := NewCacheMock(mc)

	var err error

	//register account
	manager := New(ctx, c)
	c.PutMock.Return(nil)
	c.GetMock.Return(nil, cache.ErrCacheMiss)
	manager.DirectoryUrl = testACMEServer
	client, err := manager.GetClient(ctx)
	td.CmpNoError(err)
	td.NotNil(client)

	client2, err := manager.GetClient(ctx)
	td.CmpNoError(err)
	td.True(client == client2)
}

func TestClientManagerGetFromCache(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	td := testdeep.NewT(t)

	mc := minimock.NewController(td)
	defer mc.Finish()

	c := NewCacheMock(mc)

	var err error

	manager := New(ctx, c)

	state := acmeManagerState{
		AcmeAccount: &acme.Account{},
		PrivateKey: &rsa.PrivateKey{
			D: big.NewInt(123),
		},
	}
	stateBytes, _ := json.Marshal(state)

	c.GetMock.Return(stateBytes, nil)
	client, err := manager.GetClient(ctx)
	td.CmpNoError(err)
	td.NotNil(client)
	td.CmpDeeply(client.Key, state.PrivateKey)

	client2, err := manager.GetClient(ctx)
	td.CmpNoError(err)
	td.True(client == client2)
}
