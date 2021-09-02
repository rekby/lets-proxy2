package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

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

func TestInitLogger(t *testing.T) {
	e, _, flush := th.NewEnv(t)
	defer flush()

	tmpDir := th.TmpDir(e)

	// LogToFile, logLevel
	logFile := filepath.Join(tmpDir, "log.txt")
	config := logConfig{
		EnableLogToFile: true,
		File:            logFile,
		LogLevel:        "warning",
	}
	logger := initLogger(config)
	testError := "errorTest"
	testInfo := "infoTest"
	logger.Error(testError)
	logger.Info(testInfo)
	logger.Sync()

	fileBytes, err := ioutil.ReadFile(logFile)
	e.CmpNoError(err)
	e.True(strings.Contains(string(fileBytes), testError))
	e.False(strings.Contains(string(fileBytes), testInfo))

	// DevelMode
	config = logConfig{DeveloperMode: false, LogLevel: "info", EnableLogToStdErr: true}
	logger = initLogger(config)
	logger.DPanic(testError)

	config = logConfig{DeveloperMode: true, LogLevel: "info"}
	logger = initLogger(config)
	e.CmpPanic(func() {
		logger.DPanic(testError)
	}, testError)
}
