package config

import (
	"log"
	"strings"

	"github.com/mediocregopher/lever"
	"github.com/miekg/dns"
)

// ForwardSuffix describes a forward-suffix-addr parameter specified by the user
type ForwardSuffix struct {
	ForwardAddr string
	Suffix      string
}

var (
	TCPAddr         string
	UDPAddr         string
	ForwardAddr     string
	ForwardSuffixes []ForwardSuffix
	APIAddr         string
	BackupFile      string
	MasterAddr      string
)

const forwardPort = ":53"

func init() {
	l := lever.New("gobdns", nil)

	l.Add(lever.Param{
		Name:        "--tcp-addr",
		Description: "TCP address to listen on. Set to empty to not listen on tcp",
		Default:     ":53",
	})

	l.Add(lever.Param{
		Name:        "--udp-addr",
		Description: "UDP address to listen on. Set to empty to not listen on udp",
		Default:     ":53",
	})

	l.Add(lever.Param{
		Name:        "--forward-addr",
		Description: "Address to forward requests to when no matches are found",
	})

	l.Add(lever.Param{
		Name:        "--forward-suffix-addr",
		Description: "ip[:port]/suffix, all requests not matched internally but which have the given suffix will be forwarded to the given ip:port. Can be used more than once",
	})

	l.Add(lever.Param{
		Name:        "--api-addr",
		Description: "Address for the REST API to listen on. Set to empty to disable it",
		Default:     ":8080",
	})

	l.Add(lever.Param{
		Name:        "--backup-file",
		Description: "File to read data from during startup and to write data to during runtime. Leave blank to disable persistence",
		Default:     "./gobdns.db",
	})

	l.Add(lever.Param{
		Name:        "--master-addr",
		Description: "ip:port of master instance to periodically pull snapshots from. Leave blank to disable",
	})

	l.Parse()

	TCPAddr, _ = l.ParamStr("--tcp-addr")
	UDPAddr, _ = l.ParamStr("--udp-addr")
	ForwardAddr, _ = l.ParamStr("--forward-addr")
	forwardSuffixesStrs, _ := l.ParamStrs("--forward-suffix-addr")
	APIAddr, _ = l.ParamStr("--api-addr")
	BackupFile, _ = l.ParamStr("--backup-file")
	MasterAddr, _ = l.ParamStr("--master-addr")

	if ForwardAddr != "" && !strings.HasSuffix(ForwardAddr, forwardPort) {
		ForwardAddr += forwardPort
	}

	for _, f := range forwardSuffixesStrs {
		parts := strings.SplitN(f, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			log.Fatalf("Malformed forward-suffix-addr: %s", f)
		}
		ForwardSuffixes = append(ForwardSuffixes, ForwardSuffix{
			ForwardAddr: addForwardSuffix(parts[0]),
			Suffix:      dns.Fqdn(parts[1]),
		})
	}
}

func addForwardSuffix(f string) string {
	if f != "" && !strings.HasSuffix(f, forwardPort) {
		f += forwardPort
	}
	return f
}
