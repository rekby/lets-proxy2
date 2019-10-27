//nolint:golint
package cert_manager

import (
	"strings"

	"go.uber.org/zap"
)

type certNameType string

func (n certNameType) String() string {
	return string(n)
}

func certNameFromDomain(domain DomainName) certNameType {
	return certNameType(strings.TrimPrefix(domain.String(), "www."))
}

func domainNamesFromCertificateName(name certNameType) []DomainName {
	return []DomainName{DomainName(name), DomainName("www." + name)}
}

func logCertName(certName certNameType) zap.Field {
	return zap.String("cert_name", string(certName))
}
