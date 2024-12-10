package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/miekg/dns"
)

type server struct {
}

func (s *server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}

	fmt.Printf("dns: new type %04x request: %q\n", r.Question[0].Qtype, r.Question[0].Name)

	m := new(dns.Msg)
	m.SetReply(r)

	switch r.Question[0].Qtype {
	case dns.TypeTXT:
		m.Answer = append(m.Answer, &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   r.Question[0].Name,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			Txt: []string{os.Getenv("WILD_TXT")},
		})

	case dns.TypeA:
		// TODO: make sure this is for the subdomain we care about
		m.Authoritative = true
		m.Answer = append(m.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name:   r.Question[0].Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			A: net.ParseIP(os.Getenv("WILD_A")),
		})
	}

	if err := w.WriteMsg(m); err != nil {
		panic(err)
	}
}

func main() {
	pc, err := net.ListenPacket("udp", "0.0.0.0:53")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	fmt.Printf("listening on 0.0.0.0:53...\n")

	go func() {
		if err := dns.ActivateAndServe(nil, pc, new(server)); err != nil {
			panic(err)
		}
	}()

	cfg := make(map[string]string)
	b, err := os.ReadFile("/etc/wild.json")
	if err != nil {
		panic(fmt.Errorf("read /etc/wild.json: %w", err))
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		panic(fmt.Errorf("parse /etc/wild.json: %w", err))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uuid, ok := cfg[r.Host]
		if !ok {
			w.WriteHeader(404)
			return
		}

		url := fmt.Sprintf("https://api.trieve.ai/public_page/%s", uuid)
		if r.URL.Path != "/" {
			url = fmt.Sprintf("https://api.trieve.ai%s", r.URL.Path)
		}

		fmt.Printf("http: remapping %q -> %q\n", r.URL.String(), url)
		w2, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		for key, vals := range w2.Header {
			for _, val := range vals {
				w.Header().Set(key, val)
			}
		}
		w.WriteHeader(w2.StatusCode)
		io.Copy(w, w2.Body)
	})

	for {
		if err := http.ListenAndServeTLS("0.0.0.0:443", "/etc/letsencrypt/live/demo.trytrieve.com/fullchain.pem", "/etc/letsencrypt/live/demo.trytrieve.com/privkey.pem", nil); err != nil {
			fmt.Printf("listen http failed: %v\n", err)
			time.Sleep(5 * time.Second)
		}
	}
}
