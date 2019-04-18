package cert_manager

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/rekby/lets-proxy2/internal/cache"

	"go.uber.org/zap"

	"github.com/rekby/zapcontext"

	"github.com/gojuno/minimock"

	td "github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"

	"golang.org/x/crypto/acme"
)

const testACMEServer = "http://localhost:4000/directory"
const rsaKeyLength = 2048

const certExample = `-----BEGIN CERTIFICATE-----
MIIFZzCCBE+gAwIBAgITAP+zGFBHh2pLULSUl0LmpPOCgzANBgkqhkiG9w0BAQsF
ADAfMR0wGwYDVQQDDBRoMnBweSBoMmNrZXIgZmFrZSBDQTAeFw0xOTA0MTUxNzM4
MzVaFw0xOTA3MTQxNzM4MzVaMBkxFzAVBgNVBAMTDnd3dy5vbmVjZXJ0LnJ1MIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw9lKtLOZT4byeguM+MoPh3w4
83EiJp0OzbqGNht7acC67Fum+FHeN6fWhPh+b+sq5epBYlzXTqQF9NsPIZfd12/W
gwEYynmIfdgSFOn0eO8mHwp+K0/qsifBf7HAhbSJc/o977UpVdzdiEFXfaqLvfib
uq/fA8N/EtYR7wrI/B5Olj/3bQbodlGSOLOzDqMw/+OYCsf6Pa26+1r8VFRu0H8V
dH95IiveVGsuGewiDSaQ5AxUb0rabr/QSBFH50LeO3bNmzoZnvsPZ5x1ffKd3TeW
HBmQt9YfyKzZGj9A7fcLTjtvGlVxHdsQUgnSnDrLSIEUoXs5qsPjQvgcRlyNUwID
AQABo4ICoDCCApwwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMB
BggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTrfo8PaGWv6THh1yhx
7g6m2by3SDAfBgNVHSMEGDAWgBT7eE8S+WAVgyyfF380GbMuNupBiTBkBggrBgEF
BQcBAQRYMFYwIgYIKwYBBQUHMAGGFmh0dHA6Ly8xMjcuMC4wLjE6NDAwMi8wMAYI
KwYBBQUHMAKGJGh0dHA6Ly9ib3VsZGVyOjQ0MzAvYWNtZS9pc3N1ZXItY2VydDAl
BgNVHREEHjAcggpvbmVjZXJ0LnJ1gg53d3cub25lY2VydC5ydTAnBgNVHR8EIDAe
MBygGqAYhhZodHRwOi8vZXhhbXBsZS5jb20vY3JsMGEGA1UdIARaMFgwCAYGZ4EM
AQIBMEwGAyoDBDBFMCIGCCsGAQUFBwIBFhZodHRwOi8vZXhhbXBsZS5jb20vY3Bz
MB8GCCsGAQUFBwICMBMMEURvIFdoYXQgVGhvdSBXaWx0MIIBAgYKKwYBBAHWeQIE
AgSB8wSB8ADuAHUAKHYaGJAn++880NYaAY12sFBXKcenQRvMvfYE9F1CYVMAAAFq
Iks1ewAABAMARjBEAiBd+b4P29GYXdG0a/qol4PBOGzXv/OC1OvWGhJp+vqQfwIg
AnuwNSOrRIGX1Ur3fdjUmC+S8eI+luCJswslXb0hBakAdQAW6GnB0ZXq18P4lxrj
8HYB94zhtp0xqFIYtoN/MagVCAAAAWoiSzV9AAAEAwBGMEQCIC00wO1Cg+kT442r
Ct3qsR/cxptHAQLscFGKeyr56tK5AiACMQ8/4xbPtzANO5TzWZYgvmNU+Rd5r96O
0J1+9xksFTANBgkqhkiG9w0BAQsFAAOCAQEAY8lvDieKYm5PuFM2JB2bmbjS/vsL
MvhZ4R1jE+jbWQlcrWUZPIEORXcf3FHL8t5nCrEPNB+ei7DT4/vcvAGtNOqc0JFg
PsYLHAUC8EFCVrFWLCB0gS6imD2Wby0UjhOFR0ofynTmGHK/ztrh5BLimqkigVFl
kxAEWpbiBs6NZ0oNpoSX+psrq1RLJ2i3k6IU+7+Z4Co7b7ZxXqvS9lNzAsDECaEa
D8hT6ZnLAnopthcshbPIf9UW7OheE+iJ2zXpXX1ejM/m74TK4Be+gcaVbv6K/iqt
x8I2a4dpcYBjHDo360tyC7+S4DYKt3uV31IBjUf5/+A+9vfSe0N8PUh/bQ==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIERTCCAy2gAwIBAgICElowDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAwwgY2Fj
a2xpbmcgY3J5cHRvZ3JhcGhlciBmYWtlIFJPT1QwHhcNMTYwMzIyMDI0NzUyWhcN
MjEwMzIxMDI0NzUyWjAfMR0wGwYDVQQDDBRoMnBweSBoMmNrZXIgZmFrZSBDQTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMIKR3maBcUSsncXYzQT13D5
Nr+Z3mLxMMh3TUdt6sACmqbJ0btRlgXfMtNLM2OU1I6a3Ju+tIZSdn2v21JBwvxU
zpZQ4zy2cimIiMQDZCQHJwzC9GZn8HaW091iz9H0Go3A7WDXwYNmsdLNRi00o14U
joaVqaPsYrZWvRKaIRqaU0hHmS0AWwQSvN/93iMIXuyiwywmkwKbWnnxCQ/gsctK
FUtcNrwEx9Wgj6KlhwDTyI1QWSBbxVYNyUgPFzKxrSmwMO0yNff7ho+QT9x5+Y/7
XE59S4Mc4ZXxcXKew/gSlN9U5mvT+D2BhDtkCupdfsZNCQWp27A+b/DmrFI9NqsC
AwEAAaOCAX0wggF5MBIGA1UdEwEB/wQIMAYBAf8CAQAwDgYDVR0PAQH/BAQDAgGG
MH8GCCsGAQUFBwEBBHMwcTAyBggrBgEFBQcwAYYmaHR0cDovL2lzcmcudHJ1c3Rp
ZC5vY3NwLmlkZW50cnVzdC5jb20wOwYIKwYBBQUHMAKGL2h0dHA6Ly9hcHBzLmlk
ZW50cnVzdC5jb20vcm9vdHMvZHN0cm9vdGNheDMucDdjMB8GA1UdIwQYMBaAFOmk
P+6epeby1dd5YDyTpi4kjpeqMFQGA1UdIARNMEswCAYGZ4EMAQIBMD8GCysGAQQB
gt8TAQEBMDAwLgYIKwYBBQUHAgEWImh0dHA6Ly9jcHMucm9vdC14MS5sZXRzZW5j
cnlwdC5vcmcwPAYDVR0fBDUwMzAxoC+gLYYraHR0cDovL2NybC5pZGVudHJ1c3Qu
Y29tL0RTVFJPT1RDQVgzQ1JMLmNybDAdBgNVHQ4EFgQU+3hPEvlgFYMsnxd/NBmz
LjbqQYkwDQYJKoZIhvcNAQELBQADggEBAKvePfYXBaAcYca2e0WwkswwJ7lLU/i3
GIFM8tErKThNf3gD3KdCtDZ45XomOsgdRv8oxYTvQpBGTclYRAqLsO9t/LgGxeSB
jzwY7Ytdwwj8lviEGtiun06sJxRvvBU+l9uTs3DKBxWKZ/YRf4+6wq/vERrShpEC
KuQ5+NgMcStQY7dywrsd6x1p3bkOvowbDlaRwru7QCIXTBSb8TepKqCqRzr6YREt
doIw2FE8MKMCGR2p+U3slhxfLTh13MuqIOvTuA145S/qf6xCkRc9I92GpjoQk87Z
v1uhpkgT9uwbRw0Cs5DMdxT/LgIUSfUTKU83GNrbrQNYinkJ77i6wG0=
-----END CERTIFICATE-----
`

