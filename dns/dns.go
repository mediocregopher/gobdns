package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"log"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/ips"
)

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {

	domain := r.Question[0].Name

	if ip, ok := ips.Get(domain); ok {
		a, err := dns.NewRR(fmt.Sprintf("%s IN A %s", domain, ip))
		if err != nil {
			log.Println(err)
			dns.HandleFailed(w, r)
			return
		}

		m := new(dns.Msg)
		m.SetReply(r)
		m.Answer = []dns.RR{a}
		w.WriteMsg(m)
		return
	}

	if config.ForwardAddr == "" {
		dns.HandleFailed(w, r)
		return
	}

	proxiedR, err := dns.Exchange(r, config.ForwardAddr)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteMsg(proxiedR)
}

func init() {
	handler := dns.HandlerFunc(handleRequest)
	if config.UDPAddr != "" {
		go func() {
			log.Printf("Listening on UDP %s", config.UDPAddr)
			err := dns.ListenAndServe(config.UDPAddr, "udp", handler)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	if config.TCPAddr != "" {
		go func() {
			log.Printf("Listening on TCP %s", config.TCPAddr)
			err := dns.ListenAndServe(config.TCPAddr, "tcp", handler)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}
