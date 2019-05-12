package main

import (
	"context"
	"io/ioutil"

	"github.com/gobuffalo/packr"

	"github.com/rekby/lets-proxy2/internal/tlslistener"

	"github.com/rekby/lets-proxy2/internal/proxy"

	"github.com/rekby/lets-proxy2/internal/domain_checker"

	"github.com/rekby/lets-proxy2/internal/log"

	"github.com/BurntSushi/toml"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type ConfigGeneral struct {
	IssueTimeout           int
	AutoIssueForSubdomains string
	StorageDir             string
	AcmeServer             string
}

//go:generate packr
type configType struct {
	General      ConfigGeneral
	Log          logConfig
	Proxy        proxy.Config
	CheckDomains domain_checker.Config
	Listen       tlslistener.Config
}

//nolint:maligned
type logConfig struct {
	EnableLogToFile   bool
	EnableLogToStdErr bool
	LogLevel          string
	EnableRotate      bool
	DeveloperMode     bool
	File              string
	RotateBySizeMB    int
	CompressRotated   bool
	MaxDays           int
	MaxCount          int
}

var (
	_config *configType
)

func getConfig(ctx context.Context) *configType {
	if _config == nil {
		logger := zc.LNop(ctx).With(zap.String("config_file", *configFileP))
		logger.Info("Read config")
		_config = readConfig(ctx, *configFileP)
		applyFlags(ctx, _config)
	}
	return _config
}

// Apply command line flags to config
func applyFlags(ctx context.Context, config *configType) {
	if *testAcmeServerP {
		zc.L(ctx).Info("Set test acme server by command line flag")
		config.General.AcmeServer = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}
}

func defaultConfig(ctx context.Context) []byte {
	box := packr.NewBox("static")
	configBytes, err := box.Find("default-config.toml")
	log.DebugFatalCtx(ctx, err, "Got builtin default config")
	return configBytes
}

func readConfig(ctx context.Context, file string) *configType {
	var res configType
	logger := zc.LNop(ctx).With(zap.String("config_file", file))
	var fileBytes []byte
	var err error
	fileBytes = defaultConfig(ctx)
	err = toml.Unmarshal(fileBytes, &res)
	log.DebugFatal(logger, err, "Read default config.")

	if file != "" {
		fileBytes, err = ioutil.ReadFile(file)
		log.DebugFatal(logger, err, "Read config file", zap.String("filepath", file))
		err = toml.Unmarshal(fileBytes, &res)
		log.InfoFatal(logger, err, "Parse config file", zap.String("filepath", file))
	}
	return &res
}
