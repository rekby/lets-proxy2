//nolint:golint
package cert_manager

import (
	"strings"

	"github.com/rekby/lets-proxy2/internal/domain"

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

func (n CertDescription) DomainNames() []domain.DomainName {
	domains := make([]domain.DomainName, 1, len(n.Subdomains)+1)
	domains[0] = domain.DomainName(n.MainDomain)
	for _, subdomain := range n.Subdomains {
		domains = append(domains, domain.DomainName(subdomain+n.MainDomain))
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

func CertDescriptionFromDomain(domain domain.DomainName, keyType KeyType, autoSubDomains []string) CertDescription {
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
