package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/ips"
	"github.com/mediocregopher/gobdns/snapshot"
)

var usage = `
	GET    /                     Gives you this page

	GET    /api/snapshot         Returns an encoded snapshot of this instance's
	                             data

	GET    /api/domains/all      Gives you a space separated mapping of domains
	                             to ips

	POST   /api/domains/<domain> Maps the given domain to the given ip, which
	                             can be the body data for the request, otherwise
	                             the ip the request is coming from will be used.

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
		http.HandleFunc("/api/snapshot", getSnapshot)
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
		ip := strings.TrimSpace(string(ipB))
		if ip == "" {
			portIdent := strings.LastIndex(r.RemoteAddr, ":")
			ip = r.RemoteAddr[:portIdent]
		} else if net.ParseIP(ip) == nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "invalid ip given in request body")
			return
		}
		ips.Set(domain, ip)

	case "DELETE":
		ips.Unset(domain)
	}
}

func getSnapshot(w http.ResponseWriter, r *http.Request) {
	b, err := snapshot.CreateEncoded()
	if err != nil {
		log.Printf("Creating snapshot: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, usage)
}
