package dns

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	zc "github.com/rekby/zapcontext"

	"github.com/rekby/lets-proxy2/internal/log"
	"go.uber.org/zap"

	mdns "github.com/miekg/dns"
)

const maxDNSRecursion = 10

var (
	errTruncatedResponse = errors.New("truncated answer")
	errPanic             = errors.New("panic")
)

type mDNSClient interface {
	Exchange(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error)
}

// Resolve IPs for A and AAAA records of domains
// it use direct dns query without cache
type Resolver struct {
	udp                 mDNSClient
	tcp                 mDNSClient
	server              string
	maxDNSRecursionDeep int
	lookupWithClient    func(ctx context.Context, host string, server string, recordType uint16, recursion int, client mDNSClient) ([]net.IPAddr, error)
}

// NewResolver return direct dns resolver
func NewResolver(dnsServer string) *Resolver {
	return &Resolver{
		udp:                 &mdns.Client{Net: "udp"},
		tcp:                 &mdns.Client{Net: "tcp"},
		server:              dnsServer,
		maxDNSRecursionDeep: maxDNSRecursion,
		lookupWithClient:    lookupWithClient,
	}
}

func (r *Resolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	logger := zc.L(ctx).With(zap.String("dns_server", r.server))
	ctx = zc.WithLogger(ctx, logger)
	if !strings.HasSuffix(host, ".") {
		host += "."
	}

	var wg sync.WaitGroup

	var ipAddrA, ipAddrAAAA []net.IPAddr
	var errA, errAAAA error

	wg.Add(1)
	go func() { //nolint:wsl
		defer wg.Done()
		defer log.HandlePanic(logger)

		errA = errPanic
		ipAddrA, errA = r.lookup(ctx, host, mdns.TypeA)
	}()

	wg.Add(1)
	go func() { //nolint:wsl
		defer wg.Done()
		defer log.HandlePanic(logger)

		errAAAA = errPanic
		ipAddrAAAA, errAAAA = r.lookup(ctx, host, mdns.TypeAAAA)
	}()

	wg.Wait()

	var resultErr error
	if errAAAA != nil {
		resultErr = errAAAA
	}
	if errA != nil {
		resultErr = errA
	}

	log.DebugError(logger, resultErr, "Host lookup", zap.NamedError("errA", errA),
		zap.NamedError("errAAAA", errAAAA), zap.Any("ipAddrA", ipAddrA),
		zap.Any("ipAddrAAAA", ipAddrAAAA))

	if resultErr != nil {
		return nil, resultErr
	}
	resultIPs := make([]net.IPAddr, len(ipAddrA)+len(ipAddrAAAA))
	copy(resultIPs, ipAddrA)
	copy(resultIPs[len(ipAddrA):], ipAddrAAAA)
	return resultIPs, nil
}

func (r *Resolver) lookup(ctx context.Context, host string, recordType uint16) ([]net.IPAddr, error) {
	res, err := r.lookupWithClient(ctx, host, r.server, recordType, r.maxDNSRecursionDeep, r.udp)
	if err == errTruncatedResponse {
		zc.L(ctx).Debug("fallback to tcp request")
		res, err = r.lookupWithClient(ctx, host, r.server, recordType, r.maxDNSRecursionDeep, r.tcp)
	}
	return res, err
}

//nolint:funlen
func lookupWithClient(ctx context.Context, host string, server string, recordType uint16, recursion int, client mDNSClient) (ipResults []net.IPAddr, err error) {
	logger := zc.L(ctx)

	if recursion <= 0 {
		logger.Error("Max recursion while resolve domain")
		return nil, errors.New("max recursion while resolve domain")
	}

	if ctx.Err() != nil {
		logger.Debug("Context canceled")
		return nil, ctx.Err()
	}

	defer func() {
		log.DebugError(logger, err, "Resolved ips", zap.Any("ipResults", ipResults),
			zap.Uint16("record_type", recordType))
	}()

	var msdID uint16
	for msdID == 0 {
		msdID = mdns.Id()
	}

	msg := mdns.Msg{
		MsgHdr: mdns.MsgHdr{
			Id: msdID,
		},
		Question: []mdns.Question{
			{Name: host, Qclass: mdns.ClassINET, Qtype: recordType},
		},
	}
	msg.RecursionDesired = true
	exchangeCompleted := make(chan struct {
		answer *mdns.Msg
		err    error
	}, 1)

	go func() { // nolint:wsl
		defer close(exchangeCompleted)
		defer log.HandlePanic(logger)

		dnsAnswer, _, dnsErr := client.Exchange(&msg, server)
		exchangeCompleted <- struct {
			answer *mdns.Msg
			err    error
		}{answer: dnsAnswer, err: dnsErr}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case answer := <-exchangeCompleted:
		dnsAnswer := answer.answer
		if dnsAnswer != nil && dnsAnswer.Truncated {
			return nil, errTruncatedResponse
		}

		if answer.err != nil {
			return nil, answer.err
		}

		var resIPs []net.IPAddr
		for _, r := range dnsAnswer.Answer {
			rType := r.Header().Rrtype

			switch {
			case rType == mdns.TypeA && rType == recordType:
				resIPs = append(resIPs, net.IPAddr{IP: r.(*mdns.A).A})
			case rType == mdns.TypeAAAA && rType == recordType:
				resIPs = append(resIPs, net.IPAddr{IP: r.(*mdns.AAAA).AAAA})
			case rType == mdns.TypeCNAME:
				cname := r.(*mdns.CNAME)
				zc.L(ctx).Debug("Receive CNAME record for domain.", zap.String("target", cname.Target))
				return lookupWithClient(ctx, cname.Target, server, recordType, recursion-1, client)
			default:
				// pass
			}
		}
		return resIPs, nil
	}
}
