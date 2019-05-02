package main

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/rekby/lets-proxy2/internal/domain_checker"

	"github.com/rekby/lets-proxy2/internal/log"

	"github.com/pelletier/go-toml"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type configType struct {
	IssueTimeout           int    `default:"300" comment:"Seconds for issue every certificate. Cancel issue and return error if timeout."`
	AutoIssueForSubdomains string `default:"www" comment:"Comma separated for subdomains for try get common used subdomains in one certificate."`
	HTTPSListeners         string `default:"[::]:443" comment:"Comma-separated bindings for listen https.\nSupported formats:\n1.2.3.4:443,0.0.0.0:443,[::]:443,[2001:db8::a123]:443"`
	StorageDir             string `default:"storage" comment:"Path to dir, which will store state and certificates"`
	AcmeServer             string `default:"https://acme-v01.api.letsencrypt.org/directory" comment:"Directory url of acme server.\nTest server: https://acme-staging-v02.api.letsencrypt.org/directory"`
	Log                    logConfig
	CheckDomains           domain_checker.Config
}

//nolint:maligned
type logConfig struct {
	EnableLogToFile   bool   `default:"true" comment:""`
	EnableLogToStdErr bool   `default:"true"`
	LogLevel          string `default:"info" comment:"verbose level of log, one of: debug, info, warning, error, fatal"`
	EnableRotate      bool   `default:"true" comment:"Enable self log rotating"`
	DeveloperMode     bool   `default:"false" comment:"Enable developer mode: more stacktraces and panic (stop program) on some internal errors."`
	File              string `default:"lets-proxy2.log" comment:"Path to log file"`
	RotateBySizeMB    int    `default:"100" comment:"Rotate log if current file size more than X MB"`
	CompressRotated   bool   `default:"false" comment:"Compress old log with gzip after rotate"`
	MaxDays           int    `default:"10" comment:"Delete old backups after X days. 0 for disable."`
	MaxCount          int    `default:"10" comment:"Delete old backups if old file number more then X. 0 for disable."`
}

var (
	_config *configType
)

func getConfig(ctx context.Context) *configType {
	if _config == nil {
		logger := zc.LNop(ctx).With(zap.String("config_file", *configFileP))
		logger.Info("Read config")
		config, err := readConfig(ctx, *configFileP)
		if err == nil {
			_config = &config
		} else {
			logger.Fatal("Error while read config.")
		}
		applyFlags(ctx, _config)
	}
	return _config
}

// Apply command line flags to config
func applyFlags(ctx context.Context, config *configType) {
	if *testAcmeServerP {
		zc.L(ctx).Info("Set test acme server by command line flag")
		config.AcmeServer = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}
}

func defaultConfig(ctx context.Context) []byte {
	config, _ := readConfig(ctx, "")
	buf := &bytes.Buffer{}
	err := toml.NewEncoder(buf).Order(toml.OrderPreserve).Encode(config)
	log.DebugDPanicCtx(ctx, err, "Encode default config")
	return buf.Bytes()
}

func readConfig(ctx context.Context, file string) (configType, error) {
	logger := zc.LNop(ctx).With(zap.String("config_file", file))
	var fileBytes []byte
	var err error
	if file == "" {
		logger.Info("Use default config.")
		// Workaround https://github.com/pelletier/go-toml/issues/274
		fileBytes = []byte("[Log]\n[CheckDomains]")
	} else {
		fileBytes, err = ioutil.ReadFile(file)
	}
	if err != nil {
		logger.Error("Can't read config", zap.Error(err))
		return configType{}, err
	}

	var res configType
	err = toml.Unmarshal(fileBytes, &res)
	if err != nil {
		logger.Error("Can't unmarshal config.", zap.Error(err))
		return configType{}, err
	}

	readedConfig, err := toml.Marshal(res)
	if err == nil {
		logger.Info("Read config.", zap.ByteString("config_content", readedConfig))
	} else {
		logger.Error("Can't marshal config", zap.Error(err))
	}
	return res, nil
}
