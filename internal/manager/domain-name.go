package manager

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type DomainName string // Normalized domain name.

func (d DomainName) String() string {
	return string(d)
}

func logDomain(domain DomainName) zap.Field {
	return zap.String("domain", domain.String())
}

type domainsType []DomainName

func (ss domainsType) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i].String())
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
	return DomainName(domain)
}
