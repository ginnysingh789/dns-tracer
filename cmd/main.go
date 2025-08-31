package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ginnysingh789/dns-tracer/pkg/dnsresolver"
	"github.com/ginnysingh789/dns-tracer/pkg/display"
)

func main() {
	var domainName string
	flag.StringVar(&domainName, "domain", "", "The domain name to resolve iteratively")
	flag.Parse()

	if domainName == "" {
		fmt.Println("Error: -domain flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	trace := dnsresolver.Resolve(domainName)
	display.TraceTree(trace)
}
