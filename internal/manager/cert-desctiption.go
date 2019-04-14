package manager

type certDescription struct {
	Name   string // Internal name of certificate. It may not related to domain names in future.
	Domain string // Domain in certificate, it will slice in future.
}

func describeCertificate(domain string) certDescription {
	return certDescription{
		domain,
		domain,
	}
}
