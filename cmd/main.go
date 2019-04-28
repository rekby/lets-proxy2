package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"

	_ "github.com/kardianos/minwinsvc"
	"github.com/rekby/lets-proxy2/internal/acme_client_manager"
	"github.com/rekby/lets-proxy2/internal/cache"
	"github.com/rekby/lets-proxy2/internal/cert_manager"
	"github.com/rekby/lets-proxy2/internal/log"
	"github.com/rekby/lets-proxy2/internal/proxy"
	"github.com/rekby/lets-proxy2/internal/tlslistener"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

var VERSION = "custom" // need be var becouse it redefine by --ldflags "-X main.VERSION" during autobuild

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
		fmt.Printf("Version: '%v', Os: '%v', Arch: '%v'\n", VERSION, runtime.GOOS, runtime.GOARCH)
		return
	}

	startProgram(getConfig(globalContext))
}

func startProgram(config *configType) {
	logger := initLogger(config.Log)
	ctx := zc.WithLogger(context.Background(), logger)
	httpsListeners := createHttpsListeners(ctx, config.HttpsListeners)

	if len(httpsListeners) == 0 {
		logger.Fatal("Can't start any listener")
	}

	err := os.MkdirAll(config.StorageDir, defaultDirMode)
	log.InfoFatal(logger, err, "Create storage dir")

	storage := &cache.DiskCache{Dir: config.StorageDir}
	clientManager := acme_client_manager.New(ctx, storage)
	clientManager.DirectoryUrl = config.AcmeServer
	acmeClient, err := clientManager.GetClient(ctx)
	log.DebugFatal(logger, err, "Get acme client")

	certManager := cert_manager.New(ctx, acmeClient, storage)

	tlsListener := &tlslistener.ListenersHandler{
		ListenersForHandleTls: httpsListeners,
		GetCertificate:        certManager.GetCertificate,
	}

	err = tlsListener.Start(ctx)
	log.DebugFatal(logger, err, "Start tls listener")

	p := proxy.NewHttpProxy(ctx, tlsListener)
	p.GetContext = func(req *http.Request) (i context.Context, e error) {
		localAddr := req.Context().Value(http.LocalAddrContextKey).(net.Addr)
		return tlsListener.GetConnectionContext(req.RemoteAddr, localAddr.String())
	}

	// work in background
	var a chan struct{}
	<-a
}

func createHttpsListeners(ctx context.Context, bindings string) (res []net.Listener) {
	addresses := strings.Split(bindings, ",")
	for _, address := range addresses {
		address = strings.TrimSpace(address)
		if address == "" {
			continue
		}
		var lc net.ListenConfig
		listener, err := lc.Listen(ctx, "tcp", address)
		log.InfoErrorCtx(ctx, err, "Start https listener", zap.String("address", address))
		if err == nil {
			res = append(res, listener)
		}
	}
	return res
}
