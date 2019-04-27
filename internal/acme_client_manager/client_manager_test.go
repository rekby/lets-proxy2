package acme_client_manager

import (
	"net/url"
	"testing"

	"github.com/maxatome/go-testdeep"

	"github.com/rekby/lets-proxy2/internal/cache"

	"github.com/gojuno/minimock"
	"github.com/rekby/lets-proxy2/internal/th"
)

const testACMEServer = "http://localhost:4000/directory"

func TestClientManager(t *testing.T) {
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
	parsedUrl, _ := url.Parse(testACMEServer)
	manager.DirectoryUrl = *parsedUrl
	client, err := manager.GetClient(ctx)
	td.CmpNoError(err)
	td.NotNil(client)

	client2, err := manager.GetClient(ctx)
	td.CmpNoError(err)
	td.True(client == client2)
}
