package manager

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/acme"
)

func isTlsAlpn01Hello(hello *tls.ClientHelloInfo) bool {
	return len(hello.SupportedProtos) == 1 && hello.SupportedProtos[0] == acme.ALPNProto
}

func pickChallenge(typ string, chal []*acme.Challenge) *acme.Challenge {
	for _, c := range chal {
		if c.Type == typ {
			return c
		}
	}
	return nil
}

func createCertRequest(key crypto.Signer, commonName string, dnsNames ...string) ([]byte, error) {
	req := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: commonName},
		DNSNames: dnsNames,
	}
	return x509.CreateCertificateRequest(rand.Reader, req, key)
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")
	domain = strings.ToLower(domain)
	return domain
}

// Return valid parced certificate or error
func validCertDer(domains []string, der [][]byte, key crypto.PrivateKey, now time.Time) (cert *tls.Certificate, err error) {
	// parse public part(s)
	var n int
	for _, b := range der {
		n += len(b)
	}
	buf := make([]byte, n)
	n = 0
	for _, b := range der {
		n += copy(buf[n:], b)
	}
	x509Cert, err := x509.ParseCertificates(buf)
	if err != nil || len(x509Cert) == 0 {
		return nil, errors.New("no certificate found in der bytes")
	}

	leaf := x509Cert[0]

	cert = &tls.Certificate{
		PrivateKey:  key,
		Certificate: der,
		Leaf:        leaf,
	}

	return validCertTls(cert, domains, key, now)
}

func validCertTls(cert *tls.Certificate, domains []string, key crypto.PrivateKey, now time.Time) (validCert *tls.Certificate, err error) {
	if cert == nil {
		return nil, errors.New("certificate is nil")
	}

	if cert.Leaf == nil {
		return nil, errors.New("certificate leaf is nil")
	}

	if cert.PrivateKey == nil {
		return nil, errors.New("certificate has no private key")
	}

	if cert.Leaf.PublicKey == nil {
		return nil, errors.New("certificate has no public key")
	}

	// ensure the leaf corresponds to the private key and matches the certKey type
	switch pub := cert.Leaf.PublicKey.(type) {
	case *rsa.PublicKey:
		prv, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key type does not match public key type")
		}
		if pub.N.Cmp(prv.N) != 0 {
			return nil, errors.New("private key does not match public key")
		}
	default:
		return nil, errors.New("unknown public key algorithm")
	}

	// verify the leaf is not expired and matches the domain name
	if now.Before(cert.Leaf.NotBefore) {
		return nil, errors.New("certificate is not valid yet")
	}
	if now.After(cert.Leaf.NotAfter) {
		return nil, errors.New("expired certificate")
	}

	for _, domain := range domains {
		if err := cert.Leaf.VerifyHostname(domain); err != nil {
			return nil, err
		}
	}

	return cert, nil
}
