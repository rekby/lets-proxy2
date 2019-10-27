package zc

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logLevel(level zapcore.Level, logger *zap.Logger, msg string, fields ...zap.Field) {
	const logLevelSkipCallers = 3

	ce := logger.WithOptions(zap.AddCallerSkip(logLevelSkipCallers)).Check(level, msg)
	if ce != nil {
		ce.Write(fields...)
	}
}

func logConditional(levelOK, levelError zapcore.Level, logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	var level zapcore.Level
	if err == nil {
		level = levelOK
	} else {
		level = levelError
		fields = append([]zap.Field{zap.Error(err)}, fields...)
	}
	logLevel(level, logger, msg, fields...)
}

func DebugDPanic(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.DPanicLevel, logger, err, msg, fields...)
}
func DebugDPanicCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.DPanicLevel, L(ctx), err, msg, fields...)
}

func DebugError(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.ErrorLevel, logger, err, msg, fields...)
}
func DebugErrorCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.ErrorLevel, L(ctx), err, msg, fields...)
}

func DebugFatal(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.FatalLevel, logger, err, msg, fields...)
}
func DebugFatalCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.FatalLevel, L(ctx), err, msg, fields...)
}

func DebugInfo(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.InfoLevel, logger, err, msg, fields...)
}
func DebugInfoCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.DebugLevel, zapcore.InfoLevel, L(ctx), err, msg, fields...)
}

func InfoDPanic(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.DPanicLevel, logger, err, msg, fields...)
}
func InfoDPanicCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.DPanicLevel, L(ctx), err, msg, fields...)
}

func InfoError(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.ErrorLevel, logger, err, msg, fields...)
}
func InfoErrorCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.ErrorLevel, L(ctx), err, msg, fields...)
}

func InfoFatal(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.FatalLevel, logger, err, msg, fields...)
}
func InfoFatalCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.FatalLevel, L(ctx), err, msg, fields...)
}

func InfoPanic(logger *zap.Logger, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.PanicLevel, logger, err, msg, fields...)
}
func InfoPanicCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logConditional(zapcore.InfoLevel, zapcore.PanicLevel, L(ctx), err, msg, fields...)
}
