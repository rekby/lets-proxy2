package domain_checker

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"golang.org/x/xerrors"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/rekby/lets-proxy2/internal/log"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

const getIPByExternalRequestTimeout = time.Second * 10

//nolint:funlen
func NewSelfIPChecker(ctx context.Context, config *Config) (DomainChecker, error) {
	logger := zc.L(ctx)
	logger.Info("Create ip checker", zap.String("IPSelfDetectMethod", config.IPSelfDetectMethod))

	publicIPDetected := false
	var ipSources []DomainChecker
	addBinded := func() {
		selfBindedFunc := SelfBindedPublicIPs(net.InterfaceAddrs)
		selfPublicIPList := NewIPList(ctx, selfBindedFunc)
		selfPublicIPList.StartAutoRenew()
		ipSources = append(ipSources, selfPublicIPList)
		logger.Debug("Add selfPublicIPList")

		if !publicIPDetected {
			publicIPs, _ := selfBindedFunc(ctx)
			publicIPDetected = len(publicIPs) > 0
		}
	}
	addAws := func() {
		awsList := NewIPList(ctx, AWSPublicIPs())
		awsList.StartAutoRenew()
		ipSources = append(ipSources, awsList)
		logger.Debug("Add awsList")
		if !publicIPDetected {
			publicIPs, _ := AWSPublicIPs()(ctx)
			publicIPDetected = len(publicIPs) > 0
		}
	}
	addExternal := func() {
		externalList := NewIPList(ctx, GetIPByExternalRequest(config.IPSelfExternalDetectorURL))
		externalList.StartAutoRenew()
		ipSources = append(ipSources, externalList)
		logger.Debug("Add external")
	}

	switch config.IPSelfDetectMethod {
	case "auto":
		addBinded()
		awsAvailable := AWSMetadataAvailable()
		logger.Info("Check aws metadata available", zap.Bool("available", awsAvailable))
		if awsAvailable {
			addAws()
		}
		if !publicIPDetected {
			addExternal()
		}
	case "aws", "yandex":
		awsAvailable := AWSMetadataAvailable()
		logger.Info("Check aws metadata available", zap.Bool("available", awsAvailable))
		if !awsAvailable {
			return nil, xerrors.Errorf("Aws metadata doesn't available")
		}
		addAws()
	case "bind":
		addBinded()
	case "external":
		addExternal()
	default:
		return nil, xerrors.Errorf("Unknown IPSelfDetectMethod: '%v'", config.IPSelfDetectMethod)
	}
	return NewAny(ipSources...), nil
}

func AWSMetadataAvailable() bool {
	s, err := session.NewSession()
	if err != nil {
		return false
	}
	client := ec2metadata.New(s)
	return client.Available()
}

func AWSPublicIPs() AllowedIPAddresses {
	s, err := session.NewSession()
	if err != nil {
		panic("shoud be detected early by AWSMetadataAvailable")
	}
	client := ec2metadata.New(s)
	return awsPublicIPs(client.GetMetadata)
}

func awsPublicIPs(getMetadata func(p string) (string, error)) AllowedIPAddresses {
	return func(ctx context.Context) (ips []net.IP, err error) {
		logger := zc.L(ctx)

		var errors []error

		publicIPv4S, err := getMetadata("public-ipv4")
		log.DebugError(logger, err, "Get ipv4 address", zap.String("ipv4", publicIPv4S), zap.Error(err))
		errors = append(errors, err)
		ips, err = ParseIPList(ctx, publicIPv4S, "\n")
		log.DebugError(logger, err, "Parse IPv4", zap.Any("ipv4", ips))
		errors = append(errors, err)

		macs, err := getMetadata("network/interfaces/macs/")
		log.DebugError(logger, err, "Get mac address", zap.Error(err))
		errors = append(errors, err)

		if err == nil {
			macList := strings.Split(macs, "\n")
			for _, mac := range macList {
				mac = strings.TrimSuffix(mac, "/")
				metaString := "network/interfaces/macs/" + mac + "/ipv6s"
				ipv6S, err := getMetadata(metaString)
				log.DebugError(logger, err, "Get ipv6", zap.String("mac", mac), zap.String("ipv6", ipv6S))
				errors = append(errors, err)
				ipv6Addr, err := ParseIPList(ctx, ipv6S, "\n")
				ips = append(ips, ipv6Addr...)
				errors = append(errors, err)
			}
		}

		if len(ips) > 0 {
			return ips, nil
		}
		return nil, firstError(errors...)
	}
}

