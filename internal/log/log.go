package log

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"go.uber.org/zap"
)

func Domain(domain string) zap.Field {
	return zap.String("domain", domain)
}

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
