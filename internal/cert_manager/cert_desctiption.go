//nolint:golint
package cert_manager

import (
	"strings"

	"go.uber.org/zap"
)

type CertDescription struct {
	MainDomain string
	KeyType    KeyType
	Subdomains []string
}

func (n CertDescription) CertStoreName() string {
	return n.MainDomain + "." + n.KeyType.String() + ".cer"
}

func (n CertDescription) DomainNames() []DomainName {
	domains := make([]DomainName, 1, len(n.Subdomains)+1)
	domains[0] = DomainName(n.MainDomain)
	for _, subdomain := range n.Subdomains {
		domains = append(domains, DomainName(subdomain+n.MainDomain))
	}
	return domains
}

func (n CertDescription) KeyStoreName() string {
	return n.MainDomain + "." + n.KeyType.String() + ".key"
}

func (n CertDescription) LockName() string {
	return n.MainDomain + ".lock"
}

func (n CertDescription) MetaStoreName() string {
	return n.MainDomain + "." + n.KeyType.String() + ".json"
}

func (n CertDescription) String() string {
	return n.MainDomain + "." + n.KeyType.String()
}

func (n CertDescription) ZapField() zap.Field {
	return zap.Stringer("cert_name", n)
}

func CertDescriptionFromDomain(domain DomainName, keyType KeyType, autoSubDomains []string) CertDescription {
	mainDomain := domain.String()
	for _, subdomain := range autoSubDomains {
		if strings.HasPrefix(mainDomain, subdomain) {
			mainDomain = strings.TrimPrefix(mainDomain, subdomain)
			break
		}
	}
	return CertDescription{
		MainDomain: mainDomain,
		KeyType:    keyType,
		Subdomains: autoSubDomains,
	}
}
