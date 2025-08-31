package display

import (
	"github.com/ginnysingh789/dns-tracer/pkg/dnsresolver"
	"github.com/ginnysingh789/dns-tracer/pkg/util"
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func TraceTree(trace []dnsresolver.TraceStep) {
	fmt.Printf("\nDNS Resolution Path\n")
	fmt.Printf("==================\n\n")

	totalTime := time.Duration(0)

	for _, step := range trace {
		totalTime += step.Duration

		// Natural indentation
		indent := ""
		connector := ""

		switch step.Level {
		case 0:
			connector = "ROOT: "
		case 1:
			indent = "  "
			connector = "└─ TLD: "
		case 2:
			indent = "    "
			connector = "└─ AUTH: "
		default:
			indent = strings.Repeat("  ", step.Level)
			connector = "└─ "
		}

		serverName := util.GetServerName(step.Server)

		// Show query
		fmt.Printf("%s%sAsking %s about %s\n", indent, connector, step.Server, step.Query)
		if serverName != "Unknown Server" {
			fmt.Printf("%s       (%s)\n", indent, serverName)
		}

		if step.IsTCPFallback {
			fmt.Printf("%s       Used TCP due to large response\n", indent)
		}

		fmt.Printf("%s       Response time: %v\n", indent, step.Duration)

		if step.Error != nil {
			fmt.Printf("%s       ERROR: %v\n\n", indent, step.Error)
			continue
		}
		if step.Response == nil {
			fmt.Printf("%s       No response received\n\n", indent)
			continue
		}

		if step.Response.Rcode != dns.RcodeSuccess {
			rcodeStr := dns.RcodeToString[step.Response.Rcode]
			fmt.Printf("%s       Server says: %s", indent, rcodeStr)
			if step.Response.Rcode == dns.RcodeNameError {
				fmt.Print(" (domain doesn't exist)")
			}
			fmt.Printf("\n\n")
			continue
		}

		// Show what we got back
		if len(step.Response.Answer) > 0 {
			fmt.Printf("%s       Got answer:\n", indent)

			// Check if this response has a CNAME
			hasCname := false
			for _, record := range step.Response.Answer {
				if cname, ok := record.(*dns.CNAME); ok {
					fmt.Printf("%s         %s is actually %s\n", indent,
						strings.TrimSuffix(cname.Header().Name, "."),
						strings.TrimSuffix(cname.Target, "."))
					fmt.Printf("%s         Need to look up the real name now...\n", indent)
					hasCname = true
					break // Only show the CNAME, ignore A records in same response
				}
			}

			// Only show A records if no CNAME was found
			if !hasCname {
				for _, record := range step.Response.Answer {
					if a, ok := record.(*dns.A); ok {
						fmt.Printf("%s         Final IP address: %s\n", indent, a.A)
					}
				}
			}
		} else if len(step.Response.Ns) > 0 {
			fmt.Printf("%s       Server doesn't know, but suggests asking:\n", indent)

			// Show a few name servers
			nsCount := 0
			for _, record := range step.Response.Ns {
				if ns, ok := record.(*dns.NS); ok {
					if nsCount < 2 {
						fmt.Printf("%s         %s\n", indent, strings.TrimSuffix(ns.Ns, "."))
					}
					nsCount++
				}
			}
			if nsCount > 2 {
				fmt.Printf("%s         (and %d others)\n", indent, nsCount-2)
			}

			// Show IP addresses we can use
			glueCount := 0
			fmt.Printf("%s       Here are their IP addresses:\n", indent)
			for _, record := range step.Response.Extra {
				if a, ok := record.(*dns.A); ok {
					if glueCount < 2 {
						fmt.Printf("%s         %s at %s\n", indent,
							strings.TrimSuffix(a.Header().Name, "."), a.A)
					}
					glueCount++
				}
			}
			if glueCount > 2 {
				fmt.Printf("%s         (and %d more)\n", indent, glueCount-2)
			}
		}

		fmt.Println()
	}

	fmt.Printf("Summary\n")
	fmt.Printf("-------\n")
	fmt.Printf("Total queries: %d\n", len(trace))
	fmt.Printf("Total time: %v\n", totalTime)

	// Find final result
	finalIP := "Not resolved"
	for i := len(trace) - 1; i >= 0; i-- {
		step := trace[i]
		if step.Response != nil && len(step.Response.Answer) > 0 {
			for _, ans := range step.Response.Answer {
				if a, ok := ans.(*dns.A); ok {
					finalIP = a.A.String()
					break
				}
			}
		}
		if finalIP != "Not resolved" {
			break
		}
	}
	fmt.Printf("Final result: %s\n", finalIP)
}
