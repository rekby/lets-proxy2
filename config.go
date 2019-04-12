package main

import (
	"context"
	"github.com/pelletier/go-toml"
	"github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"io/ioutil"
)

type configType struct {
	ListenHttpPort int `default:"0" comment:"Port for listen and proxy http traffic. 0 for disable."`
}

func readConfig(ctx context.Context, file string) (configType, error) {
	logger := zc.LNop(ctx).With(zap.String("config_file", file))
	fileBytes, err := ioutil.ReadFile(file)
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