func GetIPByExternalRequest(url string) AllowedIPAddresses {
	return func(ctx context.Context) (ips []net.IP, err error) {
		return getIPByExternalRequest(ctx, url)
	}
}

func getIPByExternalRequest(ctx context.Context, url string) ([]net.IP, error) {
	logger := zc.L(ctx).With(zap.String("detector", url))

	fGetIP := func(network string) (net.IP, error) {
		client := http.Client{Transport: &http.Transport{
			DialContext: func(ctx context.Context, _supressNetwork, addr string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
		},
		}
		client.Timeout = getIPByExternalRequestTimeout
		resp, err := client.Get(url)
		if resp != nil && resp.Body != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		log.DebugError(logger, err, "Request to external ip detector")
		if err != nil {
			return nil, xerrors.Errorf("request to ip detector")
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		log.DebugError(logger, err, "Read response from ip detector")
		if err != nil {
			return nil, xerrors.Errorf("response from ip detector: %w", err)
		}
		ip := net.ParseIP(strings.TrimSpace(string(respBytes)))
		logrus.Debugf("Detected ip by '%v' (%v): %v", url, network, ip)
		if ip == nil {
			return nil, xerrors.Errorf("ip detector parse response")
		}
		return ip, nil
	}

	ipsFromDetector := make([]net.IP, 2)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	var errTCP4, errTCP6 error
	go func() {
		defer wg.Done()
		defer log.HandlePanic(logger)

		ipsFromDetector[0], errTCP4 = fGetIP("tcp4")
	}()

	go func() {
		defer wg.Done()
		defer log.HandlePanic(logger)

		ipsFromDetector[1], errTCP6 = fGetIP("tcp6")
	}()
	wg.Wait()

	res := make([]net.IP, 0, 2)
	for _, ip := range ipsFromDetector {
		if ip != nil {
			res = append(res, ip)
		}
	}
	if len(res) == 0 {
		return nil, firstError(errTCP4, errTCP6)
	}
	return truncatedCopyIPs(res), nil
}

func SelfBindedPublicIPs(binded InterfacesAddrFunc) AllowedIPAddresses {
	return func(ctx context.Context) ([]net.IP, error) {
		ips := getBindedIPAddress(ctx, binded)
		ips = filterPublicOnlyIPs(ips)
		ips = truncatedCopyIPs(ips)
		return ips, nil
	}
}

func getBindedIPAddress(ctx context.Context, interfacesAddr InterfacesAddrFunc) []net.IP {
	logger := zc.L(ctx)
	binded, err := interfacesAddr()
	log.DebugDPanic(logger, err, "Get system addresses", zap.Any("addresses", binded))

	var parsed = make([]net.IP, 0, len(binded))

	for _, addr := range binded {
		addrS := addr.String()
		if addrS == "<nil>" {
			continue
		}
		ip, _, err := net.ParseCIDR(addrS)
		log.DebugDPanic(logger, err, "Parse ip from interface", zap.String("addr", addrS), zap.Any("ip", ip),
			zap.Stringer("original_ip", addr))
		if ip == nil {
			continue
		}

		logger.Debug("Parse ip", zap.Stringer("ip", ip))
		parsed = append(parsed, ip)
	}
	return parsed
}

func filterPublicOnlyIPs(ips []net.IP) []net.IP {
	var public = make([]net.IP, 0, len(ips))
	for _, ip := range ips { // nolint:wsl
		if isPublicIP(ip) {
			public = append(public, ip)
		}
	}
	return public
}
