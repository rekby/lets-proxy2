package dns

import (
	"context"
	"github.com/rekby/lets-proxy2/internal/log"
	zc "github.com/rekby/zapcontext"
	"net"
	"sync"
)

type ResolverInterface interface {
	// LookupIPAddr return ip addresses of domain. It MUST finish work when context canceled
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

type Parallel []ResolverInterface

// NewParallel return parallel resolver
func NewParallel(resolvers ...ResolverInterface) Parallel {
	state := make(Parallel, len(resolvers))
	copy(state, resolvers)
	return state
}

// LookupIPAddr return ip addresses of host, used underly resolvers in parallel
// If any of resolvers return ips - return sum array of the ips (may duplicated)
// If all resolvers return error - return any of they errors
func (p Parallel) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	logger := zc.L(ctx)

	switch len(p) {
	case 0:
		return nil, nil
	case 1:
		return p[0].LookupIPAddr(ctx, host)
	default:
		// pass
	}

	var ips = make([][]net.IPAddr, len(p))
	var errs = make([]error, len(p))

	var wg sync.WaitGroup
	wg.Add(len(p))     // nolint:wsl
	for i := range p { // nolint:wsl
		go func(i int) {
			defer wg.Done()
			defer log.HandlePanic(logger)

			ips[i], errs[i] = p[i].LookupIPAddr(ctx, host)
		}(i)
	}

	wg.Wait()

	resLen := 0
	var err error
	for i := range ips {
		resLen += len(ips[i])
		if errs[i] != nil {
			err = errs[i]
		}
	}

	if resLen == 0 {
		return nil, err
	}

	resIps := make([]net.IPAddr, 0, resLen)
	for i := range ips {
		resIps = append(resIps, ips[i]...)
	}
	return resIps, nil
}
