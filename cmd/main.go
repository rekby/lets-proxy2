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
		fmt.Println("Website: https://github.com/rekby/lets-proxy")
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

	logger.Info("Start program version", zap.String("version", version()))

	httpsListeners := createHTTPSListeners(ctx, config.HTTPSListeners)

	if len(httpsListeners) == 0 {
		logger.Fatal("Can't start any listener")
	}

	err := os.MkdirAll(config.StorageDir, defaultDirMode)
	log.InfoFatal(logger, err, "Create storage dir")

	storage := &cache.DiskCache{Dir: config.StorageDir}
	clientManager := acme_client_manager.New(ctx, storage)
	clientManager.DirectoryURL = config.AcmeServer
	acmeClient, err := clientManager.GetClient(ctx)
	log.DebugFatal(logger, err, "Get acme client")

	certManager := cert_manager.New(ctx, acmeClient, storage)

	tlsListener := &tlslistener.ListenersHandler{
		ListenersForHandleTLS: httpsListeners,
		GetCertificate:        certManager.GetCertificate,
	}

	err = tlsListener.Start(ctx)
	log.DebugFatal(logger, err, "Start tls listener")

	p := proxy.NewHTTPProxy(ctx, tlsListener)
	p.GetContext = func(req *http.Request) (i context.Context, e error) {
		localAddr := req.Context().Value(http.LocalAddrContextKey).(net.Addr)
		return tlsListener.GetConnectionContext(req.RemoteAddr, localAddr.String())
	}

	// work in background
	var a chan struct{}
	<-a
}

func createHTTPSListeners(ctx context.Context, bindings string) (res []net.Listener) {
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
