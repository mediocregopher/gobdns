package snapshot

import (
	"encoding/json"

	"github.com/mediocregopher/gobdns/ips"
)

type Snapshot struct {
	ips.Mapping
}

func Create() *Snapshot {
	return &Snapshot{
		Mapping: ips.GetAll(),
	}
}

func CreateEncoded() ([]byte, error) {
	s := Create()
	return json.Marshal(s)
}

func Load(s *Snapshot) error {
	ips.SetAll(s.Mapping)
	return nil
}

func LoadEncoded(b []byte) error {
	s := &Snapshot{}
	if err := json.Unmarshal(b, s); err != nil {
		return err
	}
	return Load(s)
}
