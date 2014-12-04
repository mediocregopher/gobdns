package main

import (
	_ "github.com/mediocregopher/gobdns/dns"
	"github.com/mediocregopher/gobdns/ips"
)

func init() {
	ips.Set("turtles", "127.0.0.1")
	ips.Set("dev.turtles", "127.0.0.2")
}

func main() {
	select {}
}
