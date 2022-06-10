package tlslistener

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"runtime"
	"sync"

	"github.com/rekby/lets-proxy2/internal/metrics"
	"golang.org/x/xerrors"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/rekby/lets-proxy2/internal/contextlabel"

	"golang.org/x/crypto/acme"

	"github.com/rekby/lets-proxy2/internal/log"

	zc "github.com/rekby/zapcontext"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type ListenersHandler struct {
	GetCertificate        func(*tls.ClientHelloInfo) (*tls.Certificate, error)
	MinTLSVersion         uint16
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
	connectionsContext   map[string]contextInfo

	connectionHandleStart  metrics.ProcessStartFunc
	connectionHandleFinish metrics.ProcessFinishFunc
}

type contextInfo struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
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

func (p *ListenersHandler) Start(ctx context.Context, r prometheus.Registerer) error {
	p.logger = zc.L(ctx)
	p.init()
	p.initMetrics(r)

	p.ctx, p.ctxCancelFunc = context.WithCancel(ctx)

	listenerClosed := make(chan struct{})

	logger := zc.L(ctx)
	logger.Info("StartAutoRenew handleListeners")

	for _, listenerForTLS := range p.ListenersForHandleTLS {
		// handlepanic: in handleConnections
		go handleConnections(ctx, listenerForTLS, p.handleTCPTLSConnection, listenerClosed)
	}

	for _, listener := range p.Listeners {
		// handlepanic: in handleConnections
		go handleConnections(ctx, listener, p.handleTCPConnection, listenerClosed)
	}

	go func() {
		defer log.HandlePanic(logger)

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

func handleConnections(ctx context.Context, l net.Listener, handleFunc func(ctx context.Context, conn net.Conn), listenerClosed chan<- struct{}) {
	logger := zc.L(ctx)
	defer log.HandlePanic(logger)

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
		// handlepanic: in handleFunc
		go handleFunc(ctx, conn)
	}
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
		MinVersion:     p.MinTLSVersion,
	}
	p.connectionsContext = make(map[string]contextInfo)
}

func (p *ListenersHandler) initMetrics(r prometheus.Registerer) {
	p.connectionHandleStart, p.connectionHandleFinish = metrics.ToefCounters(r, "registered_conn", "Registered tcp connections")
}

func (p *ListenersHandler) registerConnection(conn net.Conn, tls bool) ContextConnextion {
	key := conn.RemoteAddr().String() + "-" + conn.LocalAddr().String()

	p.connectionsContextMu.Lock()
	defer p.connectionsContextMu.Unlock()

	ctxStruct, exist := p.connectionsContext[key]
	if exist {
		p.logger.DPanic("Connection already exist in map", zap.String("key", key))
	} else {
		ctxStruct.ctx, ctxStruct.cancelFunc = context.WithCancel(context.Background())
		connectionUUID := uuid.NewV4().String()
		logger := p.logger.With(zap.String("connection_id", connectionUUID))
		ctxStruct.ctx = context.WithValue(ctxStruct.ctx, contextlabel.TLSConnection, tls)
		ctxStruct.ctx = context.WithValue(ctxStruct.ctx, contextlabel.ConnectionID, connectionUUID)
		ctxStruct.ctx = zc.WithLogger(ctxStruct.ctx, logger)
		p.connectionsContext[key] = ctxStruct
	}

	p.connectionHandleStart()

	res := ContextConnextion{
		Context:          ctxStruct.ctx,
		Conn:             conn,
		connCloseHandler: p.connectionHandleFinish,
	}

	res.CloseFunc = func() error {
		p.connectionsContextMu.Lock()
		delete(p.connectionsContext, key)
		p.connectionsContextMu.Unlock()

		runtime.SetFinalizer(&res, nil)

		zc.L(ctxStruct.ctx).WithOptions(zap.AddCallerSkip(2)).Debug("Connection closed.")
		ctxStruct.cancelFunc()

		return conn.Close()
	}
	runtime.SetFinalizer(&res, finalizeContextConnection)
	return res
}

func (p *ListenersHandler) GetConnectionContext(remoteAddr, localAddr string) (context.Context, error) {
	key := remoteAddr + "-" + localAddr

	p.connectionsContextMu.RLock()
	defer p.connectionsContextMu.RUnlock()

	if ctxStruct, ok := p.connectionsContext[key]; ok {
		return ctxStruct.ctx, nil
	}
	return nil, errors.New("not found registered connection")
}

func (p *ListenersHandler) handleTCPConnection(ctx context.Context, conn net.Conn) {
	contextConn := p.registerConnection(conn, false)
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
	contextConn := p.registerConnection(conn, true)
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

func ParseTLSVersion(s string) (uint16, error) {
	switch s {
	case "": // default
		return tls.VersionTLS10, nil
	case "1.0":
		return tls.VersionTLS10, nil
	case "1.1":
		return tls.VersionTLS11, nil
	case "1.2":
		return tls.VersionTLS12, nil
	case "1.3":
		return tls.VersionTLS13, nil
	default:
		return 0, xerrors.Errorf("Unexpected TLS version: '%v'", s)
	}
}
