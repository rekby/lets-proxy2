package proxy

import (
	"fmt"
	"github.com/egorgasay/cidranger"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestParseTcpMapPair(t *testing.T) {
	td := testdeep.NewT(t)
	var from, to string
	var err error

	from, to, err = parseTCPMapPair("")
	td.CmpDeeply(from, "")
	td.CmpDeeply(to, "")
	td.CmpError(err)

	from, to, err = parseTCPMapPair("a-b")
	td.CmpDeeply(from, "")
	td.CmpDeeply(to, "")
	td.CmpError(err)

	from, to, err = parseTCPMapPair(":123-b")
	td.CmpDeeply(from, "")
	td.CmpDeeply(to, "")
	td.CmpError(err)

	from, to, err = parseTCPMapPair("1.2.3.4-b")
	td.CmpDeeply(from, "")
	td.CmpDeeply(to, "")
	td.CmpError(err)

	from, to, err = parseTCPMapPair("1.2.3.4:123-b")
	td.CmpDeeply(from, "")
	td.CmpDeeply(to, "")
	td.CmpError(err)

	from, to, err = parseTCPMapPair("1.2.3.4:123-2.2.2.2")
	td.CmpDeeply(from, "")
	td.CmpDeeply(to, "")
	td.CmpError(err)

	from, to, err = parseTCPMapPair("1.2.3.4:123-:456")
	td.CmpDeeply(from, "")
	td.CmpDeeply(to, "")
	td.CmpError(err)

	from, to, err = parseTCPMapPair("1.2.3.4:123-2.2.2.2:456")
	td.CmpDeeply(from, "1.2.3.4:123")
	td.CmpDeeply(to, "2.2.2.2:456")
	td.CmpNoError(err)

	from, to, err = parseTCPMapPair("[::1]:123-[::2]:456")
	td.CmpDeeply(from, "[::1]:123")
	td.CmpDeeply(to, "[::2]:456")
	td.CmpNoError(err)
}

func TestConfig_getDefaultTargetDirector(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	var director Director
	var err error

	c := Config{
		DefaultTarget: "",
	}
	director, err = c.getDefaultTargetDirector(ctx)
	td.Nil(director)
	td.CmpError(err)

	c = Config{
		DefaultTarget: "asd",
	}
	director, err = c.getDefaultTargetDirector(ctx)
	td.Nil(director)
	td.CmpError(err)

	c = Config{
		DefaultTarget: ":123",
	}
	director, err = c.getDefaultTargetDirector(ctx)
	td.CmpDeeply(director, NewDirectorSameIP(123))
	td.CmpNoError(err)

	c = Config{
		DefaultTarget: "1.2.3.4",
	}
	director, err = c.getDefaultTargetDirector(ctx)
	td.CmpDeeply(director, NewDirectorHost("1.2.3.4:80"))
	td.CmpNoError(err)

	c = Config{
		DefaultTarget: "::4",
	}
	director, err = c.getDefaultTargetDirector(ctx)
	td.CmpDeeply(director, NewDirectorHost("[::4]:80"))
	td.CmpNoError(err)

	c = Config{
		DefaultTarget: "1.2.3.4:555",
	}
	director, err = c.getDefaultTargetDirector(ctx)
	td.CmpDeeply(director, NewDirectorHost("1.2.3.4:555"))
	td.CmpNoError(err)
}

func TestConfig_getHeadersDirector(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	var director Director
	var err error

	c := Config{
		Headers: []string{},
	}
	director, err = c.getHeadersDirector(ctx)
	td.Nil(director)
	td.CmpNoError(err)

	c = Config{
		Headers: []string{"asd"},
	}
	director, err = c.getHeadersDirector(ctx)
	td.Nil(director)
	td.CmpError(err)

	c = Config{
		Headers: []string{"asd:aaa", "bbb"},
	}
	director, err = c.getHeadersDirector(ctx)
	td.Nil(director)
	td.CmpError(err)

	c = Config{
		Headers: []string{"asd:aaa", "bbb:ccc:hhh", "Ajd:{{qwe}}"},
	}
	director, err = c.getHeadersDirector(ctx)
	td.CmpDeeply(director, NewDirectorSetHeaders(map[string]string{
		"asd": "aaa",
		"bbb": "ccc:hhh",
		"Ajd": "{{qwe}}",
	}))
	td.CmpNoError(err)
}

func TestConfig_getMapDirector(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	var director Director
	var err error

	c := Config{
		TargetMap: []string{},
	}
	director, err = c.getMapDirector(ctx)
	td.Nil(director)
	td.CmpNoError(err)

	c = Config{
		TargetMap: []string{"asd"},
	}
	director, err = c.getMapDirector(ctx)
	td.Nil(director)
	td.CmpError(err)

	c = Config{
		TargetMap: []string{"1.2.3.4-2.3.4.5"},
	}
	director, err = c.getMapDirector(ctx)
	td.Nil(director)
	td.CmpError(err)

	c = Config{
		TargetMap: []string{"1.2.3.4:222-2.3.4.5:333", "asd"},
	}
	director, err = c.getMapDirector(ctx)
	td.Nil(director)
	td.CmpError(err)

	c = Config{
		TargetMap: []string{"1.2.3.4:222-2.3.4.5:333", "[::2]:15-[::5]:91"},
	}
	director, err = c.getMapDirector(ctx)
	td.CmpDeeply(director, NewDirectorDestMap(map[string]string{
		"1.2.3.4:222": "2.3.4.5:333",
		"[::2]:15":    "[::5]:91",
	}))
	td.CmpNoError(err)
}

func TestConfig_getSchemeDirector(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	var director Director
	var err error

	c := &Config{
		HTTPSBackend: false,
	}
	director, err = c.getSchemaDirector(ctx)
	td.CmpNoError(err)
	td.CmpDeeply(director, NewSetSchemeDirector(ProtocolHTTP))

	c = &Config{
		HTTPSBackend: true,
	}
	director, err = c.getSchemaDirector(ctx)
	td.CmpNoError(err)
	td.CmpDeeply(director, NewSetSchemeDirector(ProtocolHTTPS))
}

func TestConfig_Apply(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	var err error
	var p = &HTTPProxy{}

	c := Config{}
	err = c.Apply(ctx, p)
	td.CmpError(err)

	c = Config{
		Headers: []string{"aaa:bbb"},
	}
	p = &HTTPProxy{}
	err = c.Apply(ctx, p)
	td.CmpError(err)

	c = Config{
		DefaultTarget: ":94",
		Headers:       []string{"aaa:bbb"},
	}
	p = &HTTPProxy{}
	err = c.Apply(ctx, p)
	td.CmpNoError(err)
	td.CmpDeeply(p.Director,
		NewDirectorChain(
			NewDirectorSameIP(94),
			NewDirectorSetHeaders(map[string]string{"aaa": "bbb"}),
			NewSetSchemeDirector(ProtocolHTTP),
		),
	)

	c = Config{
		HTTPSBackend:  true,
		DefaultTarget: "1.2.3.4:94",
		TargetMap:     []string{"1.2.3.4:33-4.5.6.7:88"},
		Headers:       []string{"aaa:bbb"},
	}
	p = &HTTPProxy{}
	err = c.Apply(ctx, p)
	td.CmpNoError(err)
	td.CmpDeeply(p.Director, NewDirectorChain(
		NewDirectorHost("1.2.3.4:94"),
		NewDirectorDestMap(map[string]string{"1.2.3.4:33": "4.5.6.7:88"}),
		NewDirectorSetHeaders(map[string]string{"aaa": "bbb"}),
		NewSetSchemeDirector(ProtocolHTTPS),
	))

	// Test backendSchemas

	c = Config{HTTPSBackendIgnoreCert: false}
	p = &HTTPProxy{}
	_ = c.Apply(ctx, p)
	transport := p.HTTPTransport.(Transport)
	transport.IgnoreHTTPSCertificate = false

	c = Config{HTTPSBackendIgnoreCert: true}
	p = &HTTPProxy{}
	_ = c.Apply(ctx, p)
	transport = p.HTTPTransport.(Transport)
	transport.IgnoreHTTPSCertificate = true
}

func TestConfig_getHeadersByIPDirector(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	tests := []struct {
		name    string
		c       Config
		wantErr bool
	}{
		{
			name: "empty",
			c:    Config{},
		},
		{
			name: "oneNetwork",
			c: Config{
				HeadersByIP: map[string][]string{
					"192.168.1.0/24": {
						"User-Agent:PostmanRuntime/7.29.2",
						"Accept:*/*",
						"Accept-Encoding:gzip, deflate, br",
					},
				},
			},
		},
		{
			name: "configError1",
			c: Config{
				HeadersByIP: map[string][]string{
					"192.168.1.0/24": {
						"User-AgentPostmanRuntime/7.29.2",
						"Accept:*/*",
						"Accept-Encoding:gzip, deflate, br",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "5Networks",
			c: Config{
				HeadersByIP: map[string][]string{
					"11.0.0.0/8": {
						"User-Agent:PostmanRuntime/7.29.2",
						"Accept:*/*",
						"Accept-Encoding:gzip, deflate, br",
					},
					"10.55.0.0/24": {
						"Connection:Keep-Alive",
						"Upgrade-Insecure-Requests:1",
						"Cache-Control:no-cache",
					},
					"10.0.1.0/24": {
						"Origin:https://example.com",
						"Content-Type:application/json",
						"Content-Length:123",
					},

					"10.2.0.0/24": {
						"Accept-Encoding:gzip, deflate, br",
						"Accept-Language:en-US,en;q=0.9",
					},
					"fe80:0000:0000:0000::/64": {
						"Accept:*/*",
						"Accept-Encoding:gzip, deflate, br",
						"Accept-Language:en-US,en;q=0.9",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.getHeadersByIPDirector(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getHeadersByIPDirector() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			} else if got == nil {
				return
			}

			ranger := got.(cidranger.Ranger[HTTPHeaders])
			for network, headers := range tt.c.HeadersByIP {
				_, ipnet, err := net.ParseCIDR(network)
				if err != nil {
					t.Fatalf("ParseCIDR error %v", err)
				}

				ip, err := netToIP(ipnet)
				if err != nil {
					t.Fatalf("netToIP error %v", err)
				}

				if ok, err := ranger.Contains(ip); err != nil {
					t.Fatalf("contains error %v", err)
				} else if !ok {
					t.Fatalf("network %s not found", network)
				}

				gotHeaders := make([]string, 0, len(headers))

				err = ranger.IterByIncomingNetworks(ip, func(network net.IPNet, h HTTPHeaders) error {
					if headers == nil {
						return nil
					}

					for _, header := range h {
						gotHeaders = append(gotHeaders, fmt.Sprintf("%s:%s", header.Name, header.Value))
					}
					return nil
				})
				if err != nil {
					t.Fatalf("IterByIncomingNetworks error %v", err)
				}

				if !isTheSameArray(headers, gotHeaders) {
					t.Fatalf("want \n%v \ngot \n%v", headers, gotHeaders)
				}
			}
		})
	}
}

func isTheSameArray[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	tmp := make(map[T]struct{})
	for _, el := range a {
		tmp[el] = struct{}{}
	}
	for _, el := range b {
		if _, ok := tmp[el]; !ok {
			return false
		}
	}
	return true
}

func netToIP(ipnet *net.IPNet) (net.IP, error) {
	ip := ipnet.IP.String()

	sep := ":"
	if ipnet.IP.To4() != nil {
		sep = "."
	}

	split := strings.Split(ip, sep)

	if sep == ":" && split[len(split)-1] == "" {
		return net.ParseIP(strings.Join(split, sep) + "1"), nil
	}
	num, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		return nil, err
	}

	num++
	split[3] = strconv.Itoa(num)

	return net.ParseIP(strings.Join(split, sep)), nil
}
