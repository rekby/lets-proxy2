package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rekby/lets-proxy2/internal/contextlabel"

	"github.com/rekby/lets-proxy2/internal/log"

	zc "github.com/rekby/zapcontext"

	"go.uber.org/zap"
)

const (
	ConnectionID = "{{CONNECTION_ID}}"
	HTTPProto    = "{{HTTP_PROTO}}"
	SourceIP     = "{{SOURCE_IP}}"
	SourcePort   = "{{SOURCE_PORT}}"
	SourceIPPort = "{{SOURCE_IP}}:{{SOURCE_PORT}}"
)

const (
	ProtocolHTTP  = "http"
	ProtocolHTTPS = "https"
)

type DirectorChain []Director

func (c DirectorChain) Director(request *http.Request) error {
	for _, d := range c {
		err := d.Director(request)
		if err != nil {
			return err
		}
	}
	return nil
}

// skip nil directors
func NewDirectorChain(directors ...Director) DirectorChain {
	cnt := 0

	for _, item := range directors {
		if item != nil {
			cnt++
		}
	}

	ownDirectors := make(DirectorChain, 0, cnt)

	for _, item := range directors {
		if item != nil {
			ownDirectors = append(ownDirectors, item)
		}
	}

	return ownDirectors
}

type DirectorSameIP struct {
	Port string
}

func NewDirectorSameIP(port int) DirectorSameIP {
	return DirectorSameIP{strconv.Itoa(port)}
}

func (s DirectorSameIP) Director(request *http.Request) error {
	localAddr := request.Context().Value(http.LocalAddrContextKey).(*net.TCPAddr)
	if request.URL == nil {
		request.URL = &url.URL{}
	}
	request.URL.Host = localAddr.IP.String() + ":" + s.Port
	zc.L(request.Context()).Debug("Set target as same ip",
		zap.Stringer("local_addr", localAddr), zap.String("dest_host", request.Host))
	return nil
}

type DirectorDestMap map[string]string

func (d DirectorDestMap) Director(request *http.Request) error {
	ctx := request.Context()

	type Stringer interface {
		String() string
	}

	localAddr := ctx.Value(http.LocalAddrContextKey).(Stringer).String()
	var dest string
	var ok bool
	if dest, ok = d[localAddr]; !ok {
		zc.L(ctx).Debug("Map director no matches, skip.")
		return nil
	}

	if request.URL == nil {
		request.URL = &url.URL{}
	}
	request.URL.Host = dest
	zc.L(ctx).Debug("Map director set dest", zap.String("host", request.URL.Host))
	return nil
}

