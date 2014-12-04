package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/ips"
)

var usage = `
	GET    /                     Gives you this page

	GET    /api/domains/all      Gives you a space separated mapping of domains
	                             to ips

	POST   /api/domains/<domain> Maps the given domain to the given ip, which
	                             will be the body data for the request

	PUT    /api/domains/<domain> Same as POST'ing

	DELETE /api/domains/<domain> Removes the domain->ip mapping for the given
	                             domain

`

func init() {
	if config.APIAddr == "" {
		return
	}

	go func() {
		log.Printf("API Listening on %s", config.APIAddr)
		http.HandleFunc("/api/domains/all", getAll)
		http.HandleFunc("/api/domains/", putDelete)
		http.HandleFunc("/", root)
		http.ListenAndServe(config.APIAddr, nil)
	}()
}

func getAll(w http.ResponseWriter, r *http.Request) {
	for domain, ip := range ips.GetAll() {
		fmt.Fprintf(w, "%s %s\n", domain, ip)
	}
}

func putDelete(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/domains/" {
		w.WriteHeader(404)
		return
	}

	domain := path.Base(r.URL.Path)

	switch r.Method {
	case "PUT", "POST":
		ipB, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		ip := string(ipB)
		if ip == "" {
			w.WriteHeader(400)
			return
		}
		ips.Set(domain, ip)

	case "DELETE":
		ips.Unset(domain)
	}

	r.Body.Close()
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, usage)
}
