package proxy

import (
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
