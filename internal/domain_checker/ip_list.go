//nolint:golint
package domain_checker

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	zc "github.com/rekby/zapcontext"

	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/log"
)

var (
	nonPublicIPNetworks = []net.IPNet{
		// list networks from https://en.wikipedia.org/wiki/Reserved_IP_addresses
		mustParseNet("0.0.0.0/8"),
		mustParseNet("10.0.0.0/8"),
		mustParseNet("100.64.0.0/10"),
		mustParseNet("127.0.0.0/8"),
		mustParseNet("169.254.0.0/16"),
		mustParseNet("172.16.0.0/12"),
		mustParseNet("192.0.0.0/24"),
		mustParseNet("192.0.2.0/24"),
		mustParseNet("192.88.99.0/24"), // Is global Anycast addresses, can't handle TCP on this
		mustParseNet("192.168.0.0/16"),
		mustParseNet("198.18.0.0/15"),
		mustParseNet("198.51.100.0/24"),
		mustParseNet("203.0.113.0/24"),
		mustParseNet("224.0.0.0/4"),
		mustParseNet("240.0.0.0/4"),
		mustParseNet("255.255.255.255/32"),
		//mustParseNet("::/0"),
		mustParseNet("::/128"),
		mustParseNet("::1/128"),
		//mustParseNet("::ffff:0:0/96"),
		mustParseNet("::ffff:0:0:0/96"),
		//mustParseNet("64:ff9b::/96"),
		mustParseNet("100::/64"),
		//mustParseNet("2001::/32"),
		mustParseNet("2001:20::/28"),
		mustParseNet("2001:db8::/32"),
		mustParseNet("2002::/16"),
		mustParseNet("fc00::/7"),
		mustParseNet("fe80::/10"),
		mustParseNet("ff00::/8"),
	}

	defaultResolver Resolver = net.DefaultResolver
)

func SetDefaultResolver(resolver Resolver) {
	defaultResolver = resolver
}

type IPList struct {
	Addresses          AllowedIPAddresses
	Resolver           Resolver
	AutoUpdateInterval time.Duration // Set zero for disable autorenew.

	ctx     context.Context
	mu      sync.RWMutex
	ips     []net.IP
	started bool
}

type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

type AllowedIPAddresses func(ctx context.Context) ([]net.IP, error)

// After create can change settings fields, than can call StartAutoRenew
// struct fields MUST NOT changes after call StartAutoRenew or concurrency with usage.
func NewIPList(ctx context.Context, addresses AllowedIPAddresses) *IPList {
	res := &IPList{
		ctx:                ctx,
		Addresses:          addresses,
		Resolver:           defaultResolver,
		AutoUpdateInterval: time.Hour,
	}
	res.updateIPs()
	return res
}

func (s *IPList) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	if s.ctx.Err() != nil {
		return false, errors.New("iplist main context canceled")
	}

	logger := zc.L(ctx)
	ips, err := s.Resolver.LookupIPAddr(ctx, domain)
	log.DebugInfo(logger, err, "Resolve domain ip addresses", zap.Any("ips", ips))
	if err != nil {
		return false, err
	}

	if len(ips) == 0 {
		logger.Info("Doesn't allow domain without ip address")
		return false, errors.New("domain has no ip address")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

hostIP:
	for _, ip := range ips {
		for _, bindedIP := range s.ips {
			if ip.IP.Equal(bindedIP) {
				continue hostIP
			}
		}
		ip := ip
		logger.Debug("Non self or private ip", zap.Stringer("checked_ip", &ip))
		return false, nil
	}
	return true, nil
}

// Can called most once - for autorenew internal ips
func (s *IPList) StartAutoRenew() {
	s.updateIPs()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		zc.L(s.ctx).DPanic("Double started self public ip")
	}
	s.started = true
	if s.AutoUpdateInterval > 0 {
		// handlepanic: in updateIPsByTimer
		go s.updateIPsByTimer()
	}
}

func (s *IPList) updateIPs() {
	ips, err := s.Addresses(s.ctx)
	log.DebugDPanicCtx(s.ctx, err, "Got ips while auto update", zap.Any("ips", ips))
	if err != nil {
		return
	}

	s.mu.Lock()
	s.ips = ips
	s.mu.Unlock()
}

func (s *IPList) updateIPsByTimer() {
	contextDone := s.ctx.Done()
	ticker := time.NewTicker(s.AutoUpdateInterval)
	defer ticker.Stop()

	logger := zc.L(s.ctx)

	for {
		select {
		case <-contextDone:
			return
		case <-ticker.C:
			func() {
				defer log.HandlePanic(logger)

				s.updateIPs()
			}()
		}
	}
}

type InterfacesAddrFunc func() ([]net.Addr, error)

func ParseIPList(ctx context.Context, s, sep string) ([]net.IP, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	logger := zc.L(ctx)
	parts := strings.Split(s, sep)
	var parsed = make([]net.IP, 0, len(parts))
	for _, stringIP := range parts {
		stringIP = strings.TrimSpace(stringIP)
		if stringIP == "" {
			continue
		}
		ip := net.ParseIP(stringIP)
		if ip == nil {
			logger.Error("Can't parse ip", zap.String("ip", stringIP))
			return nil, errors.New("can't parse ip")
		}
		logger.Debug("Parse ip", zap.Stringer("ip", ip))
		parsed = append(parsed, ip)
	}

	res := truncatedCopyIPs(parsed)
	return res, nil
}

func mustParseNet(s string) net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if ipnet == nil || err != nil {
		panic(err)
	}
	return *ipnet
}

func isPublicIP(ip net.IP) bool {
	if len(ip) == 0 {
		return false
	}

	for _, ipNet := range nonPublicIPNetworks {
		if ipNet.Contains(ip) {
			return false
		}
	}
	return true
}

// return copy of ips, with cap truncated for len
// return nil if len(ips) == 0
func truncatedCopyIPs(ips []net.IP) []net.IP {
	if len(ips) == 0 {
		return nil
	}

	var res = make([]net.IP, len(ips))

	copy(res, ips)

	return res
}

func firstError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
