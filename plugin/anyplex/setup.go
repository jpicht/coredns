package any

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"

	"github.com/caddyserver/caddy"
)

func init() { plugin.Register("anyplex", setup) }

func setup(c *caddy.Controller) error {
	a := &AnyPlex{
		Types: make([]dns.Type, 0, 1),
	}

	c.Next() // skip anyplex
	for c.Next() {
		if rt, ok := dns.StringToType[c.Val()]; ok {
			a.Types = append(a.Types, dns.Type(rt))
		} else if c.Val() == "extra" {
			a.Extra = true
		} else if c.Val() == "authority" {
			a.Authority = true
		} else {
			return c.Errf("unknown property '%s'", c.Val())
		}
	}

	if len(a.Types) == 0 {
		return c.Errf("ANYplex needs at least one type.")
	}

	log.Infof("ANYplex with %d types: %v", len(a.Types), a.Types)

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		log.Infof("ANYplex registered next: %#v", next)
		a.Next = next
		return a
	})

	return nil
}
