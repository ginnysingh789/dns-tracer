package util

import (
	"strings"

	"github.com/miekg/dns"
)

func GetServerName(ip string) string {
	knownServers := map[string]string{
		"198.41.0.4":     "a.root-servers.net",
		"199.9.14.201":   "b.root-servers.net",
		"192.33.4.12":    "c.root-servers.net",
		"199.7.91.13":    "d.root-servers.net",
		"192.203.230.10": "e.root-servers.net",
		"192.5.5.241":    "f.root-servers.net",
		"192.112.36.4":   "g.root-servers.net",
		"198.97.190.53":  "h.root-servers.net",
		"192.36.148.17":  "i.root-servers.net",
		"192.58.128.30":  "j.root-servers.net",
		"193.0.14.129":   "k.root-servers.net",
		"199.7.83.42":    "l.root-servers.net",
		"202.12.27.33":   "m.root-servers.net",
	}
	if name, ok := knownServers[ip]; ok {
		return name
	}
	return "Unknown Server"
}

func DetermineServerTypeAndLevel(ns []dns.RR, currentLevel int) (string, int) {
	if len(ns) == 0 {
		return "Unknown Server", currentLevel + 1
	}

	domain := ns[0].Header().Name
	domain = strings.TrimSuffix(domain, ".")
	parts := strings.Split(domain, ".")

	if len(parts) == 1 {
		return "TLD Server", 1
	} else if len(parts) >= 2 {
		return "Authoritative Server", 2
	}

	return "Unknown Server", currentLevel + 1
}
