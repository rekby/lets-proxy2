//nolint:golint
package acme_client_manager

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	zc "github.com/rekby/zapcontext"

	"golang.org/x/xerrors"

	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/log"

	"github.com/rekby/lets-proxy2/internal/cache"
	"golang.org/x/crypto/acme"
)

const rsaKeyLength = 2048
const renewAccountInterval = time.Hour * 24
const disableDuration = time.Hour

var errClosed = xerrors.Errorf("acmeManager already closed")

//nolint:maligned
type AcmeManager struct {
	IgnoreCacheLoad      bool
	DirectoryURL         string
	AgreeFunction        func(tosurl string) bool
	RenewAccountInterval time.Duration

	ctx                   context.Context
	ctxCancel             context.CancelFunc
	ctxAutorenewCompleted context.Context
	cache                 cache.Bytes
	httpClient            *http.Client

	background       sync.WaitGroup
	mu               sync.Mutex
	lastAccountIndex int
	accounts         []clientAccount
	stateLoaded      bool
	closed           bool
}

type clientAccount struct {
	client  *acme.Client
	account *acme.Account
	enabled bool
}

func New(ctx context.Context, cache cache.Bytes) *AcmeManager {
	ctx, ctxCancel := context.WithCancel(ctx)
	return &AcmeManager{
		ctx:                  ctx,
		ctxCancel:            ctxCancel,
		cache:                cache,
		AgreeFunction:        acme.AcceptTOS,
		RenewAccountInterval: renewAccountInterval,
		httpClient:           http.DefaultClient,
		lastAccountIndex:     -1,
	}
}

func (m *AcmeManager) Close() error {
	logger := zc.L(m.ctx)
	logger.Debug("Start close")
	m.mu.Lock()
	alreadyClosed := m.closed
	ctxAutorenewCompleted := m.ctxAutorenewCompleted
	m.closed = true
	m.ctxCancel()
	m.mu.Unlock()
	logger.Debug("Set closed flag", zap.Any("autorenew_context", ctxAutorenewCompleted))

	if alreadyClosed {
		return xerrors.Errorf("close: %w", errClosed)
	}

	if ctxAutorenewCompleted != nil {
		logger.Debug("Start waiting for complete autorenew")
		<-ctxAutorenewCompleted.Done()
		logger.Debug("Autorenew context closed")
	}
	m.background.Wait()
	return nil
}

func (m *AcmeManager) GetClient(ctx context.Context) (_ *acme.Client, disableFunc func(), err error) {
	if ctx.Err() != nil {
		return nil, nil, errors.New("acme manager context closed")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, nil, xerrors.Errorf("GetClient: %w", errClosed)
	}

	createDisableFunc := func(index int) func() {
		return func() {
			wasEnabled := m.disableAccountSelfSync(index)
			if wasEnabled {
				time.AfterFunc(disableDuration, func() {
					m.accountEnableSelfSync(index)
				})
			}
		}
	}

	if !m.stateLoaded && m.cache != nil && !m.IgnoreCacheLoad {
		err := m.loadFromCache(ctx)
		if err != nil && err != cache.ErrCacheMiss {
			return nil, nil, err
		}
		m.stateLoaded = true
	}

	if index, ok := m.nextEnabledClientIndex(); ok {
		return m.accounts[index].client, createDisableFunc(index), nil
	}

	acc, err := m.registerAccount(ctx)
	m.accounts = append(m.accounts, acc)

	m.background.Add(1)
	// handlepanic: in accountRenewSelfSync
	go func(index int) {
		defer m.background.Done()
		m.accountRenewSelfSync(index)
	}(len(m.accounts) - 1)

	if err != nil {
		return nil, nil, err
	}

	if err = m.saveState(ctx); err != nil {
		return nil, nil, err
	}

	return acc.client, createDisableFunc(len(m.accounts) - 1), nil
}

