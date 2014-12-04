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
	APIAddr     string
	BackupFile  string
	MasterAddr  string
)

const forwardSuffix = ":53"

func init() {
	fc := flagconfig.New("gobdns")

	fc.StrParam("tcp-addr", "TCP address to listen on. Set to empty to not listen on tcp", ":53")
	fc.StrParam("udp-addr", "UDP address to listen on. Set to empty to not listen on udp", ":53")
	fc.StrParam("forward-addr", "Address to forward requests to when no matches are found", "")
	fc.StrParam("api-addr", "Address for the REST API to listen on. Set to empty to disable it", ":8080")
	fc.StrParam("backup-file", "File to read data from during startup and to write data to during runtime. Leave blank to disable persistance", "./gobdns.db")
	fc.StrParam("master-addr", "ip:port of master instance to periodically pull snapshots from. Leave blank to disable", "")

	if err := fc.Parse(); err != nil {
		log.Fatal(err)
	}

	TCPAddr = fc.GetStr("tcp-addr")
	UDPAddr = fc.GetStr("udp-addr")
	ForwardAddr = fc.GetStr("forward-addr")
	APIAddr = fc.GetStr("api-addr")
	BackupFile = fc.GetStr("backup-file")
	MasterAddr = fc.GetStr("master-addr")

	if ForwardAddr != "" && !strings.HasSuffix(ForwardAddr, forwardSuffix) {
		ForwardAddr += forwardSuffix
	}
}
