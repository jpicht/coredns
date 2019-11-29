package any

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

// AnyPlex is a plugin that returns a configurable subset of entries for ANY requests
type AnyPlex struct {
	Next plugin.Handler

	Types     []dns.Type
	Extra     bool
	Authority bool
}

func addToSet(dest []dns.RR, item dns.RR) []dns.RR {
	for _, e := range dest {
		if dns.IsDuplicate(e, item) {
			return dest
		}
	}
	return append(dest, item)
}

func addToSetMany(dest []dns.RR, src []dns.RR) []dns.RR {
	for _, o := range src {
		dest = addToSet(dest, o)
	}
	return dest
}

// ServeDNS implements the plugin.Handler interface.
func (a *AnyPlex) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// var s = ctx.Value(dnsserver.Key{}).(*dnsserver.Server)

	if r.Question[0].Qtype != dns.TypeANY {
		return plugin.NextOrFailure(a.Name(), a.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	hdr := dns.RR_Header{Name: r.Question[0].Name, Ttl: 8482, Class: dns.ClassINET, Rrtype: dns.TypeHINFO}
	m.Answer = []dns.RR{&dns.HINFO{Hdr: hdr, Cpu: "ANY obsoleted", Os: "See RFC 8482"}}

	if a.Next == nil {
		w.WriteMsg(m)
		return 0, nil
	}

	var (
		fakeRequest dns.Msg
	)

	for _, t := range a.Types {
		r.CopyTo(&fakeRequest)
		fakeRequest.Question[0].Qtype = uint16(t)
		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		n, err := a.Next.ServeDNS(ctx, rec, &fakeRequest)
		if err != nil {
			return n, err
		}
		// s.ServeDNS(ctx, rec, &fakeRequest)
		m.Answer = addToSetMany(m.Answer, rec.Msg.Answer)
		if a.Extra {
			m.Extra = addToSetMany(m.Extra, rec.Msg.Extra)
		}
		if a.Authority {
			m.Ns = addToSetMany(m.Ns, rec.Msg.Ns)
		}
	}

	w.WriteMsg(m)
	return 0, nil
}

// Name implements the Handler interface.
func (a *AnyPlex) Name() string { return "anyplex" }
