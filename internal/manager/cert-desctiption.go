package manager

import "go.uber.org/zap"

type certNameType string

type certDescription struct {
	Name   certNameType // Internal name of certificate. It may not related to domain names in future.
	Domain string       // Domain in certificate, it will slice in future.
}

func certNameFromDomain(domain DomainName) certNameType {
	return certNameType(domain)
}

func domainNamesFromCertificateName(name certNameType) []DomainName {
	return []DomainName{DomainName(name), DomainName("www." + name)}
}

func logCetName(certName certNameType) zap.Field {
	return zap.String("cert_name", string(certName))
}
