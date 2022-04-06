package cert_manager

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"
	"golang.org/x/crypto/acme"
	"golang.org/x/xerrors"
)

// Test nil in getAuthorized
func TestIssue_134(t *testing.T) {
	t.Parallel()

	ctx, ctxCancel := th.TestContext(t)
	defer ctxCancel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	td := testdeep.NewT(t)

	const testDomain = "test.com"
	testURL := "http://test"
	client := NewAcmeClientMock(mc)
	client.AuthorizeOrderMock.Set(func(ctx context.Context, id []acme.AuthzID, opt ...acme.OrderOption) (op1 *acme.Order, err error) {
		if len(id) == 1 && id[0].Value == testDomain {
			return &acme.Order{
				Status:    acme.StatusPending,
				AuthzURLs: []string{testURL},
			}, nil
		}
		t.Fatalf("Unexpected args: %#v", id)
		return nil, errors.New("Unexpected args")
	})

	testErr := xerrors.New("testErr")
	client.GetAuthorizationMock.Set(func(ctx context.Context, url string) (ap1 *acme.Authorization, err error) {
		if url != testURL {
			t.Fatalf("Unexpected args: %#v", url)
		}
		return nil, testErr
	})

	m := &Manager{}
	res, err := m.createOrderForDomains(ctx, client, testDomain)
	td.Nil(res)
	td.True(xerrors.Is(err, testErr))
}
