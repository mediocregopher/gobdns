package ips

import (
	"sync"
	"strings"
)

var m = map[string]string{}
var mLock sync.RWMutex

// Retrieves the best matching ip for the given domain. Or not. Will return the
// most specific match.
func GetIP(domain string) (string, bool) {
	for {
		mLock.RLock()
		ip, ok := m[domain]
		mLock.RUnlock()
		if ok {
			return ip, true
		}

		i := strings.IndexAny(domain, ".-")
		if i == -1 {
			return "", false
		}
		domain = domain[i+1:]
	}
}

// Get's a copy of the map which maps all known domains to all known ips
func GetAllIPs() map[string]string {
	m2 := map[string]string{}
	mLock.RLock()
	defer mLock.RUnlock()
	for domain, ip := range m {
		m2[domain] = ip
	}
	return m2
}

// Sets the given domain to point to the given ip. If the given domain doesn't
// end in a period one will be appended to it
func SetIP(domain, ip string) {
	if domain == "" {
		return
	}
	if domain[len(domain)-1] != '.' {
		domain = domain + "."
	}

	mLock.Lock()
	m[domain] = ip
	mLock.Unlock()
}

// If the given domain is set to an ip, unsets it
func UnsetIP(domain string) {
	mLock.Lock()
	delete(m, domain)
	mLock.Unlock()
}