func (m *AcmeManager) accountRenewSelfSync(index int) {
	logger := zc.L(m.ctx)
	ctx, ctxCancel := context.WithCancel(m.ctx)
	defer ctxCancel()

	m.mu.Lock()
	m.ctxAutorenewCompleted = ctx
	acc := m.accounts[index]
	m.mu.Unlock()

	if m.ctx.Err() != nil {
		return
	}

	logger.Debug("Start account autorenew")

	ticker := time.NewTicker(m.RenewAccountInterval)
	ctxDone := m.ctx.Done()

	for {
		select {
		case <-ctxDone:
			log.InfoCtx(m.ctx, "Stop renew acme account because cancel context", zap.Error(m.ctx.Err()))
			return
		case <-ticker.C:
			var newAccount *acme.Account
			func() {
				defer log.HandlePanic(logger)

				newAccount = renewTos(m.ctx, acc.client, acc.account)
			}()
			acc.account = newAccount
			m.mu.Lock()
			m.accounts[index] = acc
			m.mu.Unlock()
		}
	}
}

func (m *AcmeManager) disableAccountSelfSync(index int) (wasEnabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.accounts[index].enabled {
		m.accounts[index].enabled = false
		return true
	}

	return false
}

func (m *AcmeManager) accountEnableSelfSync(index int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.accounts[index].enabled = true
}

func (m *AcmeManager) initClient() *acme.Client {
	return &acme.Client{DirectoryURL: m.DirectoryURL, HTTPClient: m.httpClient}
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

	content, err := m.cache.Get(ctx, stateName(m.DirectoryURL))
	if err != nil {
		return err
	}

	var state acmeManagerState
	_, err = state.Load(content)
	if err != nil { // nolint:wsl
		return err
	}

	if len(state.Accounts) == 0 {
		return xerrors.Errorf("no accounts in state")
	}

	m.accounts = make([]clientAccount, 0, len(state.Accounts))
	for index, stateAccount := range state.Accounts {
		client := m.initClient()
		client.Key = stateAccount.PrivateKey
		acc := clientAccount{
			client:  client,
			account: stateAccount.AcmeAccount,
			enabled: true,
		}

		m.background.Add(1)
		// handlepanic inside accountRenewSelfSync
		go func(index int) {
			defer m.background.Done()
			m.accountRenewSelfSync(index)
		}(index)
		m.accounts = append(m.accounts, acc)
	}

	return nil
}

func (m *AcmeManager) nextEnabledClientIndex() (int, bool) {
	switch {
	case len(m.accounts) == 0:
		return 0, false
	case len(m.accounts) == 1 && m.accounts[0].enabled:
		return 0, true
	default:
		// pass
	}

	startIndex := m.lastAccountIndex
	if startIndex < 0 {
		startIndex = len(m.accounts) - 1
	}
	index := startIndex
	for {
		index++
		if index >= len(m.accounts) {
			index = 0
		}
		if m.accounts[index].enabled {
			m.lastAccountIndex = index
			return index, true
		}
		if index == startIndex {
			return 0, false
		}
	}
}

func (m *AcmeManager) registerAccount(ctx context.Context) (clientAccount, error) {
	// create account
	client := m.initClient()

	account, err := createAcmeAccount(ctx, client, m.AgreeFunction)
	log.InfoErrorCtx(ctx, err, "Create acme account")
	if err != nil {
		return clientAccount{}, err
	}

	acc := clientAccount{
		client:  client,
		account: account,
		enabled: true,
	}

	return acc, nil
}

func (m *AcmeManager) saveState(ctx context.Context) error {
	var state acmeManagerState
	state.Accounts = make([]acmeAccountState, 0, len(m.accounts))

	for _, acc := range m.accounts {
		state.Accounts = append(state.Accounts, acmeAccountState{PrivateKey: acc.client.Key.(*rsa.PrivateKey), AcmeAccount: acc.account})
	}

	stateBytes, err := json.Marshal(state)
	log.InfoPanicCtx(ctx, err, "Marshal account state to json")

	if m.cache != nil {
		err = m.cache.Put(ctx, stateName(m.DirectoryURL), stateBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// createAcmeAccount create account on acme server and store private key in client.Key
func createAcmeAccount(ctx context.Context, client *acme.Client, agreeFunction func(tosurl string) bool) (*acme.Account, error) {
	key, err := rsa.GenerateKey(rand.Reader, rsaKeyLength)
	log.InfoDPanicCtx(ctx, err, "Generate account key")

	client.Key = key
	account := &acme.Account{}
	account, err = client.Register(ctx, account, agreeFunction)
	log.InfoErrorCtx(ctx, err, "Register acme account")
	return account, err
}

func stateName(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	sum := hasher.Sum(nil)
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
