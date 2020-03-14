package main

import (
	"errors"
	"math"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rekby/lets-proxy2/internal/log"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

type logWriteSyncer struct {
	*lumberjack.Logger
}

func (w logWriteSyncer) Sync() error {
	return w.Logger.Close()
}

func initLogger(config logConfig) *zap.Logger {
	var writers []zapcore.WriteSyncer
	if config.EnableLogToFile {
		lr := &lumberjack.Logger{
			Filename: config.File,
			Compress: config.CompressRotated,
			MaxSize:  config.RotateBySizeMB, MaxAge: config.MaxDays,
			MaxBackups: config.MaxCount,
		}

		if !config.EnableRotate {
			lr.MaxSize = int(math.MaxInt32) // about 2 Petabytes. Really no reachable in this scenario.
		}

		writeSyncer := logWriteSyncer{lr}
		writers = append(writers, writeSyncer)
	}

	if config.EnableLogToStdErr {
		writer, _, err := zap.Open("stderr")
		if err != nil {
			panic("Can't open stderr to log")
		}

		writers = append(writers, writer)
	}

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	logLevel, errLogLevel := parseLogLevel(config.LogLevel)

	core := zapcore.NewCore(encoder, zap.CombineWriteSyncers(writers...), logLevel)
	logger := zap.New(core, getLogOptions(config)...)

	log.InfoError(logger, errLogLevel, "Initialize log on level", zap.Stringer("level", logLevel))

	return logger
}

func parseLogLevel(logLevelS string) (zapcore.Level, error) {
	logLevelS = strings.TrimSpace(logLevelS)
	logLevelS = strings.ToLower(logLevelS)
	switch logLevelS { //nolint:wsl
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, errors.New("undefined log level")
	}
}

func getLogOptions(config logConfig) (res []zap.Option) {
	res = []zap.Option{
		zap.AddCaller(),
	}
	if config.DeveloperMode {
		res = append(res, zap.AddStacktrace(zapcore.WarnLevel), zap.Development())
	} else {
		res = append(res, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	return res
}
