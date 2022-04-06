package dns

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
	mdns "github.com/miekg/dns"
)

var (
	_ ResolverInterface = &Resolver{}
)

func TestNewResolver(t *testing.T) {
	td := testdeep.NewT(t)
	resolver := NewResolver("1.2.3.4:53")
	td.CmpDeeply(resolver.server, "1.2.3.4:53")
	td.CmpDeeply(resolver.maxDNSRecursionDeep, 10)
	td.NotNil(resolver.tcp)
	td.NotNil(resolver.udp)
	td.NotNil(resolver.lookupWithClient)
}

func TestLookupWithClient(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)

	client := NewMDNSClientMock(mc)
	client.ExchangeMock.Set(func(m *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error) {
		if m.Id == 0 {
			td.Error("Unset msg id")
		}
		switch m.Question[0].Qtype {
		case mdns.TypeA:
			switch m.Question[0].Name {
			case "alias.com":
				td.CmpDeeply(m.Question, []mdns.Question{{Name: "alias.com", Qtype: mdns.TypeA, Qclass: mdns.ClassINET}})
				return &mdns.Msg{
					MsgHdr: mdns.MsgHdr{Id: m.Id},
					Answer: []mdns.RR{
						&mdns.CNAME{
							Hdr: mdns.RR_Header{
								Rrtype: mdns.TypeCNAME,
							},
							Target: "alias-target.com",
						},
					},
				}, 0, nil
			case "alias-target.com":
				td.CmpDeeply(m.Question, []mdns.Question{{Name: "alias-target.com", Qtype: mdns.TypeA, Qclass: mdns.ClassINET}})
				return &mdns.Msg{
					MsgHdr: mdns.MsgHdr{Id: m.Id},
					Answer: []mdns.RR{
						&mdns.A{
							Hdr: mdns.RR_Header{Rrtype: mdns.TypeA},
							A:   net.IPv4(3, 4, 5, 6),
						},
					},
				}, 0, nil
			case "asd.com":
				td.CmpDeeply(m.Question, []mdns.Question{{Name: "asd.com", Qtype: mdns.TypeA, Qclass: mdns.ClassINET}})
				return &mdns.Msg{
					MsgHdr: mdns.MsgHdr{Id: m.Id},
					Answer: []mdns.RR{
						&mdns.A{
							Hdr: mdns.RR_Header{Rrtype: mdns.TypeA},
							A:   net.IPv4(1, 2, 3, 4),
						},
						&mdns.A{
							Hdr: mdns.RR_Header{Rrtype: mdns.TypeA},
							A:   net.IPv4(5, 6, 7, 8),
						},
					},
				}, 0, nil
			default:
				td.Error("Unexpected domain")
				return nil, 0, errors.New("unexpected domain")
			}
		case mdns.TypeAAAA:
			td.CmpDeeply(m.Question, []mdns.Question{{Name: "asd.com", Qtype: mdns.TypeAAAA, Qclass: mdns.ClassINET}})
			return &mdns.Msg{
				MsgHdr: mdns.MsgHdr{Id: m.Id},
				Answer: []mdns.RR{
					&mdns.AAAA{
						Hdr:  mdns.RR_Header{Rrtype: mdns.TypeAAAA},
						AAAA: net.ParseIP("::1a"),
					},
					&mdns.AAAA{
						Hdr:  mdns.RR_Header{Rrtype: mdns.TypeAAAA},
						AAAA: net.ParseIP("::1f"),
					},
				},
			}, 0, nil
		case mdns.TypeAVC:
			return &mdns.Msg{
				MsgHdr: mdns.MsgHdr{Id: m.Id},
				Answer: []mdns.RR{
					&mdns.AVC{
						Hdr: mdns.RR_Header{Rrtype: mdns.TypeAVC},
						Txt: []string{"aaa"},
					},
				},
			}, 0, nil

		default:
			td.Error("Unexpected record type", m.Question)
			return nil, 0, errors.New("unexpected")
		}
	})

	ips, err := lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeA, 1, client)
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.IPv4(1, 2, 3, 4)}, {IP: net.IPv4(5, 6, 7, 8)}})

	ips, err = lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeAAAA, 1, client)
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.ParseIP("::1a")}, {IP: net.ParseIP("::1f")}})

	ips, err = lookupWithClient(ctx, "alias.com", "1.2.3.4:53", mdns.TypeA, 2, client)
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.IPv4(3, 4, 5, 6)}})

	ips, err = lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeAVC, 1, client)
	td.CmpNoError(err)
	td.Nil(ips)

	ips, err = lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeAVC, 0, client)
	td.CmpError(err)
	td.Nil(ips)

	client.ExchangeMock.Set(func(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error) {
		return &mdns.Msg{MsgHdr: mdns.MsgHdr{Id: msg.Id, Truncated: true}}, 0, nil
	})
	ips, err = lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeA, 1, client)
	td.CmpDeeply(err, errTruncatedResponse)
	td.Nil(ips)

	client.ExchangeMock.Set(func(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error) {
		return &mdns.Msg{MsgHdr: mdns.MsgHdr{Id: msg.Id, Truncated: true}}, 0, errors.New("asd")
	})
	ips, err = lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeA, 1, client)
	td.CmpDeeply(err, errTruncatedResponse)
	td.Nil(ips)

	client.ExchangeMock.Set(func(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error) {
		return &mdns.Msg{MsgHdr: mdns.MsgHdr{Id: msg.Id}}, 0, errors.New("asd")
	})
	ips, err = lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeA, 1, client)
	td.CmpDeeply(err, errors.New("asd"))
	td.Nil(ips)

	client.ExchangeMock.Set(func(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error) {
		return &mdns.Msg{
			MsgHdr: mdns.MsgHdr{Id: msg.Id},
			Answer: []mdns.RR{
				&mdns.AAAA{
					Hdr:  mdns.RR_Header{Rrtype: mdns.TypeAAAA},
					AAAA: net.ParseIP("::1a"),
				},
				&mdns.AAAA{
					Hdr:  mdns.RR_Header{Rrtype: mdns.TypeAAAA},
					AAAA: net.ParseIP("::1f"),
				},
				&mdns.A{
					Hdr: mdns.RR_Header{Rrtype: mdns.TypeA},
					A:   net.IPv4(1, 2, 3, 4),
				},
				&mdns.A{
					Hdr: mdns.RR_Header{Rrtype: mdns.TypeA},
					A:   net.IPv4(5, 6, 7, 8),
				},
			},
		}, 0, nil
	})
	ips, err = lookupWithClient(ctx, "asd.com", "1.2.3.4:53", mdns.TypeA, 1, client)
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.IPv4(1, 2, 3, 4)}, {IP: net.IPv4(5, 6, 7, 8)}})

	client.ExchangeMock.Set(func(msg *mdns.Msg, address string) (r *mdns.Msg, rtt time.Duration, err error) {
		time.Sleep(time.Second)
		return &mdns.Msg{
			MsgHdr: mdns.MsgHdr{Id: msg.Id},
			Answer: []mdns.RR{
				&mdns.A{
					Hdr: mdns.RR_Header{Rrtype: mdns.TypeA},
					A:   net.IPv4(1, 2, 3, 4),
				},
				&mdns.A{
					Hdr: mdns.RR_Header{Rrtype: mdns.TypeA},
					A:   net.IPv4(5, 6, 7, 8),
				},
			},
		}, 0, nil
	})

	timeoutCtx, timeoutCancelCtx := context.WithTimeout(ctx, time.Millisecond*10)
	defer timeoutCancelCtx()

	ips, err = lookupWithClient(timeoutCtx, "asd.com", "1.2.3.4:53", mdns.TypeA, 1, client)
	td.CmpError(err)
	td.Nil(ips)

	timeoutCancelled, timeoutCancelledCancelCtx := context.WithCancel(ctx)
	timeoutCancelledCancelCtx()
	ips, err = lookupWithClient(timeoutCancelled, "asd.com", "1.2.3.4:53", mdns.TypeA, 1, client)
	td.CmpError(err)
	td.Nil(ips)
}

