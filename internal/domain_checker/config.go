//nolint:golint
package domain_checker

import (
	"context"
	"net"
	"regexp"

	zc "github.com/rekby/zapcontext"

	"github.com/rekby/lets-proxy2/internal/log"
	"go.uber.org/zap"
)

type Config struct {
	IPSelf      bool   `default:"true" comment:"Allow domain if it resolver for one of public IPs of this server."`
	IPWhiteList string `default:"" comment:"Allow domain if it resolver for one of the ips."`
	BlackList   string `default:"" comment:"Regexp in golang syntax of blacklisted domain for issue certificate.\nThis list overrided by whitelist."`
	WhiteList   string `default:"" comment:"Regexp in golang syntax of whitelist domains for issue certificate.\nWhitelist need for allow part of domains, which excluded by blacklist.\n"`
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

	var ipCheckers Any

	if c.IPSelf {
		selfPublicIpList := NewIPList(ctx, CreateGetSelfPublicBinded(net.InterfaceAddrs))
		selfPublicIpList.StartAutoRenew()
		ipCheckers = append(ipCheckers, selfPublicIpList)
	}

	if c.IPWhiteList != "" {
		ips, err := ParseIPs(ctx, c.IPWhiteList)
		log.DebugError(logger, err, "Parse ip whitelist")
		if err != nil {
			return nil, err
		}
		whiteIpList := NewIPList(ctx, func(ctx context.Context) ([]net.IP, error) {
			return ips, nil
		})
		// ipList.StartAutoRenew() - doesn't need renew, because list static
		ipCheckers = append(ipCheckers, whiteIpList)
	}

	// If no ip checks - allow domain without ip check
	// If have one or more ip checks - allow
	if len(ipCheckers) == 0 {
		ipCheckers = NewAny(True{})
	}

	res := NewAll(listCheckers, ipCheckers)
	return res, nil
}
