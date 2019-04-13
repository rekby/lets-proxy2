package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/rekby/zapcontext"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type proxyType struct {
	GetCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)
	TLSListener    net.Listener // listener for accept TLS connections TODO: allow multiply bindings
	TargetAddr     string       // TODO: allow programmaticaly select target addr

	tlsConfig  tls.Config
	proxy      httputil.ReverseProxy
	httpServer http.Server

	listenerMu      sync.Map
	connListenProxy listenerType
}

// Block until finish work (by context or start error)
// It can stop by cancel context.
// Now Start return immedantly after cancel context - without wait to finish background processes.
// It can change in future.
func (p *proxyType) Start(ctx context.Context) error {
	p.init()

	logger := zc.L(ctx)
	logger.Info("Start proxy")

	go func() {
		_ = p.httpServer.Serve(&p.connListenProxy)
	}()

	go func() {
		<-ctx.Done()
		logger.Debug("Close listener proxy becouse cancel context")
		//_ = p.connListenProxy.Close() - will close by httpserver
		_ = p.httpServer.Shutdown(ctx) // No wait real shutdown
	}()

	for {
		tcpConn, err := p.TLSListener.Accept()
		if err != nil && ctx.Err() != nil {
			if tcpConn != nil {
				//noinspection GoUnhandledErrorResult
				go tcpConn.Close()
			}
			return ctx.Err()
		}

		go p.handleTcpTLSConnection(ctx, tcpConn, err)
	}
}

func (p *proxyType) init() {
	p.proxy.Director = p.director
	p.connListenProxy.connections = make(chan net.Conn)
	p.tlsConfig = tls.Config{GetCertificate: p.GetCertificate}
	p.httpServer.Handler = &p.proxy
}

func (p *proxyType) director(request *http.Request) {
	if request.URL == nil {
		request.URL = &url.URL{}
	}
	request.URL.Scheme = "http"
	request.URL.Host = p.TargetAddr
}

func (p *proxyType) handleTcpTLSConnection(ctx context.Context, conn net.Conn, acceptError error) {
	logger := zc.L(ctx)
	if acceptError != nil {
		logger.Error("Can't accept connection", zap.Error(acceptError))
	}

	connectionUUID := uuid.NewV1()
	logger = logger.With(zap.String("connection_id", connectionUUID.String()))

	ctx = zc.WithLogger(ctx, logger)

	// TODO: save logger/context to some map - for extract in GetCertificate by Conn

	logger.Debug("Accept connection", zap.String("remote_addr", conn.RemoteAddr().String()),
		zap.String("local_addr", conn.LocalAddr().String()))

	tlsConn := tls.Server(conn, &p.tlsConfig)
	err := tlsConn.Handshake()
	if err == nil {
		logger.Debug("Handshake ok")
	} else {
		logger.Info("Can't handshake", zap.Error(err))
	}

	err = p.connListenProxy.Put(tlsConn)
	if err != nil {
		if ctx.Err() != nil {
			logger.Warn("Can't put connection to proxy. Close it.", zap.Error(err))
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
	l.mu.Lock()
	defer l.mu.Unlock()

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
	} else {
		return nil, errors.New("listener closed")
	}
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
