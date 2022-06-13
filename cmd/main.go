package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/rekby/lets-proxy2/internal/config"

	"github.com/rekby/lets-proxy2/internal/secrethandler"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/rekby/lets-proxy2/internal/metrics"

	"go.uber.org/zap/zapcore"

	"github.com/rekby/lets-proxy2/internal/cert_manager"
	"github.com/rekby/lets-proxy2/internal/profiler"

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

func startMetrics(ctx context.Context, r prometheus.Gatherer, config config.Config, getCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)) error {
	if !config.Enable {
		return nil
	}

	loggerLocal := zc.L(ctx).Named("startMetrics")

	listener := &tlslistener.ListenersHandler{GetCertificate: getCertificate}
	err := config.GetListenConfig().Apply(ctx, listener)
	log.DebugFatal(loggerLocal, err, "Apply listen config")
	if err != nil {
		return xerrors.Errorf("apply config settings to metrics listener: %w", err)
	}

	err = listener.Start(zc.WithLogger(ctx, zc.L(ctx).Named("metrics_listener")), nil)
	log.DebugFatal(loggerLocal, err, "start metrics listener")
	if err != nil {
		return xerrors.Errorf("start metrics listener: %w", err)
	}

	m := metrics.New(zc.L(ctx).Named("metrics"), r)

	secretMetric := secrethandler.New(zc.L(ctx).Named("metrics_secret"), config.GetSecretHandlerConfig(), m)
	go func() {
		defer log.HandlePanic(loggerLocal)

		err := http.Serve(listener, secretMetric)
		var effectiveError = err
		if effectiveError == http.ErrServerClosed {
			effectiveError = nil
		}
		log.DebugDPanic(loggerLocal, effectiveError, "Handle metric stopped")
	}()
	return nil
}

//nolint:funlen
func startProgram(config *configType) {
	logger := initLogger(config.Log)
	ctx := zc.WithLogger(context.Background(), logger)

	logger.Info("StartAutoRenew program version", zap.String("version", version()))

	var registry *prometheus.Registry
	if config.Metrics.Enable {
		registry = prometheus.NewRegistry()
	}

	startProfiler(ctx, config.Profiler)
	err := os.MkdirAll(config.General.StorageDir, defaultDirMode)
	log.InfoFatal(logger, err, "Create storage dir", zap.String("dir", config.General.StorageDir))

	storage := &cache.DiskCache{Dir: config.General.StorageDir}
	clientManager := acme_client_manager.New(ctx, storage)

	clientManager.DirectoryURL = config.General.AcmeServer
	logger.Info("Acme directory", zap.String("url", config.General.AcmeServer))

	_, _, err = clientManager.GetClient(ctx)
	log.InfoFatal(logger, err, "Get acme client")

	certManager := cert_manager.New(clientManager, storage, registry)
	certManager.CertificateIssueTimeout = time.Duration(config.General.IssueTimeout) * time.Second
	certManager.SaveJSONMeta = config.General.StoreJSONMetadata

	certManager.AllowECDSACert = config.General.AllowECDSACert
	certManager.AllowRSACert = config.General.AllowRSACert
	certManager.AllowInsecureTLSChipers = config.General.AllowInsecureTLSChipers

	for _, subdomain := range config.General.Subdomains {
		subdomain = strings.TrimSpace(subdomain)
		subdomain = strings.TrimSuffix(subdomain, ".") + "." // must ends with dot
		certManager.AutoSubdomains = append(certManager.AutoSubdomains, subdomain)
	}

	certManager.DomainChecker, err = config.CheckDomains.CreateDomainChecker(ctx)
	log.DebugFatal(logger, err, "Config domain checkers.")

	err = startMetrics(ctx, registry, config.Metrics, certManager.GetCertificate)
	log.InfoFatalCtx(ctx, err, "start metrics")

	tlsListener := &tlslistener.ListenersHandler{
		GetCertificate: certManager.GetCertificate,
	}

	err = config.Listen.Apply(ctx, tlsListener)
	log.DebugFatal(logger, err, "Config listeners")

	err = tlsListener.Start(ctx, registry)
	log.DebugFatal(logger, err, "StartAutoRenew tls listener")

	config.Proxy.EnableAccessLog = config.Log.EnableAccessLog
	p := proxy.NewHTTPProxy(ctx, tlsListener)
	p.GetContext = func(req *http.Request) (i context.Context, e error) {
		localAddr := req.Context().Value(http.LocalAddrContextKey).(net.Addr)
		return tlsListener.GetConnectionContext(req.RemoteAddr, localAddr.String())
	}

	err = config.Proxy.Apply(ctx, p)
	log.InfoFatal(logger, err, "Apply proxy config")

	go func() {
		defer log.HandlePanic(logger)

		<-ctx.Done()
		err := p.Close()
		log.DebugError(logger, err, "Stop proxy")
	}()

	err = p.Start()
	var effectiveError = err
	if effectiveError == http.ErrServerClosed {
		effectiveError = nil
	}
	log.DebugErrorCtx(ctx, effectiveError, "Handle request stopped")
}

func startProfiler(ctx context.Context, config profiler.Config) {
	logger := zc.L(ctx)

	if !config.Enable {
		logger.Info("Profiler disabled")
		return
	}

	go func() {
		defer log.HandlePanic(logger)

		httpServer := http.Server{
			Addr:    config.BindAddress,
			Handler: profiler.New(logger.Named("profiler"), config),
		}

		logger.Info("Start profiler", zap.String("bind_address", httpServer.Addr))
		err := httpServer.ListenAndServe()
		var logLevel zapcore.Level
		if err == http.ErrServerClosed {
			logLevel = zapcore.InfoLevel
		} else {
			logLevel = zapcore.ErrorLevel
		}
		log.LevelParam(logger, logLevel, "Profiler stopped")
	}()
}
