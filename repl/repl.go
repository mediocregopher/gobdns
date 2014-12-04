package repl

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/ips"
)

func init() {
	if config.MasterAddr == "" {
		return
	}

	go spin()
}

func spin() {
	tick := time.Tick(5 * time.Second)
	snapshot()
	for {
		select {
		case <-tick:
			snapshot()
		}
	}
}

func snapshot() {
	url := "http://" + config.MasterAddr + "/api/domains/all"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error connecting to master: %s", err)
		return
	}
	defer resp.Body.Close()

	// TODO probably make a snapshot package to do all this so we're not
	// duplicating with persist
	m := map[string]string{}
	buf := bufio.NewReader(resp.Body)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("error reading from master: %s", err)
			return
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			log.Printf("corrupted snapshot at %q", line)
			return
		}

		m[parts[0]] = strings.TrimSpace(parts[1])
	}

	for domain, ip := range m {
		ips.Set(domain, ip)
	}
	return
}
