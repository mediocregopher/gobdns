package main

import (
	_ "github.com/mediocregopher/gobdns/dns"
	"github.com/mediocregopher/gobdns/ips"
)

func init() {
	ips.SetIP("turtles", "127.0.0.1")
	ips.SetIP("dev.turtles", "127.0.0.2")
}

func main() {
	select {}
}
