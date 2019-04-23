package log

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	zc "github.com/rekby/zapcontext"

	"go.uber.org/zap"
)

type certLogger x509.Certificate

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
	} else {
		return CertX509(cert.Leaf)
	}
}

func CertX509(cert *x509.Certificate) zap.Field {
	return zap.Stringer("certificate", (*certLogger)(cert))
}

func DebugError(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	debugError(logger, err, mess, fields...)
}

func DebugErrorCtx(ctx context.Context, err error, mess string, fields ...zap.Field) {
	debugError(zc.L(ctx), err, mess, fields...)
}

func debugError(logger *zap.Logger, err error, mess string, fields ...zap.Field) {
	logger = logger.WithOptions(zap.AddCallerSkip(2))
	if err == nil {
		logger.Debug(mess, fields...)
	} else {
		logger.Error(mess, append(fields, zap.Error(err))...)
	}
}
