package persist

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mediocregopher/gobdns/config"
	"github.com/mediocregopher/gobdns/snapshot"
)

func init() {
	if config.BackupFile == "" {
		return
	}

	initialRead()
	go spin()
}

func initialRead() {
	log.Printf("Reading snapshot from %s", config.BackupFile)
	b, err := ioutil.ReadFile(config.BackupFile)
	if os.IsNotExist(err) {
		log.Printf("%s not found, not reading from it", config.BackupFile)
		return
	} else if err != nil {
		log.Fatalf("%s could not be read from: %s", config.BackupFile, err)
	}

	if err := snapshot.LoadEncoded(b); err != nil {
		log.Fatalf("Could not load snapshot from %s: %s", config.BackupFile, err)
	}
}

func writeSnapshot() {
	tmpfn := config.BackupFile + ".tmp"
	tmpf, err := os.Create(tmpfn)
	if err != nil {
		log.Printf("Couldn't open %s for writing: %s", tmpfn, err)
		return
	}
	defer tmpf.Close()

	b, err := snapshot.CreateEncoded()
	if err != nil {
		log.Printf("Couldn't create snapshot: %s", err)
		return
	}

	if _, err := tmpf.Write(b); err != nil {
		log.Printf("Couldn't write snapshot to %s: %s", tmpfn, err)
		return
	}

	if err := os.Rename(tmpfn, config.BackupFile); err != nil {
		log.Printf("Couldn't rename %s: %s", tmpfn, err)
		return
	}
}

func spin() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			writeSnapshot()
		}
	}
}
