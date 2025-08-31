package dnsresolver

import (
	"github.com/ginnysingh789/dns-tracer/pkg/util"
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type TraceStep struct {
	Server        string
	Query         string
	Response      *dns.Msg
	Duration      time.Duration
	Error         error
	IsTCPFallback bool
	ServerType    string
	Level         int
}

func Resolve(domainName string) []TraceStep {
	rootServers := []string{"198.41.0.4"} // Starting with a root server
	serversToQuery := rootServers
	fqdn := dns.Fqdn(domainName)
	currentServerType := "Root Server"
	currentLevel := 0

	trace := []TraceStep{}

resolveLoop:
	for {
		if len(serversToQuery) == 0 {
			break
		}

		server := serversToQuery[0]
		m := new(dns.Msg)
		m.SetQuestion(fqdn, dns.TypeA)

		queryString := fmt.Sprintf("%s A", strings.TrimSuffix(fqdn, "."))

		clientUDP := new(dns.Client)
		start := time.Now()
		in, _, err := clientUDP.Exchange(m, server+":53")
		duration := time.Since(start)

		if err != nil {
			trace = append(trace, TraceStep{
				Server:     server,
				Query:      queryString,
				Duration:   duration,
				Error:      err,
				ServerType: currentServerType,
				Level:      currentLevel,
			})
			if len(serversToQuery) > 1 {
				serversToQuery = serversToQuery[1:]
				continue
			} else {
				break
			}
		}

		isTCP := false
		if in.Truncated {
			isTCP = true
			clientTCP := &dns.Client{Net: "tcp"}
			tcpStart := time.Now()
			in, _, err = clientTCP.Exchange(m, server+":53")
			tcpDuration := time.Since(tcpStart)
			duration += tcpDuration

			if err != nil {
				trace = append(trace, TraceStep{
					Server:        server,
					Query:         queryString,
					Duration:      duration,
					Error:         err,
					IsTCPFallback: true,
					ServerType:    currentServerType,
					Level:         currentLevel,
				})
				if len(serversToQuery) > 1 {
					serversToQuery = serversToQuery[1:]
					continue
				} else {
					break
				}
			}
		}

		step := TraceStep{
			Server:        server,
			Query:         queryString,
			Response:      in,
			Duration:      duration,
			IsTCPFallback: isTCP,
			ServerType:    currentServerType,
			Level:         currentLevel,
		}
		trace = append(trace, step)

		if in.Rcode == dns.RcodeNameError {
			break
		}
		if in.Rcode == dns.RcodeServerFailure {
			if len(serversToQuery) > 1 {
				serversToQuery = serversToQuery[1:]
				continue
			} else {
				break
			}
		}

		// Process answers - CNAME takes absolute priority
		if len(in.Answer) > 0 {
			// First pass: look for CNAME only
			for _, record := range in.Answer {
				if cname, ok := record.(*dns.CNAME); ok {
					fqdn = cname.Target
					serversToQuery = rootServers
					currentServerType = "Root Server"
					currentLevel = 0
					continue resolveLoop // Immediately restart, ignore everything else
				}
			}

			// Only if no CNAME was found, look for final answers
			for _, record := range in.Answer {
				if _, ok := record.(*dns.A); ok {
					break resolveLoop // Found final answer
				}
				if _, ok := record.(*dns.AAAA); ok {
					break resolveLoop // Found final answer
				}
			}
		} else {
			// Follow referral
			nextServers := []string{}
			for _, rr := range append(in.Ns, in.Extra...) {
				if a, ok := rr.(*dns.A); ok {
					nextServers = append(nextServers, a.A.String())
				}
			}
			if len(nextServers) > 0 {
				serversToQuery = nextServers
				currentServerType, currentLevel = util.DetermineServerTypeAndLevel(in.Ns, currentLevel)
			} else {
				break
			}
		}
	}
	return trace
}
