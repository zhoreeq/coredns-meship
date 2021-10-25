// Package meship implements a plugin that returns details about the resolving
// querying it.
package meship

import (
	"context"

	_meshname "github.com/zhoreeq/meshname/pkg/meshname"

	"github.com/miekg/dns"
)

const name = "meship"

// Meship is a plugin that resolves .meship domains 
type Meship struct{}

// ServeDNS implements the plugin.Handler interface.
func (mn Meship) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	for _, q := range r.Question {
		labels := dns.SplitDomainName(q.Name)
		if len(labels) != 2 || q.Qtype != dns.TypeAAAA || q.Qclass != dns.ClassINET {
			//s.log.Debugln("Error: invalid resource requested")
			continue
		}

		if resolvedAddr, err := _meshname.IPFromDomain(&labels[0]); err == nil {
			answer := new(dns.AAAA)
			answer.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600}
			answer.AAAA = resolvedAddr

			a.Answer = append(a.Answer, answer)
		}
	}

	w.WriteMsg(a)
	return 0, nil
}

// Name implements the Handler interface.
func (mn Meship) Name() string { return name }
