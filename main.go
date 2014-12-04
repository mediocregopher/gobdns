package main

import (
	_ "github.com/mediocregopher/gobdns/dns"
	_ "github.com/mediocregopher/gobdns/http"
	_ "github.com/mediocregopher/gobdns/persist"
	_ "github.com/mediocregopher/gobdns/repl"
)

func main() {
	select {}
}
