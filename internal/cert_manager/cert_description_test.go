package cert_manager

import (
	"testing"

	"go.uber.org/zap"

	"github.com/maxatome/go-testdeep"
)

func TestCertDescription_CertStoreName(t *testing.T) {
	td := testdeep.NewT(t)
	td.Cmp(CertDescription{MainDomain: "asd.ru", KeyType: KeyRSA}.CertStoreName(), "asd.ru.rsa.cer")
}

func TestCertDescription_DomainNames(t *testing.T) {
	td := testdeep.NewT(t)
	td.Cmp(CertDescription{MainDomain: "asd.ru", KeyType: KeyRSA, Subdomains: []string{"www."}}.DomainNames(), []DomainName{"asd.ru", "www.asd.ru"})
}

func TestCertDescription_KeyStoreName(t *testing.T) {
	td := testdeep.NewT(t)
	td.Cmp(CertDescription{MainDomain: "asd.ru", KeyType: KeyRSA}.KeyStoreName(), "asd.ru.rsa.key")
}

func TestCertDescription_LockName(t *testing.T) {
	td := testdeep.NewT(t)
	td.Cmp(CertDescription{MainDomain: "asd.ru", KeyType: KeyRSA}.LockName(), "asd.ru.lock")
}

func TestCertDescription_MetaStoreName(t *testing.T) {
	td := testdeep.NewT(t)
	td.Cmp(CertDescription{MainDomain: "asd.ru", KeyType: KeyRSA}.MetaStoreName(), "asd.ru.rsa.json")
}

func TestCertDescription_String(t *testing.T) {
	td := testdeep.NewT(t)
	td.Cmp(CertDescription{MainDomain: "asd.ru", KeyType: KeyRSA}.String(), "asd.ru.rsa")
}

func TestCertDescription_ZapField(t *testing.T) {
	td := testdeep.NewT(t)
	cd := CertDescription{MainDomain: "asd.ru", KeyType: KeyRSA}
	td.Cmp(cd.ZapField(), zap.Stringer("cert_name", cd))
}
