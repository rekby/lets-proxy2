package manager

type certNameType string

type certDescription struct {
	Name   certNameType // Internal name of certificate. It may not related to domain names in future.
	Domain string       // Domain in certificate, it will slice in future.
}

func describeCertificate(domain string) certDescription {
	return certDescription{
		certNameType(normalizeDomain(domain)),
		normalizeDomain(domain),
	}
}
