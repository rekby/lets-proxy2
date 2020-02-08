//nolint:golint
package cert_manager

import (
	"strings"

	"go.uber.org/zap"
)

type CertDescription struct {
	MainDomain string
	KeyType    KeyType
}

func (n CertDescription) CertStoreName() string {
	return n.MainDomain + "." + n.KeyType.String() + ".cer"
}

func (n CertDescription) DomainNames() []DomainName {
	return []DomainName{DomainName(n.MainDomain), DomainName("www." + n.MainDomain)}
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

func CertDescriptionFromDomain(domain DomainName, keyType KeyType) CertDescription {
	return CertDescription{
		MainDomain: strings.TrimPrefix(domain.String(), "www."),
		KeyType:    keyType,
	}
}
