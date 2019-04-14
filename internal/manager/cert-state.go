package manager

import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/rekby/zapcontext"
)

type certState struct {
	mu sync.RWMutex

	issueContext       context.Context // nil if no issue process now
	issueContextCancel func()
	cert               *tls.Certificate
}

// Try to lock state for issue certificate.
// It return true if success.
// It return true if state already locked for issue
func (s *certState) StartIssue(context.Context) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.issueContext == nil {
		s.issueContext, s.issueContextCancel = context.WithCancel(context.Background())
		return true
	} else {
		return false
	}
}

// Must call after StartIssue issue certificate
func (s *certState) FinishIssue(ctx context.Context, cert *tls.Certificate) {
	logger := zc.L(ctx)

	s.mu.Lock()
	s.cert = cert
	oldContext, oldCancel := s.issueContext, s.issueContextCancel
	s.issueContext, s.issueContextCancel = nil, nil
	s.mu.Unlock()

	if oldContext == nil {
		zc.L(ctx).DPanic("Finish issue certificate without start it.")
	} else {
		logger.Debug("Cert state set as issue finished")
		oldCancel()
	}
}

func (s *certState) WaitFinishIssue(ctx context.Context) error {
	s.mu.RLock()
	issueContext := s.issueContext
	s.mu.RUnlock()

	if issueContext == nil {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-issueContext.Done():
		return nil
	}
}

func (s *certState) Cert() *tls.Certificate {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.cert
}