const keyExample = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAw9lKtLOZT4byeguM+MoPh3w483EiJp0OzbqGNht7acC67Fum
+FHeN6fWhPh+b+sq5epBYlzXTqQF9NsPIZfd12/WgwEYynmIfdgSFOn0eO8mHwp+
K0/qsifBf7HAhbSJc/o977UpVdzdiEFXfaqLvfibuq/fA8N/EtYR7wrI/B5Olj/3
bQbodlGSOLOzDqMw/+OYCsf6Pa26+1r8VFRu0H8VdH95IiveVGsuGewiDSaQ5AxU
b0rabr/QSBFH50LeO3bNmzoZnvsPZ5x1ffKd3TeWHBmQt9YfyKzZGj9A7fcLTjtv
GlVxHdsQUgnSnDrLSIEUoXs5qsPjQvgcRlyNUwIDAQABAoIBAFS57FfAWtLMzpl9
5b67q3wxgXHPv7Z0u7LEvsspmHpnpnYaMGG9CSWKtoNP/WLtmeFdNmwXPg4HZ4xG
OIWP7akF+QczskXlzeajUy85B0pKK3PCVlLmf+IS0OMtQtyU/eHuoFzTQs6ifjQ0
EGWNImdM5hIdg51dNdwwQBHp2Ik40Rhamr5C0QqrnTDssCipI9Vw+Tw4t8311d/P
gCmOM101te3IMkWbg5AdBigCuJ/LysWqKPojrsW82QxJzZP2w6P+U7qTsnp8MrOM
wcUE0zOFZVY6Hab6CPArfO21OcDbsRzxqw7CIu7NTg43/GZho1tSCGYRaQWDew3b
wFPHwAECgYEA0UVNQv2L9oGQh0Uvt/VA3sIulsLpa4lAMUidpkq9tj4QGEeBSElg
kh//06T6ajfyuRahYcykvEcQqGVbXb/qkG16SdxwWZfxzmEepMJi2AyJV87zVdqa
E7XVL5HRr41XcypXTEr6kw5vmd+98aeNbYBFggnKlO/uW+9s0kC7qIECgYEA75S+
PAMmb8/M0XvOOnkVWcQqLoZoM1b6N3GuBcZs2kZcUvWIl5Wiw259pPUnKoTWa5Ti
uazvt38guEXKfBaMGY3S+24bZuf8tedKewvTvl8xxm5h3F8LEqCG7hKqd3ZjnyQn
rOtDTboJQHSNqTSM6s/XRvY6nr56OAAfrPzjK9MCgYEAseapLLbgYhl45RXS4B61
G+mVs4JU7p8KHBtwMaq+JgwSoKFA7VO7rd1YHPLWErAnPmNXpA1VSd1b8tvfVQ5O
eKMo31tvgyqhXGHBrCy33JSjuSrsP+MLMpBUgBEFYajVW8j115yx8YvHIddL4QAg
QaNW85ohRoXFaxBZwU9YeIECgYEAnBbioiLDRhGytcDdqcb9nBBsEwfKl7hRKRJN
eMHAZa11tS73IRuCgaVZAsIeFFubf1fvJ11+iKSw4p3FwHbILFX0YY9pFvCJ+tGH
+wbHm75VpZyA8ZySkD456p4Kpe5iFWru1oAox1kvcej96oGsVce30CnYI1iiNB4Q
hRn1v7UCgYEAur5ZbW8rxIWUekA5kMCAqt0Hz6GtigwBlsH/ps+L3+4fKOxUkop0
3rc2j04gu7erPAHk9MvhrZFOMn6UlgEqPnck/sx/wcTICL5tJ8bxrVuNBeRTplbf
RJWdLy/H7EyBc86Ak/0zK4WdIHNHQheP2RPMMuT0RFgeZSqcjM4j1wg=
-----END RSA PRIVATE KEY-----
`

type contextConnection struct {
	net.Conn
	context.Context
}

func (c contextConnection) GetContext() context.Context {
	return c.Context
}

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zc.SetDefaultLogger(logger)
}

func TestManager_GetCertificate(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	ctx, flush := th.TestContext()
	defer flush()

	mc := minimock.NewController(t)
	defer mc.Finish()

	cacheMock := NewCacheMock(mc)
	cacheMock.GetMock.Set(func(ctx context.Context, key string) (ba1 []byte, err error) {
		zc.L(ctx).Debug("Cache mock get", zap.String("key", key))
		return nil, cache.ErrCacheMiss
	})
	cacheMock.PutMock.Set(func(ctx context.Context, key string, data []byte) (err error) {
		zc.L(ctx).Debug("Cache mock put", zap.String("key", key))
		return nil
	})

	manager := New(ctx, createTestClient(t))
	manager.Cache = cacheMock

	lisneter, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5001})

	if err != nil {
		t.Fatal(err)
	}
	defer lisneter.Close()

	go func() {
		counter := 0
		for {
			conn, err := lisneter.Accept()
			if conn != nil {
				t.Log("incoming connection")
				ctx := zc.WithLogger(context.Background(), logger.With(zap.Int("connection_id", counter)))

				tlsConn := tls.Server(contextConnection{conn, ctx}, &tls.Config{
					NextProtos: []string{
						"h2", "http/1.1", // enable HTTP/2
						acme.ALPNProto, // enable tls-alpn ACME challenges
					},
					GetCertificate: manager.GetCertificate,
				})

				err := tlsConn.Handshake()
				if err == nil {
					t.Log("Handshake ok")
				} else {
					t.Error(err)
				}

				err = conn.Close()
				if err != nil {
					t.Error(err)
				}
			}
			if err != nil {
				break
			}
		}
	}()

	t.Run("OneCert", func(t *testing.T) {
		domain := "onecert.ru"

		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
		if cert.Leaf.VerifyHostname("www."+domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	t.Run("punycode-domain", func(t *testing.T) {
		domain := "xn--80adjurfhd.xn--p1ai" // проверка.рф

		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
		if cert.Leaf.VerifyHostname("www."+domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	t.Run("OneCertCamelCase", func(t *testing.T) {
		domain := "onecertCamelCase.ru"
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	t.Run("ParallelCert", func(t *testing.T) {
		// change top loevel logger
		// no parallelize
		oldLogger := logger
		logger = zap.NewNop()
		defer func() {
			logger = oldLogger
		}()

		domain := "ParallelCert.ru"
		const cnt = 100

		chanCerts := make(chan *tls.Certificate, cnt)

		var wg sync.WaitGroup
		wg.Add(cnt)

		for i := 0; i < cnt; i++ {
			go func() {
				cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
				if err != nil {
					t.Fatal(err)
				}
				chanCerts <- cert
				wg.Done()
			}()
		}

		wg.Wait()
		cert := <-chanCerts
		for i := 0; i < len(chanCerts)-1; i++ {
			cert2 := <-chanCerts
			td.CmpDeeply(t, cert2, cert)
		}
	})
}

func createTestClient(t *testing.T) *acme.Client {
	_, err := http.Get(testACMEServer)
	if err != nil {
		t.Fatalf("Can't connect to buoulder server: %q", err)
	}

	client := acme.Client{}
	client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client.DirectoryURL = testACMEServer
	client.Key, _ = rsa.GenerateKey(rand.Reader, rsaKeyLength)
	_, err = client.Register(context.Background(), &acme.Account{}, func(tosURL string) bool {
		return true
	})

	if err != nil {
		t.Fatal("Can't initialize acme client.")
	}
	return &client
}

func TestStoreCertificate(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	key, _ := rsa.GenerateKey(rand.Reader, 512)

	cert := &tls.Certificate{Certificate: [][]byte{
		{1, 2, 3},
		{4, 5, 6},
	},
		PrivateKey: key,
	}

	mc := minimock.NewController(t)
	cacheMock := NewCacheMock(mc).PutMock.Set(func(ctx context.Context, key string, data []byte) (err error) {
		fmt.Println(key)
		fmt.Println(string(data))
		return nil
	})

	storeCertificate(ctx, cacheMock, "asd", cert)
}
