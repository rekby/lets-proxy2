//nolint:golint
package domain_checker

import (
	"context"
	"net"
	"regexp"
	"strings"

	"golang.org/x/xerrors"

	"github.com/pkg/errors"

	"github.com/rekby/lets-proxy2/internal/dns"

	zc "github.com/rekby/zapcontext"

	"github.com/rekby/lets-proxy2/internal/log"
	"go.uber.org/zap"
)

type Config struct {
	IPSelf                    bool
	IPSelfDetectMethod        string
	IPSelfExternalDetectorURL string
	IPWhiteList               string
	BlackList                 string
	WhiteList                 string
	Resolver                  string
}

func (c *Config) CreateDomainChecker(ctx context.Context) (DomainChecker, error) {
	logger := zc.L(ctx)

	var listCheckers DomainChecker = True{}

	if c.BlackList != "" {
		r, err := regexp.Compile(c.BlackList)
		log.InfoError(logger, err, "Compile blacklist regexp", zap.String("regexp", c.BlackList))
		if err != nil {
			return nil, err
		}
		listCheckers = NewAll(NewNot(NewRegexp(r)), listCheckers)
	}

	if c.WhiteList != "" {
		r, err := regexp.Compile(c.WhiteList)
		log.InfoError(logger, err, "Compile whitelist regexp", zap.String("regexp", c.WhiteList))
		if err != nil {
			return nil, err
		}
		listCheckers = NewAny(listCheckers, NewRegexp(r))
	}

	resolver, err := c.createResolver(logger)
	if err != nil {
		log.DebugError(logger, err, "Create resolver")
		return nil, err
	}
	SetDefaultResolver(resolver)

	var ipCheckers Any

	if c.IPSelf {
		selfIPChecker, err := NewSelfIPChecker(ctx, c)
		if err != nil {
			return nil, xerrors.Errorf("create self ip checkers: %w", err)
		}
		ipCheckers = append(ipCheckers, selfIPChecker)
	}

	if c.IPWhiteList != "" {
		ips, err := ParseIPList(ctx, c.IPWhiteList, ",")
		log.DebugError(logger, err, "Parse ip whitelist")
		if err != nil {
			return nil, err
		}
		whiteIPList := NewIPList(ctx, func(ctx context.Context) ([]net.IP, error) {
			return ips, nil
		})
		// ipList.StartAutoRenew() - doesn't need renew, because list static
		ipCheckers = append(ipCheckers, whiteIPList)
	}

	// If no ip checks - allow domain without ip check
	// If have one or more ip checks - allow
	if len(ipCheckers) == 0 {
		ipCheckers = NewAny(True{})
	}

	res := NewAll(listCheckers, ipCheckers)
	return res, nil
}

func (c *Config) createResolver(logger *zap.Logger) (Resolver, error) {
	var resolver Resolver
	if strings.TrimSpace(c.Resolver) == "" {
		resolver = net.DefaultResolver
	} else {
		stringAddresses := strings.Split(c.Resolver, ",")
		var resolvers = make([]dns.ResolverInterface, 0, len(stringAddresses))
		for _, addr := range stringAddresses {
			addr = strings.TrimSpace(addr)
			if addr == "" {
				continue
			}
			tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				logger.Error("Can't resolve dns server address string", zap.String("addr", addr), zap.Error(err))
				return nil, err
			}
			if len(tcpAddr.IP) == 0 {
				logger.Error("Can't resolve dns server address ip - it is empty.", zap.String("addr", addr))
				return nil, errors.New("empty ip address")
			}
			if tcpAddr.Port == 0 {
				tcpAddr.Port = 53 // default dns port
			}
			tcpAddrString := tcpAddr.String()
			resolvers = append(resolvers, dns.NewResolver(tcpAddrString))
		}
		resolver = dns.NewParallel(resolvers...)
	}
	return resolver, nil
}
