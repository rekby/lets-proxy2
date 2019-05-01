//nolint:golint
package domain_checker

import (
	"context"
	"net"
	"sync"
	"time"

	zc "github.com/rekby/zapcontext"

	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/log"
)

var (
	nonPublicIpNetworks = []net.IPNet{
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
)

type SelfPublicIP struct {
	Addresses          InterfacesAddrFunc
	Resolver           Resolver
	AutoUpdateInterval time.Duration

	ctx     context.Context
	mu      sync.RWMutex
	ips     []net.IP
	started bool
}

type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

type InterfacesAddrFunc func() ([]net.Addr, error)

// After create can change settings fields, than must call Start()
// Fields must not change after start.
func NewSelfIP(ctx context.Context) *SelfPublicIP {
	res := &SelfPublicIP{
		ctx:                ctx,
		Addresses:          net.InterfaceAddrs,
		Resolver:           net.DefaultResolver,
		AutoUpdateInterval: time.Hour,
	}
	res.updateIPs()
	return res
}

func (s *SelfPublicIP) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	logger := zc.L(ctx)
	s.mu.RLock()
	defer s.mu.RUnlock()

	ips, err := s.Resolver.LookupIPAddr(ctx, domain)
	log.DebugInfo(logger, err, "Resolve domain ip addresses", zap.Any("ips", ips))
	if err != nil {
		return false, err
	}

hostIP:
	for _, ip := range ips {
		for _, bindedIp := range s.ips {
			if ip.IP.Equal(bindedIp) {
				continue hostIP
			}
		}
		ip := ip
		logger.Debug("Non self or private ip", zap.Stringer("checked_ip", &ip))
		return false, nil
	}
	return true, nil
}

func (s *SelfPublicIP) Start() {
	s.updateIPs()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		zc.L(s.ctx).DPanic("Double started self public ip")
	}
	s.started = true

	go s.updateIPsByTimer()
}

func (s *SelfPublicIP) updateIPs() {
	ips := getSelfPublicIPs(s.ctx, s.Addresses)

	s.mu.Lock()
	s.ips = ips
	s.mu.Unlock()
}

func (s *SelfPublicIP) updateIPsByTimer() {
	contextDone := s.ctx.Done()
	ticker := time.NewTicker(s.AutoUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-contextDone:
			return
		case <-ticker.C:
			s.updateIPs()
		}
	}
}

func getSelfPublicIPs(ctx context.Context, interfacesAddr InterfacesAddrFunc) []net.IP {
	logger := zc.L(ctx)
	binded, err := interfacesAddr()
	log.DebugDPanic(logger, err, "Get system addresses", zap.Any("addresses", binded))

	var parsed = make([]net.IP, 0, len(binded))
	for _, addr := range binded {
		ip, _, err := net.ParseCIDR(addr.String())
		log.DebugDPanic(logger, err, "Parse ip from interface", zap.Any("ip", ip),
			zap.Stringer("original_ip", addr))
		if ip == nil {
			continue
		}
		logger.Debug("Parse ip", zap.Stringer("ip", ip))
		parsed = append(parsed, ip)
	}

	var public = make([]net.IP, 0, len(parsed))
	for _, ip := range parsed {
		if isPublicIp(ip) {
			public = append(public, ip)
		}
	}

	// Truncate pre_allocated_memory
	var res = make([]net.IP, len(public))
	copy(res, public)
	return res
}

func mustParseNet(s string) net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if ipnet == nil || err != nil {
		panic(err)
	}
	return *ipnet
}

func isPublicIp(ip net.IP) bool {
	if len(ip) == 0 {
		return false
	}

	for _, ipNet := range nonPublicIpNetworks {
		if ipNet.Contains(ip) {
			return false
		}
	}
	return true
}
