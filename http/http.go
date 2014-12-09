package http

import (
	"fmt"
	"io/ioutil"
	"github.com/elazarl/go-bindata-assetfs"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/ips"
	"github.com/mediocregopher/gobdns/snapshot"
)

func init() {
	if config.APIAddr == "" {
		return
	}

	go func() {
		log.Printf("API Listening on %s", config.APIAddr)
		http.HandleFunc("/api/domains/all", getAll)
		http.HandleFunc("/api/domains/", putDelete)
		http.HandleFunc("/api/snapshot", getSnapshot)

		assetFS := assetfs.AssetFS{
			Asset: Asset,
			AssetDir: AssetDir,
		}

		http.Handle("/", http.FileServer(&assetFS))
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
			ip = getIP(r)
		}
		if net.ParseIP(ip) == nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "invalid ip: %s", ip)
			return
		}
		ips.Set(domain, ip)

	case "DELETE":
		ips.Unset(domain)
	}
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	portIdent := strings.LastIndex(r.RemoteAddr, ":")
	return r.RemoteAddr[:portIdent]
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
