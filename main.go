package main

import (
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
)

type server struct {
	addr net.IP
}

func (s *server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}

	fmt.Printf("new DNS request: %q\n", r.Question[0].Name)

	// TODO: make sure this is for the subdomain we care about
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = append(m.Answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		A: s.addr,
	})

	if err := w.WriteMsg(m); err != nil {
		panic(err)
	}
}

func main() {
	s := &server{addr: net.ParseIP("5.78.103.252")}
	if val := os.Getenv("WILD_ADDR"); val != "" {
		s.addr = net.ParseIP(val)
		if s.addr == nil {
			panic(fmt.Sprintf("parse WILD_ADDR failed"))
		}
	}

	pc, err := net.ListenPacket("udp", "0.0.0.0:53")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	fmt.Printf("listening on 0.0.0.0:53...\n")

	if err := dns.ActivateAndServe(nil, pc, s); err != nil {
		panic(err)
	}
}
