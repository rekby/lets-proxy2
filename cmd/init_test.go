package main

import (
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zc.SetDefaultLogger(logger)
}
