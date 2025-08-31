package main

import (
	"flag"
	"fmt"

	"github.com/miekg/dns"
)

func main() {
	// Define command-line flags
	//Define the flag

	//Iterative Process
	//Firs show any server,iterative over the server if the server become there is not next serve

	serverQuery := []string{"198.41.0.4"}

	var domainName string

	flag.StringVar(&domainName, "domain", "google.com", "Domain to Resolve (Default-> google.com)")
	flag.Parse()
	if domainName == "" {
		fmt.Println("Error domain flag is required")
		flag.Usage()
		return
	}

	//DNS
	for {
		fmt.Println(len(serverQuery))
		if len(serverQuery) == 0 {
			fmt.Println("There are No Server ")
			break
		}

		currentServer := serverQuery[0]
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(domainName), dns.TypeA)
		ans, err := dns.Exchange(m, currentServer+":53")

		if err != nil {
			fmt.Println("Error in the query", err)
			return
		}
		if len(ans.Answer) > 0 {
			fmt.Println("Answer is found")
			for _, respone := range ans.Answer {
				fmt.Println(respone)
			}
			break
		} else {
			nextServer := []string{}
			//The IP of next server always present in the Extra section of the answer
			for _, rr := range ans.Extra {
				if a, ok := rr.(*dns.A); ok {
					fmt.Println("Found refereal IP (glue record)", a.A.String())
					nextServer = append(nextServer, a.A.String())
				}
			}
			serverQuery = nextServer
		}

		// if len(ans.Answer) == 0 {
		// 	fmt.Println("No record found ")
		// }
	}

}