func TestResolver_Lookup(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)

	clientUDP := NewMDNSClientMock(mc)
	clientTCP := NewMDNSClientMock(mc)

	resolver := &Resolver{
		udp:                 clientUDP,
		tcp:                 clientTCP,
		server:              "dns",
		maxDNSRecursionDeep: 13,
	}

	resolver.lookupWithClient = func(ctx context.Context, host string, server string, recordType uint16, recursion int, client mDNSClient) (addrs []net.IPAddr, e error) {
		td.CmpDeeply(recursion, 13)
		answer, _, err := client.Exchange(&mdns.Msg{Question: []mdns.Question{
			{Name: host, Qtype: recordType},
		}}, server)
		if err != nil {
			return nil, err
		}
		return []net.IPAddr{{IP: answer.Answer[0].(*mdns.A).A}}, nil
	}

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "1", Qtype: mdns.TypeA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.A{A: net.IPv4(1, 2, 3, 4)}}}, 0, nil)
	ips, err := resolver.lookup(ctx, "1", mdns.TypeA)
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.IPv4(1, 2, 3, 4)}})

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "2", Qtype: mdns.TypeA},
	}}, "dns").
		Then(nil, 0, errTruncatedResponse)
	clientTCP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "2", Qtype: mdns.TypeA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.A{A: net.IPv4(1, 2, 3, 4)}}}, 0, nil)
	ips, err = resolver.lookup(ctx, "2", mdns.TypeA)
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.IPv4(1, 2, 3, 4)}})

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "3", Qtype: mdns.TypeA},
	}}, "dns").
		Then(nil, 0, errors.New("test3"))
	ips, err = resolver.lookup(ctx, "3", mdns.TypeA)
	td.CmpDeeply(err, errors.New("test3"))
	td.Nil(ips)

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "4", Qtype: mdns.TypeA},
	}}, "dns").
		Then(nil, 0, errTruncatedResponse)
	clientTCP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "4", Qtype: mdns.TypeA},
	}}, "dns").
		Then(nil, 0, errors.New("test4"))
	ips, err = resolver.lookup(ctx, "4", mdns.TypeA)
	td.CmpDeeply(err, errors.New("test4"))
	td.Nil(ips)
}

