package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"fortio.org/fortio/log"
	"github.com/miekg/dns"
)

func usage() {
	fmt.Println("Usage: dnsping query server\neg:\tdnsping www.google.com. 127.0.0.1")
	os.Exit(1)
}

func main() {
	portFlag := flag.Int("p", 53, "`Port` to connect to")
	intervalFlag := flag.Duration("i", 1*time.Second, "How long to `wait` between requests")
	countFlag := flag.Int("c", 10, "How many `requests` to make")
	queryTypeFlag := flag.String("t", "A", "Query `type` to use (A, SOA, CNAME...)")
	flag.Parse()
	qt, exists := dns.StringToType[*queryTypeFlag]
	if !exists {
		log.Errf("Invalid type name %q", *queryTypeFlag)
		os.Exit(1)
	}
	args := flag.Args()
	nArgs := len(args)
	log.LogVf("got %d arguments: %v", nArgs, args)
	if nArgs != 2 {
		usage()
	}
	addrStr := fmt.Sprintf("%s:%d", args[1], *portFlag)
	m := new(dns.Msg)
	m.SetQuestion(args[0], qt)
	log.Infof("Will query server: %s for %s (%d) record for %s", addrStr, *queryTypeFlag, qt, args[0])
	log.LogVf("Query is: %v", m)
	for i := 1; i <= *countFlag; i++ {
		r, err := dns.Exchange(m, addrStr)
		if err != nil {
			log.Fatalf("failed to exchange: %v", err)
		}
		if r == nil {
			log.Fatalf("response is nil")
		}
		log.LogVf("response is %v", r)
		if r.Rcode != dns.RcodeSuccess {
			log.Errf("failed to get an valid answer: %v", r)
		}
		log.Infof("%d: %v", i, r.Answer)
		time.Sleep(*intervalFlag)
	}
}
