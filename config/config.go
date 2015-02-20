package config

import (
	"strings"

	"github.com/mediocregopher/lever"
)

var (
	TCPAddr     string
	UDPAddr     string
	ForwardAddr string
	APIAddr     string
	BackupFile  string
	MasterAddr  string
)

const forwardSuffix = ":53"

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
	APIAddr, _ = l.ParamStr("--api-addr")
	BackupFile, _ = l.ParamStr("--backup-file")
	MasterAddr, _ = l.ParamStr("--master-addr")

	if ForwardAddr != "" && !strings.HasSuffix(ForwardAddr, forwardSuffix) {
		ForwardAddr += forwardSuffix
	}
}
