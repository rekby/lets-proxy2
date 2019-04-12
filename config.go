package main

import (
	"context"
	"github.com/pelletier/go-toml"
	"github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"io/ioutil"
)

type configType struct {
	A              int `toml:"-" comment:"sss"`
	ListenHttpPort int `default:"0" comment:"Port for listen and proxy http traffic. 0 for disable."`
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

}

func defaultConfig(ctx context.Context) []byte {
	config, _ := readConfig(ctx, "")
	configBytes, _ := toml.Marshal(&config)
	return configBytes
}

func readConfig(ctx context.Context, file string) (configType, error) {
	logger := zc.LNop(ctx).With(zap.String("config_file", file))
	var fileBytes []byte
	var err error
	if file == "" {
		logger.Info("Use default config.")
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
