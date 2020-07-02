package log

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

type certLogger x509.Certificate

const addSkipCallers = 2

func (c *certLogger) String() string {
	cert := (*x509.Certificate)(c)
	if cert == nil {
		return "x509 nil"
	}

	return fmt.Sprintf("Common name: %q, Domains: %q, Expire: %q, SerialNumber: %q",
		cert.Subject.CommonName, cert.DNSNames, cert.NotAfter, cert.Subject.SerialNumber)
}

func Cert(cert *tls.Certificate) zap.Field {
	if cert == nil {
		return zap.String("certificate", "tls nil")
	}
	return CertX509(cert.Leaf)
}

func CertX509(cert *x509.Certificate) zap.Field {
	return zap.Stringer("certificate", (*certLogger)(cert))
}

func DebugInfo(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	debugInfo(logger, err, mess, fields...)
}

func DebugInfoCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	debugInfo(zc.L(ctx), err, mess, fields...)
}

func DebugWarning(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	debugWarning(logger, err, mess, fields...)
}

func debugWarning(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if err == nil {
		logger.Debug(mess, fields...)
	} else {
		logger.Warn(mess, append(fields, zap.Error(err))...)
	}
}

func DebugError(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	debugError(logger, err, mess, fields...)
}

func DebugDPanic(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	debugDpanic(logger, err, mess, fields...)
}

func DebugDPanicCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	debugDpanic(zc.L(ctx), err, mess, fields...)
}

func debugDpanic(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))

	if err == nil {
		logger.Debug(mess, fields...)
	} else {
		logger.DPanic(mess, append(fields, zap.Error(err))...)
	}
}

func DebugFatal(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	debugFatal(logger, err, mess, fields...)
}

func DebugFatalCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	debugFatal(zc.L(ctx), err, mess, fields...)
}

func debugFatal(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))

	if err == nil {
		logger.Debug(mess, fields...)
	} else {
		logger.Fatal(mess, append(fields, zap.Error(err))...)
	}
}

func InfoCtx(ctx context.Context, mess string, fields ...zap.Field) {
	logger := zc.L(ctx)
	logger.WithOptions(zap.AddCallerSkip(1)).Info(mess, fields...)
}

func InfoError(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	infoError(logger, err, mess, fields...)
}

func InfoErrorCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	infoError(zc.L(ctx), err, mess, fields...)
}

func InfoFatal(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	infoFatal(logger, err, mess, fields...)
}

func InfoFatalCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	infoFatal(zc.L(ctx), err, mess, fields...)
}

func infoFatal(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if err == nil {
		logger.Info(mess, fields...)
	} else {
		logger.Fatal(mess, append(fields, zap.Error(err))...)
	}
}

func InfoPanic(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	infoPanic(logger, err, mess, fields...)
}

func InfoDPanicCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	infoDPanic(zc.L(ctx), err, mess, fields...)
}

func InfoPanicCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	infoPanic(zc.L(ctx), err, mess, fields...)
}

func DebugErrorCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	debugError(zc.L(ctx), err, mess, fields...)
}

func debugInfo(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if err == nil {
		logger.Debug(mess, fields...)
	} else {
		logger.Info(mess, append(fields, zap.Error(err))...)
	}
}

func debugError(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if err == nil {
		logger.Debug(mess, fields...)
	} else {
		logger.Error(mess, append(fields, zap.Error(err))...)
	}
}

func infoError(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if err == nil {
		logger.Info(mess, fields...)
	} else {
		logger.Error(mess, append(fields, zap.Error(err))...)
	}
}

func infoDPanic(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if err == nil {
		logger.Info(mess, fields...)
	} else {
		logger.DPanic(mess, append(fields, zap.Error(err))...)
	}
}

func infoPanic(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if err == nil {
		logger.Info(mess, fields...)
	} else {
		logger.Panic(mess, append(fields, zap.Error(err))...)
	}
}

func LevelParam(logger *zap.Logger, level zapcore.Level, mess string, fields ...zap.Field) {
	levelParam(logger, level, mess, fields...)
}

func LevelParamCtx(ctx context.Context, level zapcore.Level, mess string, fields ...zap.Field) {
	logger := zc.L(ctx)
	levelParam(logger, level, mess, fields...)
}

func levelParam(logger *zap.Logger, level zapcore.Level, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(addSkipCallers))
	if ce := logger.Check(level, mess); ce != nil {
		ce.Write(fields...)
	}
}