func NewDirectorDestMap(m map[string]string) DirectorDestMap {
	res := make(DirectorDestMap, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

type DirectorHost string

func (d DirectorHost) Director(request *http.Request) error {
	if request.URL == nil {
		request.URL = &url.URL{}
	}
	request.URL.Host = string(d)
	return nil
}

func NewDirectorHost(host string) DirectorHost {
	return DirectorHost(host)
}

type DirectorSetHeaders map[string]string

func NewDirectorSetHeaders(m map[string]string) DirectorSetHeaders {
	res := make(DirectorSetHeaders, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

func (h DirectorSetHeaders) Director(request *http.Request) error {
	ctx := request.Context()
	host, port, errHostPort := net.SplitHostPort(request.RemoteAddr)
	log.DebugDPanicCtx(ctx, errHostPort, "Parse remote addr for headers", zap.String("host", host), zap.String("port", port))

	for name, headerVal := range h {
		var value string

		switch headerVal {
		case ConnectionID:
			value = request.Context().Value(contextlabel.ConnectionID).(string)
		case HTTPProto:
			if tls, ok := ctx.Value(contextlabel.TLSConnection).(bool); ok {
				if tls {
					value = ProtocolHTTPS
				} else {
					value = ProtocolHTTP
				}
			} else {
				value = "error protocol detection"
			}
		case SourceIP:
			value = host
		case SourceIPPort:
			value = host + ":" + port
		case SourcePort:
			value = port
		default:
			value = headerVal
		}

		if request.Header == nil {
			request.Header = make(http.Header)
		}

		request.Header.Set(name, value)
	}
	return nil
}

type HTTPHeader struct {
	Name  string
	Value string
}
type HTTPHeaders []HTTPHeader
type NetHeaders struct {
	IPNet   net.IPNet
	Headers HTTPHeaders
}
type DirectorSetHeadersByIP []NetHeaders

func NewDirectorSetHeadersByIP(m map[string]HTTPHeaders) (DirectorSetHeadersByIP, error) {
	res := make(DirectorSetHeadersByIP, 0, len(m))
	for k, v := range m {
		_, subnet, err := net.ParseCIDR(k)
		if err != nil {
			return nil, fmt.Errorf("can't parse CIDR: %v %w", k, err)
		}

		res = append(res, NetHeaders{
			IPNet:   *subnet,
			Headers: v,
		})
	}

	return sortByIPNet(res), nil
}

// sortByIPNet sorts by CIDR using quicksort algorithm.
func sortByIPNet(d DirectorSetHeadersByIP) DirectorSetHeadersByIP {
	ipv4 := make(DirectorSetHeadersByIP, 0, len(d))
	ipv6 := make(DirectorSetHeadersByIP, 0, len(d))
	for _, item := range d {
		if item.IPNet.IP.To4() != nil {
			ipv4 = append(ipv4, item)
		} else {
			ipv6 = append(ipv6, item)
		}
	}

	ipv4 = quickSortByIPNet(ipv4)
	ipv6 = quickSortByIPNet(ipv6)

	return append(ipv4, ipv6...)
}

// quickSortByIPNet sorts by CIDR using quicksort algorithm.
// The result is sorted by IPNet.
// example:
//
//	IN -> DirectorSetHeadersByIP{
//				{IPNet: net.IPNet{IP: net.ParseIP("192.168.88.0"), Mask: net.CIDRMask(24, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("192.0.0.0"), Mask: net.CIDRMask(8, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(16, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("192.168.99.0"), Mask: net.CIDRMask(24, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("172.0.0.0"), Mask: net.CIDRMask(8, 32)}},
//			},
//
//	OUT <- DirectorSetHeadersByIP{
//				{IPNet: net.IPNet{IP: net.ParseIP("192.0.0.0"), Mask: net.CIDRMask(8, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("192.168.88.0"), Mask: net.CIDRMask(24, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("192.168.99.0"), Mask: net.CIDRMask(24, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("172.0.0.0"), Mask: net.CIDRMask(8, 32)}},
//				{IPNet: net.IPNet{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(16, 32)}},
//			},
func quickSortByIPNet(d DirectorSetHeadersByIP) DirectorSetHeadersByIP {
	if len(d) <= 1 {
		return d
	}

	mid := len(d) / 2
	left := d[:mid]
	right := d[mid:]

	left = quickSortByIPNet(left)
	right = quickSortByIPNet(right)

	return mergeByIPNet(left, right)
}

// mergeByIPNet merges two sorted arrays with CIDRs.
// The result is sorted by IPNet.
func mergeByIPNet(left, right DirectorSetHeadersByIP) DirectorSetHeadersByIP {
	res := make(DirectorSetHeadersByIP, 0, len(left)+len(right))
	for len(left) > 0 || len(right) > 0 {
		if len(left) > 0 && len(right) > 0 {
			if left[0].IPNet.Contains(right[0].IPNet.IP) {
				res = append(res, left[0])
				left = left[1:]
			} else if right[0].IPNet.Contains(left[0].IPNet.IP) {
				res = append(res, right[0])
				right = right[1:]
			} else {
				res = append(res, left[0], right[0])
				left = left[1:]
				right = right[1:]
			}
		} else if len(left) > 0 {
			res = append(res, left[0])
			left = left[1:]
		} else if len(right) > 0 {
			res = append(res, right[0])
			right = right[1:]
		}
	}
	return res
}

func (h DirectorSetHeadersByIP) Director(request *http.Request) error {
	if request == nil {
		return fmt.Errorf("request is nil")
	}

	ctx := request.Context()
	host, port, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		zc.L(ctx).Debug("Split host port error", zap.Error(err), zap.String("host", host),
			zap.String("port", port))
	}

	ip := net.ParseIP(host)

	for _, ipHeaders := range h {
		if !ipHeaders.IPNet.Contains(ip) {
			continue
		}

		if request.Header == nil {
			request.Header = make(http.Header)
		}

		for _, header := range ipHeaders.Headers {
			request.Header.Set(header.Name, header.Value)
		}
	}

	return nil
}

type DirectorSetScheme string

func (d DirectorSetScheme) Director(req *http.Request) error {
	if req.URL == nil {
		req.URL = &url.URL{}
	}
	req.URL.Scheme = string(d)
	return nil
}

func NewSetSchemeDirector(scheme string) DirectorSetScheme {
	return DirectorSetScheme(scheme)
}
