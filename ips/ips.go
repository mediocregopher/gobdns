package ips

import (
	"sync"
	"strings"
)

type Mapping map[string]string

var m = Mapping{}
var mLock sync.RWMutex

// Retrieves the best matching ip for the given domain. Or not. Will return the
// most specific match.
func Get(domain string) (string, bool) {
	for {
		ip, ok := GetExact(domain)
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

// Retrieves an ip exactly matching the given domain (although this will append
// the period to the end if necessary)
func GetExact(domain string) (string, bool) {
	mLock.RLock()
	defer mLock.RUnlock()
	ip, ok := m[appendPeriod(domain)]
	return ip, ok
}

// Get's a copy of the map which maps all known domains to all known ips
func GetAll() Mapping {
	m2 := Mapping{}
	mLock.RLock()
	defer mLock.RUnlock()
	for domain, ip := range m {
		m2[domain] = ip
	}
	return m2
}

// Sets the given domain to point to the given ip. If the given domain doesn't
// end in a period one will be appended to it
func Set(domain, ip string) {
	if domain == "" {
		return
	}
	domain = appendPeriod(domain)

	mLock.Lock()
	m[domain] = ip
	mLock.Unlock()
}

// Sets the current snapshot to a copy of the given one
func SetAll(m2 Mapping) {
	mLock.Lock()
	defer mLock.Unlock()
	m = Mapping{}
	for domain, ip := range m2 {
		m[domain] = ip
	}
}

// If the given domain is set to an ip, unsets it
func Unset(domain string) {
	mLock.Lock()
	delete(m, domain)
	mLock.Unlock()
}

func appendPeriod(domain string) string {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}
	return domain
}
