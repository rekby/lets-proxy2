//nolint:golint
package domain_checker

import (
	"context"
	"regexp"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rekby/lets-proxy2/internal/log"
)

type True struct{}

func (True) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	log.DebugCtx(ctx, "Allow by 'true' rule")
	return true, nil
}

type False struct{}

func (False) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	log.InfoCtx(ctx, "Deny by 'false' rule")
	return false, nil
}

type Not struct {
	origin DomainChecker
}

func (n Not) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	subRuleRes, err := n.origin.IsDomainAllowed(ctx, domain)
	res := !subRuleRes
	if err == nil {
		logLevel := zapcore.DebugLevel
		if !res {
			logLevel = zapcore.InfoLevel
		}
		log.LevelParamCtx(ctx, logLevel, "'Not' filter (details in debug log)", zap.Bool("result", res))
		return res, nil
	}
	return false, err
}

type Regexp regexp.Regexp

func (r *Regexp) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	reg := (*regexp.Regexp)(r)
	result := reg.MatchString(domain)
	logLevel := zapcore.DebugLevel
	if !result {
		logLevel = zapcore.InfoLevel
	}
	log.LevelParamCtx(ctx, logLevel, "Check if domain allowed by regexp", zap.String("regexp", reg.String()), zap.Bool("result", result))
	return result, nil
}

func NewRegexp(r *regexp.Regexp) *Regexp {
	return (*Regexp)(r)
}
