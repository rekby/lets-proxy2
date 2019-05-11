package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/rekby/lets-proxy2/internal/cert_manager"

	_ "github.com/kardianos/minwinsvc"
	"github.com/rekby/lets-proxy2/internal/acme_client_manager"
	"github.com/rekby/lets-proxy2/internal/cache"
	"github.com/rekby/lets-proxy2/internal/log"
	"github.com/rekby/lets-proxy2/internal/proxy"
	"github.com/rekby/lets-proxy2/internal/tlslistener"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

var VERSION = "custom" // need be var because it redefine by --ldflags "-X main.VERSION" during autobuild

const defaultDirMode = 0700

func main() {
	flag.Parse()

	z, _ := zap.NewProduction()
	globalContext := zc.WithLogger(context.Background(), z)

	if *defaultConfigP {
		fmt.Println(string(defaultConfig(globalContext)))
		os.Exit(0)
	}

	if *versionP {
		fmt.Println(version())
		fmt.Println("Website: https://github.com/rekby/lets-proxy2")
		fmt.Println("Developer: timofey@koolin.ru")
		return
	}

	startProgram(getConfig(globalContext))
}

func version() string {
	return fmt.Sprintf("Version: '%v', Os: '%v', Arch: '%v'", VERSION, runtime.GOOS, runtime.GOARCH)
}

func startProgram(config *configType) {
	logger := initLogger(config.Log)
	ctx := zc.WithLogger(context.Background(), logger)

	logger.Info("StartAutoRenew program version", zap.String("version", version()))

	err := os.MkdirAll(config.General.StorageDir, defaultDirMode)
	log.InfoFatal(logger, err, "Create storage dir")

	storage := &cache.DiskCache{Dir: config.General.StorageDir}
	clientManager := acme_client_manager.New(ctx, storage)
	clientManager.DirectoryURL = config.General.AcmeServer
	acmeClient, err := clientManager.GetClient(ctx)
	log.DebugFatal(logger, err, "Get acme client")

	certManager := cert_manager.New(acmeClient, storage)

	certManager.DomainChecker, err = config.CheckDomains.CreateDomainChecker(ctx)
	log.DebugFatal(logger, err, "Config domain checkers.")

	tlsListener := &tlslistener.ListenersHandler{
		GetCertificate: certManager.GetCertificate,
	}

	err = config.Listen.Apply(ctx, tlsListener)
	log.DebugFatal(logger, err, "Config listeners")

	err = tlsListener.Start(ctx)
	log.DebugFatal(logger, err, "StartAutoRenew tls listener")

	p := proxy.NewHTTPProxy(ctx, tlsListener)
	p.GetContext = func(req *http.Request) (i context.Context, e error) {
		localAddr := req.Context().Value(http.LocalAddrContextKey).(net.Addr)
		return tlsListener.GetConnectionContext(req.RemoteAddr, localAddr.String())
	}
	err = config.Proxy.Apply(ctx, p)
	log.InfoFatal(logger, err, "Apply proxy config")

	err = p.Start()
	var effectiveError = err
	if effectiveError == http.ErrServerClosed {
		effectiveError = nil
	}
	log.DebugErrorCtx(ctx, effectiveError, "Handle request stopped")
}