func TestResolver_LookupIPAddr(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	mc := minimock.NewController(td)

	clientUDP := NewMDNSClientMock(mc)
	clientTCP := NewMDNSClientMock(mc)

	resolver := &Resolver{
		udp:                 clientUDP,
		tcp:                 clientTCP,
		server:              "dns",
		maxDNSRecursionDeep: 13,
	}

	resolver.lookupWithClient = func(ctx context.Context, host string, server string, recordType uint16, recursion int, client mDNSClient) (addrs []net.IPAddr, e error) {
		td.CmpDeeply(recursion, 13)
		dnsAnswer, _, err := client.Exchange(&mdns.Msg{Question: []mdns.Question{
			{Name: host, Qtype: recordType},
		}}, server)
		if err != nil {
			return nil, err
		}
		var resIPs []net.IPAddr
		for _, r := range dnsAnswer.Answer {
			switch recordType {
			case mdns.TypeA:
				resIPs = append(resIPs, net.IPAddr{IP: r.(*mdns.A).A})
			case mdns.TypeAAAA:
				resIPs = append(resIPs, net.IPAddr{IP: r.(*mdns.AAAA).AAAA})
			default:
				// pass
			}
		}
		return resIPs, nil
	}

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "1.", Qtype: mdns.TypeA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.A{A: net.IPv4(1, 2, 3, 4)}}}, 0, nil)
	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "1.", Qtype: mdns.TypeAAAA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.AAAA{AAAA: net.ParseIP("::bb")}}}, 0, nil)
	ips, err := resolver.LookupIPAddr(ctx, "1")
	td.CmpNoError(err)
	td.CmpDeeply(ips, []net.IPAddr{{IP: net.IPv4(1, 2, 3, 4)}, {IP: net.ParseIP("::bb")}})

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "2.", Qtype: mdns.TypeA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.A{A: net.IPv4(1, 2, 3, 4)}}}, 0, errors.New("err2-1"))
	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "2.", Qtype: mdns.TypeAAAA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.AAAA{AAAA: net.ParseIP("::bb")}}}, 0, nil)
	ips, err = resolver.LookupIPAddr(ctx, "2")
	td.CmpDeeply(err, errors.New("err2-1"))
	td.Nil(ips)

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "3.", Qtype: mdns.TypeA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.A{A: net.IPv4(1, 2, 3, 4)}}}, 0, errors.New("err3-1"))
	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "3.", Qtype: mdns.TypeAAAA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.AAAA{AAAA: net.ParseIP("::bb")}}}, 0, errors.New("err3-2"))
	ips, err = resolver.LookupIPAddr(ctx, "3")
	td.CmpDeeply(err, errors.New("err3-1"))
	td.Nil(ips)

	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "4.", Qtype: mdns.TypeA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.A{A: net.IPv4(1, 2, 3, 4)}}}, 0, nil)
	clientUDP.ExchangeMock.When(&mdns.Msg{Question: []mdns.Question{
		{Name: "4.", Qtype: mdns.TypeAAAA},
	}}, "dns").
		Then(&mdns.Msg{Answer: []mdns.RR{&mdns.AAAA{AAAA: net.ParseIP("::bb")}}}, 0, errors.New("err4-2"))
	ips, err = resolver.LookupIPAddr(ctx, "4")
	td.CmpDeeply(err, errors.New("err4-2"))
	td.Nil(ips)
}

func TestResolverReal(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)

	var ips []net.IPAddr
	var err error

	r := NewResolver("8.8.8.8:53")
	ips, err = r.LookupIPAddr(ctx, "one.one.one.one")
	td.CmpNoError(err)
	td.Contains(ips,
		testdeep.Any(
			net.IPAddr{IP: net.IPv4(1, 1, 1, 1).To4()},
			net.IPAddr{IP: net.IPv4(1, 1, 1, 1).To16()},
		),
	)

	r = NewResolver("1.1.1.1:53")
	ips, err = r.LookupIPAddr(ctx, "test.l.rekby.ru")
	td.CmpNoError(err)
	td.Contains(ips,
		testdeep.Any(
			net.IPAddr{IP: net.ParseIP("127.0.0.1").To4()},
			net.IPAddr{IP: net.ParseIP("127.0.0.1").To16()},
		),
	)
}
