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
	SelfIP    bool   `default:"true",comment:"Enable check domain for it resolved to IPs of this server only.\n"`
	BlackList string `default:"",comment:"Regexp in golang syntax of blacklisted domain for issue certificate.\nThis list overrided by whitelist."`
	WhiteList string `default:"",comment:"Regexp in golang syntax of whitelist domains for issue certificate.\nWhitelist need for allow part of domains, which excluded by blacklist.\n"`
}

func (c *Config) CreateDomainChecker(ctx context.Context) (DomainChecker, error) {
	logger := zc.L(ctx)
	var res DomainChecker = True{}

	if c.BlackList != "" {
		r, err := regexp.Compile(c.BlackList)
		log.InfoError(logger, err, "Compile blacklist regexp", zap.String("regexp", c.BlackList))
		if err != nil {
			return nil, err
		}
		res = NewAll(NewNot(NewRegexp(r)), res)
	}

	if c.WhiteList != "" {
		r, err := regexp.Compile(c.WhiteList)
		log.InfoError(logger, err, "Compile whitelist regexp", zap.String("regexp", c.WhiteList))
		if err != nil {
			return nil, err
		}
		res = NewAny(NewRegexp(r), res)
	}

	if c.SelfIP {
		ipList := NewIPList(ctx, CreateGetSelfPublicBinded(net.InterfaceAddrs))
		ipList.StartAutoRenew()
		res = NewAll(res, ipList)
	}
	return res, nil
}
