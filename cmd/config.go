package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rekby/lets-proxy2/internal/config"
	"github.com/rekby/lets-proxy2/internal/domain_checker"
	"github.com/rekby/lets-proxy2/internal/log"
	"github.com/rekby/lets-proxy2/internal/profiler"
	"github.com/rekby/lets-proxy2/internal/proxy"
	"github.com/rekby/lets-proxy2/internal/tlslistener"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

//go:embed static/default-config.toml
var defaultConfigContent []byte

type configType struct {
	General      configGeneral
	Log          logConfig
	Proxy        proxy.Config
	CheckDomains domain_checker.Config
	Listen       tlslistener.Config

	Profiler profiler.Config
	Metrics  config.Config
}

type configGeneral struct {
	IssueTimeout            int
	StorageDir              string
	Subdomains              []string
	AcmeServer              string
	StoreJSONMetadata       bool
	IncludeConfigs          []string
	MaxConfigFilesRead      int
	AllowRSACert            bool
	AllowECDSACert          bool
	AllowInsecureTLSChipers bool
	MinTLSVersion           string
}

//nolint:maligned
type logConfig struct {
	EnableLogToFile   bool
	EnableLogToStdErr bool
	LogLevel          string
	EnableAccessLog   bool
	EnableRotate      bool
	DeveloperMode     bool
	File              string
	RotateBySizeMB    int
	CompressRotated   bool
	MaxDays           int
	MaxCount          int
}

var (
	_config           *configType
	parsedConfigFiles = 0
)

func getConfig(ctx context.Context) *configType {
	if _config == nil {
		logger := zc.LNop(ctx).With(zap.String("config_file", *configFileP))
		logger.Info("Read config")
		_config = &configType{}
		mergeConfigBytes(ctx, _config, defaultConfig(ctx), "default")
		mergeConfigByTemplate(ctx, _config, *configFileP)
		applyMoveConfigDetails(_config)
		applyFlags(ctx, _config)
		logger.Info("Parse configs finished", zap.Int("readed_files", parsedConfigFiles),
			zap.Int("max_read_files", _config.General.MaxConfigFilesRead))

		if *debugLog {
			_config.Log.LogLevel = "debug"
		}
	}
	return _config
}

// Apply command line flags to config
func applyFlags(ctx context.Context, config *configType) {
	if *testAcmeServerP {
		zc.L(ctx).Info("Set test acme server by command line flag")
		config.General.AcmeServer = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}
	if *manualAcmeServer != "" {
		zc.L(ctx).Info("Set force acme server address", zap.String("server", *manualAcmeServer))
		config.General.AcmeServer = *manualAcmeServer
	}
}

func applyMoveConfigDetails(cfg *configType) {
	cfg.Listen.MinTLSVersion = cfg.General.MinTLSVersion
}

func defaultConfig(ctx context.Context) []byte {
	configBytes := bytes.Replace(defaultConfigContent, []byte("\r\n"), []byte("\n"), -1)
	log.DebugCtx(ctx, "Got builtin default config")
	return configBytes
}

func mergeConfigByTemplate(ctx context.Context, c *configType, filepathTemplate string) {
	logger := zc.LNop(ctx).With(zap.String("config_file", filepathTemplate))
	if !hasMeta(filepathTemplate) {
		mergeConfigByFilepath(ctx, c, filepathTemplate)
		return
	}

	filenames, err := filepath.Glob(filepathTemplate)
	log.DebugFatal(logger, err, "Expand config file template",
		zap.String("filepathTemplate", filepathTemplate), zap.Strings("files", filenames))
	for _, filename := range filenames {
		mergeConfigByFilepath(ctx, c, filename)
	}
}

func mergeConfigByFilepath(ctx context.Context, c *configType, filename string) {
	logger := zc.LNop(ctx).With(zap.String("config_file", filename))
	if parsedConfigFiles > c.General.MaxConfigFilesRead {
		logger.Fatal("Exceed max config files read count", zap.Int("MaxConfigFilesRead", c.General.MaxConfigFilesRead))
	}
	parsedConfigFiles++

	var err error
	if !filepath.IsAbs(filename) {
		var filepathNew string
		filepathNew, err = filepath.Abs(filename)
		log.DebugFatal(logger, err, "Convert filepath to absolute",
			zap.String("old", filename), zap.String("new", filepathNew))
		filename = filepathNew
	}

	content, err := ioutil.ReadFile(filename)
	log.DebugFatal(logger, err, "Read filename")

	dir, err := os.Getwd()
	log.DebugFatal(logger, err, "Current workdir", zap.String("dir", dir))
	fileDir := filepath.Dir(filename)
	if dir != fileDir {
		err = os.Chdir(fileDir)
		log.DebugFatal(logger, err, "Chdir to config filename directory")
		defer func() {
			err = os.Chdir(dir)
			log.DebugFatal(logger, err, "Restore workdir to", zap.String("dir", dir))
		}()
	}

	mergeConfigBytes(ctx, c, content, filename)
}

// hasMeta reports whether path contains any of the magic characters
// recognized by Match.
// copy from filepath module
func hasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS != "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}

func mergeConfigBytes(ctx context.Context, c *configType, content []byte, file string) {
	// for prevent loop by existed included
	c.General.IncludeConfigs = nil

	meta, err := toml.Decode(string(content), c)
	if err == nil && len(meta.Undecoded()) > 0 {
		err = fmt.Errorf("unknown fields: %v", meta.Undecoded())
	}
	log.InfoFatal(zc.L(ctx), err, "Parse config file", zap.String("config_file", file))

	if len(c.General.IncludeConfigs) > 0 {
		includeConfigs := c.General.IncludeConfigs // need save because it will reset while merging
		for _, file := range includeConfigs {
			mergeConfigByTemplate(ctx, c, file)
		}
	}
}
