package repl

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/snapshot"
)

func init() {
	if config.MasterAddr == "" {
		return
	}

	go spin()
}

func spin() {
	tick := time.Tick(5 * time.Second)
	snapshotRequest()
	for {
		select {
		case <-tick:
			snapshotRequest()
		}
	}
}

func snapshotRequest() {
	url := "http://" + config.MasterAddr + "/api/snapshot"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error connecting to master: %s", err)
		return
	}
	defer resp.Body.Close()

	b := make([]byte, resp.ContentLength)
	if _, err := io.ReadFull(resp.Body, b); err != nil {
		log.Printf("error reading from master: %s", err)
		return
	}

	if err := snapshot.LoadEncoded(b); err != nil {
		log.Printf("error loading snapshot: %s", err)
		return
	}
}
