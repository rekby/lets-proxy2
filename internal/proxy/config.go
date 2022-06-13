package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rekby/lets-proxy2/internal/log"

	"go.uber.org/zap"

	zc "github.com/rekby/zapcontext"
)

const defaultHTTPPort = 80

//nolint:lll
type Config struct {
	DefaultTarget           string
	TargetMap               []string
	Headers                 []string
	KeepAliveTimeoutSeconds int
	HTTPSBackend            bool
	HTTPSBackendIgnoreCert  bool
	EnableAccessLog         bool
}

func (c *Config) Apply(ctx context.Context, p *HTTPProxy) error {
	var resErr error

	var chain []Director
	appendDirector := func(f func(ctx context.Context) (Director, error)) {
		if resErr != nil {
			return
		}
		director, err := f(ctx)
		resErr = err

		chain = append(chain, director)
	}

	appendDirector(c.getDefaultTargetDirector)
	appendDirector(c.getMapDirector)
	appendDirector(c.getHeadersDirector)
	appendDirector(c.getSchemaDirector)
	p.HTTPTransport = Transport{c.HTTPSBackendIgnoreCert}
	p.EnableAccessLog = c.EnableAccessLog

	if resErr != nil {
		zc.L(ctx).Error("Can't parse proxy config", zap.Error(resErr))
		return resErr
	}

	chainDirector := NewDirectorChain(chain...)
	p.Director = chainDirector
	p.IdleTimeout = time.Duration(c.KeepAliveTimeoutSeconds) * time.Second
	return nil
}

func (c *Config) getDefaultTargetDirector(ctx context.Context) (Director, error) {
	logger := zc.L(ctx)

	var defaultTarget *net.TCPAddr
	s := strings.TrimSpace(c.DefaultTarget)
	if s == "" {
		return nil, errors.New("empty default target")
	}
	defaultTarget, err := net.ResolveTCPAddr("tcp", c.DefaultTarget)
	logger.Debug("Parse default target as tcp address", zap.Stringer("default_target", defaultTarget), zap.Error(err))

	if err != nil {
		defaultTargetIP, err := net.ResolveIPAddr("ip", c.DefaultTarget)
		logger.Debug("Parse default target as ip address", zap.Stringer("default_target", defaultTarget), zap.Error(err))
		if err != nil {
			logger.Error("Error parse default target address")
			return nil, err
		}
		defaultTarget = &net.TCPAddr{IP: defaultTargetIP.IP, Port: defaultHTTPPort}
	}

	if len(defaultTarget.IP) == 0 {
		logger.Info("Create same ip director", zap.Int("port", defaultTarget.Port))
		return NewDirectorSameIP(defaultTarget.Port), nil
	}

	logger.Info("Create host ip director", zap.Int("port", defaultTarget.Port))
	return NewDirectorHost(defaultTarget.String()), nil
}

//can return nil,nil
func (c *Config) getHeadersDirector(ctx context.Context) (Director, error) {
	logger := zc.L(ctx)

	if len(c.Headers) == 0 {
		return nil, nil
	}

	m := make(map[string]string)

	for _, line := range c.Headers {
		line = strings.TrimSpace(line)
		lineParts := strings.SplitN(line, ":", 2)
		if len(lineParts) != 2 {
			logger.Error("Can't split header line to parts", zap.String("line", line))
			return nil, errors.New("can't parse headers proxy config")
		}
		m[lineParts[0]] = lineParts[1]
	}

	logger.Info("Create headers director", zap.Any("headers", m))
	return NewDirectorSetHeaders(m), nil
}

// can return nil, nil
func (c *Config) getMapDirector(ctx context.Context) (Director, error) {
	logger := zc.L(ctx)
	if len(c.TargetMap) == 0 {
		return nil, nil
	}

	m := make(map[string]string)
	for _, line := range c.TargetMap {
		from, to, err := parseTCPMapPair(line)
		log.DebugError(logger, err, "Parse target map", zap.String("line", line),
			zap.String("from", from), zap.String("to", to))
		if err != nil {
			return nil, err
		}
		m[from] = to
	}

	logger.Info("Add target-map director", zap.Any("map", m))
	return NewDirectorDestMap(m), nil
}

func (c *Config) getSchemaDirector(ctx context.Context) (Director, error) {
	if c.HTTPSBackend {
		return NewSetSchemeDirector(ProtocolHTTPS), nil
	}
	return NewSetSchemeDirector(ProtocolHTTP), nil
}

func parseTCPMapPair(line string) (from, to string, err error) {
	line = strings.TrimSpace(line)
	lineParts := strings.Split(line, "-")
	if len(lineParts) != 2 {
		return "", "", errors.New("can't split tcp map to pair")
	}
	fromTCP, err := net.ResolveTCPAddr("tcp", lineParts[0])
	if err != nil {
		return "", "", fmt.Errorf("from addr can't resolve: %v", err.Error())
	}
	if len(fromTCP.IP) == 0 {
		return "", "", errors.New("from addr has no ip")
	}
	toTCP, err := net.ResolveTCPAddr("tcp", lineParts[1])
	if err != nil {
		return "", "", fmt.Errorf("to line can't resolve addr: %v", err.Error())
	}
	if len(toTCP.IP) == 0 {
		return "", "", errors.New("to addr has no ip")
	}

	from = fromTCP.String()
	to = toTCP.String()
	return from, to, nil
}
