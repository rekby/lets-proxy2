package main

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"go.uber.org/zap/zapcore"
)

func TestParseLogLever(t *testing.T) {
	td := testdeep.NewT(t)

	res, err := parseLogLevel("debug")
	td.CmpNoError(err)
	td.CmpDeeply(res, zapcore.DebugLevel)

	res, err = parseLogLevel("info")
	td.CmpNoError(err)
	td.CmpDeeply(res, zapcore.InfoLevel)

	res, err = parseLogLevel("warning")
	td.CmpNoError(err)
	td.CmpDeeply(res, zapcore.WarnLevel)

	res, err = parseLogLevel("error")
	td.CmpNoError(err)
	td.CmpDeeply(res, zapcore.ErrorLevel)

	res, err = parseLogLevel("fatal")
	td.CmpNoError(err)
	td.CmpDeeply(res, zapcore.FatalLevel)

	res, err = parseLogLevel("")
	td.CmpError(err)
	td.CmpDeeply(res, zapcore.InfoLevel)
}
