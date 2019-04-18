package cert_manager

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/idna"
)

type DomainName string // Normalized domain name.

func (d DomainName) String() string {
	return string(d)
}

func (d DomainName) ASCII() string {
	ascii, err := idna.ToASCII(string(d))
	if err != nil {
		ascii += "[err: " + err.Error() + "]"
	}
	return ascii
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

func logDomain(domain DomainName) zap.Field {
	return zap.String("domain", domain.FullString())
}

type domainsType []DomainName

func (ss domainsType) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i].FullString())
	}
	return nil
}

func logDomains(domains []DomainName) zap.Field {
	return logDomainsNamed("domains", domains)
}

func logDomainsNamed(name string, domains []DomainName) zap.Field {
	return zap.Array(name, domainsType(domains))
}

func normalizeDomain(domain string) DomainName {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")
	domain = strings.ToLower(domain)
	domain = DomainName(domain).ASCII()
	return DomainName(domain)
}
