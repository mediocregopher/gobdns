package persist

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/ips"
)

func init() {
	if config.BackupFile == "" {
		return
	}

	if initialRead() {
		go spin()
	}
}

func initialRead() bool {
	f, err := os.Open(config.BackupFile)
	if os.IsNotExist(err) {
		log.Printf("%s not found, not reading from it", config.BackupFile)
		return true
	} else if err != nil {
		log.Fatalf("%s could not be read from: %s", config.BackupFile, err)
	}

	log.Printf("Reading mappings from %s", config.BackupFile)
	buf := bufio.NewReader(f)
	m := map[string]string{}
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("%s could not be read from: %s", config.BackupFile, err)
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			log.Fatalf("%s is corrupted, found %q", config.BackupFile, line)
		}
		m[parts[0]] = strings.TrimSpace(parts[1])
	}

	for domain, ip := range m {
		ips.Set(domain, ip)
	}

	return true
}

func spin() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			tmpfn := config.BackupFile + ".tmp"
			tmpf, err := os.Create(tmpfn)
			if err != nil {
				log.Printf("Couldn't open %s for writing: %s", tmpfn, err)
				break
			}
			for domain, ip := range ips.GetAll() {
				fmt.Fprintf(tmpf, "%s %s\n", domain, ip)
			}
			tmpf.Close()

			if err := os.Rename(tmpfn, config.BackupFile); err != nil {
				log.Printf("Couldn't rename %s: %s", tmpfn, err)
				break
			}
		}
	}
}
