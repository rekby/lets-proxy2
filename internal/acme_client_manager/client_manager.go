package acme_client_manager

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"sync"

	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/log"

	"github.com/rekby/lets-proxy2/internal/cache"
	"golang.org/x/crypto/acme"
)

const rsaKeyLength = 2048
const stateNameForCache = "acme_account.state.json"

type AcmeManager struct {
	IgnoreCacheLoad bool

	cache cache.Cache

	mu     sync.Mutex
	client *acme.Client
}

func New(cache cache.Cache) *AcmeManager {
	return &AcmeManager{
		cache: cache,
	}
}

type acmeManagerState struct {
	PrivateKey  *rsa.PrivateKey
	AcmeAccount *acme.Account
}

func (m *AcmeManager) GetClient(ctx context.Context) (*acme.Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client != nil {
		return m.client, nil
	}

	if m.cache != nil && !m.IgnoreCacheLoad {
		err := m.loadFromCache(ctx)
		if err != cache.ErrCacheMiss {
			return m.client, err
		}
	}

	err := m.createAccount()
	return m.client, err
}

func (m *AcmeManager) loadFromCache(ctx context.Context) (err error) {
	defer func() {
		var effectiveError error
		if err == cache.ErrCacheMiss {
			effectiveError = nil
		} else {
			effectiveError = err
		}
		log.DebugErrorCtx(ctx, effectiveError, "Load acme manager from cache.", zap.NamedError("raw_err", err))
	}()

	content, err := m.cache.Get(ctx, stateNameForCache)
	if err != nil {
		return err
	}

	var state acmeManagerState
	err = json.Unmarshal(content, &state)
	if err != nil {
		return err
	}

	if state.PrivateKey == nil {
		return errors.New("empty private key")
	}
	if state.AcmeAccount == nil {
		return errors.New("empty account info")
	}

	m.client = &acme.Client{Key: state.PrivateKey}
	return nil
}
