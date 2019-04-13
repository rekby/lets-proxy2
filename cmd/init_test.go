package main

import (
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

func init() {
	zc.SetDefaultLogger(zap.NewNop())
}
