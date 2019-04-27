package acme_client_manager

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/log"

	"github.com/rekby/lets-proxy2/internal/cache"
	"golang.org/x/crypto/acme"
)

const rsaKeyLength = 2048

type AcmeManager struct {
	IgnoreCacheLoad      bool
	DirectoryUrl         string
	AgreeFunction        func(tosurl string) bool
	RenewAccountInterval time.Duration

	ctx   context.Context
	cache cache.Cache

	mu      sync.Mutex
	client  *acme.Client
	account *acme.Account
}

func New(ctx context.Context, cache cache.Cache) *AcmeManager {
	return &AcmeManager{
		ctx:                  ctx,
		cache:                cache,
		AgreeFunction:        acme.AcceptTOS,
		RenewAccountInterval: time.Hour * 24,
	}
}

type acmeManagerState struct {
	PrivateKey  *rsa.PrivateKey
	AcmeAccount *acme.Account
}

func (m *AcmeManager) GetClient(ctx context.Context) (*acme.Client, error) {
	if ctx.Err() != nil {
		return nil, errors.New("acme manager context closed")
	}

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

	client := &acme.Client{DirectoryURL: m.DirectoryUrl}
	key, account, err := createAccount(ctx, client, m.AgreeFunction)
	if err != nil {
		return nil, err
	}
	m.account = account
	m.client = client
	state := acmeManagerState{PrivateKey: key, AcmeAccount: account}
	stateBytes, err := json.Marshal(state)
	log.InfoPanicCtx(ctx, err, "Marshal account state to json")
	if m.cache != nil {
		err = m.cache.Put(ctx, certName(m.DirectoryUrl), stateBytes)
		if err != nil {
			return nil, err
		}
	}

	if m.client != nil {
		go m.accountRenew()
	}

	return m.client, err
}

func (m *AcmeManager) accountRenew() {
	ticker := time.NewTicker(m.RenewAccountInterval)
	ctxDone := m.ctx.Done()
	for {
		select {
		case <-ctxDone:
			log.InfoCtx(m.ctx, "Stop renew acme account becouse cancel context", zap.Error(m.ctx.Err()))
			return
		case <-ticker.C:
			newAccount := renewTos(m.ctx, m.client, m.account)
			m.mu.Lock()
			m.account = newAccount
			m.mu.Unlock()
		}
	}
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

	content, err := m.cache.Get(ctx, certName(m.DirectoryUrl))
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

	m.client = &acme.Client{DirectoryURL: m.DirectoryUrl, Key: state.PrivateKey}
	m.account = state.AcmeAccount
	go m.accountRenew()
	return nil
}

func createAccount(ctx context.Context, client *acme.Client, agreeFunction func(tosurl string) bool) (*rsa.PrivateKey, *acme.Account, error) {
	key, err := rsa.GenerateKey(rand.Reader, rsaKeyLength)
	log.InfoDPanicCtx(ctx, err, "Generate account key")

	client.Key = key
	account := &acme.Account{}
	account, err = client.Register(ctx, account, agreeFunction)
	log.InfoErrorCtx(ctx, err, "Register acme account")
	return key, account, err
}

func certName(url string) string {
	sum := md5.Sum([]byte(url))
	sumPrefix := sum[:4]
	return fmt.Sprintf("account_info_%x.client_manager.json", sumPrefix)
}

func renewTos(ctx context.Context, client *acme.Client, account *acme.Account) *acme.Account {
	newAccount, err := client.UpdateReg(ctx, account)
	log.InfoErrorCtx(ctx, err, "Renew acme account")
	if err == nil {
		return newAccount
	}
	return account
}
