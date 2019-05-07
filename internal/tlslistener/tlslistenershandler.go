package tlslistener

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"runtime"
	"sync"

	"github.com/rekby/lets-proxy2/internal/contextlabel"

	"golang.org/x/crypto/acme"

	"github.com/rekby/lets-proxy2/internal/log"

	zc "github.com/rekby/zapcontext"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type ListenersHandler struct {
	GetCertificate        func(*tls.ClientHelloInfo) (*tls.Certificate, error)
	ListenersForHandleTLS []net.Listener // listener which will handle TLS

	// listeners which will not handle TLS, but will proxy to self listener. It is comfortable for listen https and http
	// ports and translate it to one proxy
	Listeners []net.Listener

	NextProtos []string

	ctx           context.Context
	ctxCancelFunc func()
	tlsConfig     tls.Config
	logger        *zap.Logger

	connListenProxy listenerType

	connectionsContextMu sync.RWMutex
	connectionsContext   map[string]context.Context
}

// Implement net.Listener
func (p *ListenersHandler) Accept() (net.Conn, error) {
	return p.connListenProxy.Accept()
}

func (p *ListenersHandler) Close() error {
	p.ctxCancelFunc()
	return p.connListenProxy.Close()
}

func (p *ListenersHandler) Addr() net.Addr {
	return dummyAddr{}
}

// Block until finish work (by context or start error)
// It can stop by cancel context.
// Now StartAutoRenew return immedantly after cancel context - without wait to finish background processes.
// It can change in future.
func (p *ListenersHandler) Start(ctx context.Context) error {
	p.logger = zc.L(ctx)
	p.init()

	p.ctx, p.ctxCancelFunc = context.WithCancel(ctx)

	listenerClosed := make(chan struct{})

	logger := zc.L(ctx)
	logger.Info("StartAutoRenew handleListeners")
	for _, listenerForTLS := range p.ListenersForHandleTLS {
		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				if err != nil {
					if ctx.Err() != nil {
						err = nil
					}
					log.InfoError(logger, err, "Close listener", zap.String("local_addr", l.Addr().String()))
					err = l.Close()
					log.DebugError(logger, err, "Listener closed", zap.String("local_addr", l.Addr().String()))
					listenerClosed <- struct{}{}
					return
				}
				go p.handleTCPTLSConnection(ctx, conn)
			}
		}(listenerForTLS)
	}

	for _, listener := range p.Listeners {
		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				if err != nil {
					if ctx.Err() != nil {
						err = nil
					}
					log.InfoError(logger, err, "Close listener", zap.String("local_addr", l.Addr().String()))
					err = l.Close()
					log.DebugError(logger, err, "Listener closed", zap.String("local_addr", l.Addr().String()))
					listenerClosed <- struct{}{}
					return
				}
				go p.handleTCPConnection(ctx, conn)
			}
		}(listener)

	}

	go func() {
		listenersCount := len(p.ListenersForHandleTLS) + len(p.Listeners)
		for i := 0; i < listenersCount; i++ {
			select {
			case <-ctx.Done():
				return
			case <-listenerClosed:
			}
		}
		if ctx.Err() == nil {
			logger.Warn("All listeners closed. Close Listener handler.")
			_ = p.Close()
		}
	}()

	return nil
}

func (p *ListenersHandler) init() {
	p.connListenProxy.connections = make(chan net.Conn)
	var nextProtos = p.NextProtos
	if nextProtos == nil {
		nextProtos = []string{"h2", "http/1.1"}
	}

	p.tlsConfig = tls.Config{
		GetCertificate: p.GetCertificate,
		NextProtos:     append(nextProtos, acme.ALPNProto),
	}
	p.connectionsContext = make(map[string]context.Context)
}

func (p *ListenersHandler) registerConnection(conn net.Conn) ContextConnextion {
	key := conn.RemoteAddr().String() + "-" + conn.LocalAddr().String()

	p.connectionsContextMu.Lock()
	defer p.connectionsContextMu.Unlock()

	ctx, exist := p.connectionsContext[key]
	if !exist {
		connectionUUID := uuid.NewV4().String()
		logger := p.logger.With(zap.String("connection_id", connectionUUID))
		ctx = context.WithValue(ctx, contextlabel.ConnectionID, connectionUUID)
		ctx = zc.WithLogger(p.ctx, logger)
		p.connectionsContext[key] = ctx
	}

	res := ContextConnextion{
		Context: ctx,
		Conn:    conn,
	}

	// Set finalizer and deregister only for first registered connection
	// for subconnection - return same context
	if exist {
		return res
	}

	res.CloseFunc = func() error {
		p.connectionsContextMu.Lock()
		delete(p.connectionsContext, key)
		p.connectionsContextMu.Unlock()

		runtime.SetFinalizer(&res, nil)

		zc.L(ctx).WithOptions(zap.AddCallerSkip(2)).Debug("Connection closed.")

		return conn.Close()
	}
	runtime.SetFinalizer(&res, finalizeContextConnection)
	return res
}

func (p *ListenersHandler) GetConnectionContext(remoteAddr, localAddr string) (context.Context, error) {
	key := remoteAddr + "-" + localAddr

	p.connectionsContextMu.RLock()
	defer p.connectionsContextMu.RUnlock()

	if ctx, ok := p.connectionsContext[key]; ok {
		return ctx, nil
	}
	return nil, errors.New("not found registered connection")
}

func (p *ListenersHandler) handleTCPConnection(ctx context.Context, conn net.Conn) {
	contextConn := p.registerConnection(conn)
	logger := zc.L(contextConn.Context)

	logger.Debug("Accept connection", zap.String("remote_addr", conn.RemoteAddr().String()),
		zap.String("local_addr", conn.LocalAddr().String()))

	err := p.connListenProxy.Put(contextConn)
	if err != nil {
		if ctx.Err() != nil {
			logger.Error("Can't put connection to proxy. Close it.", zap.Error(err))
		}
		_ = contextConn.Close()
	}
}

func (p *ListenersHandler) handleTCPTLSConnection(ctx context.Context, conn net.Conn) {
	contextConn := p.registerConnection(conn)
	logger := zc.L(contextConn.Context)

	logger.Debug("Accept tls connection", zap.String("remote_addr", conn.RemoteAddr().String()),
		zap.String("local_addr", conn.LocalAddr().String()))

	tlsConn := tls.Server(contextConn, &p.tlsConfig)
	err := tlsConn.Handshake()
	log.DebugInfo(logger, err, "TLS Handshake")

	err = p.connListenProxy.Put(tlsConn)
	if err != nil {
		if ctx.Err() != nil {
			logger.Error("Can't put tls connection to proxy. Close it.", zap.Error(err))
		}
		_ = tlsConn.Close()
	}
}

// listenerType is proxy type - for use already handled connection and send it to http server
// caller MUST NOT Put any connection after Close() call
type listenerType struct {
	connections chan net.Conn

	mu     sync.RWMutex
	closed bool
}

func (l *listenerType) Put(conn net.Conn) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.closed {
		return errors.New("listener closed")
	}

	l.connections <- conn
	return nil
}

func (l *listenerType) Accept() (net.Conn, error) {
	conn, ok := <-l.connections
	if ok {
		return conn, nil
	}
	return nil, errors.New("listener closed")
}

func (l *listenerType) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return errors.New("close closed proxy listener")
	}

	l.closed = true
	close(l.connections)
	return nil
}

func (l *listenerType) Addr() net.Addr {
	return dummyAddr{}
}

type dummyAddr struct{}

func (dummyAddr) Network() string { return "dummy net" }
func (dummyAddr) String() string  { return "dummy addr" }
