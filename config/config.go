package config

import (
	"github.com/mediocregopher/flagconfig"
	"log"
	"strings"
)

var (
	TCPAddr     string
	UDPAddr     string
	ForwardAddr string
)

const forwardSuffix = ":53"

func init() {
	fc := flagconfig.New("gobdns")

	fc.StrParam("tcp-addr", "TCP address to listen on. Set to empty to not listen on tcp", ":53")
	fc.StrParam("udp-addr", "UDP address to listen on. Set to empty to not listen on udp", ":53")
	fc.StrParam("forward-addr", "Address to forward requests to when no matches are found", "")

	if err := fc.Parse(); err != nil {
		log.Fatal(err)
	}

	TCPAddr = fc.GetStr("tcp-addr")
	UDPAddr = fc.GetStr("udp-addr")
	ForwardAddr = fc.GetStr("forward-addr")

	if ForwardAddr != "" && !strings.HasSuffix(ForwardAddr, forwardSuffix) {
		ForwardAddr += forwardSuffix
	}
}
