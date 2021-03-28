package domain

import (
	"net"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"golang.org/x/net/idna"
	"golang.org/x/xerrors"
)

type DomainName string // Normalized domain name.

func (d DomainName) String() string {
	return string(d)
}

func (d DomainName) ASCII() string {
	return string(d)
}

func (d DomainName) Unicode() string {
	unicode, err := idna.ToUnicode(string(d))
	if err != nil {
		unicode += "[err: " + err.Error() + "]"
	}
	return unicode
}

func (d DomainName) FullString() string {
	return d.Unicode() + " (punycode:" + d.ASCII() + ")"
}

var domainNormalizationProfile = idna.New(idna.ValidateForRegistration(), idna.MapForLookup())

func NormalizeDomain(domain string) (DomainName, error) {
	if strings.Contains(domain, ":") {
		host, _, err := net.SplitHostPort(domain)
		if err != nil {
			return "", xerrors.Errorf("split domain host, port: %w", err)
		}
		domain = host
	}
	domain, err := domainNormalizationProfile.ToASCII(domain)
	domain = strings.TrimSuffix(domain, ".")
	return DomainName(domain), err
}

func LogDomain(domain DomainName) zap.Field {
	return zap.String("domain", domain.FullString())
}

type domainsType []DomainName

func (ss domainsType) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i].FullString())
	}
	return nil
}

func LogDomains(domains []DomainName) zap.Field {
	return logDomainsNamed("domains", domains)
}

func logDomainsNamed(name string, domains []DomainName) zap.Field {
	return zap.Array(name, domainsType(domains))
}
