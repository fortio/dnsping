package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"fortio.org/dnsping/version"
	"fortio.org/fortio/log"
	"fortio.org/fortio/stats"
	"github.com/miekg/dns"
)

func usage() {
	fmt.Fprintln(flag.CommandLine.Output(),
		"dnsping "+version.Version+" usage:\n\tdnsping [flags] query server\neg:\tdnsping www.google.com. 127.0.0.1\nwith flags:")
	flag.PrintDefaults()
}

func main() {
	portFlag := flag.Int("p", 53, "`Port` to connect to")
	intervalFlag := flag.Duration("i", 1*time.Second, "How long to `wait` between requests")
	timeoutFlag := flag.Duration("t", 700*time.Millisecond, "`Timeout` for each query")
	countFlag := flag.Int("c", 0, "How many `requests` to make. Default is to run until ^C")
	queryTypeFlag := flag.String("q", "A", "Query `type` to use (A, AAAA, SOA, CNAME...)")
	versionFlag := flag.Bool("v", false, "Display version and exit.")
	// make logger be less about debug by default
	lcf := flag.Lookup("logcaller")
	lcf.DefValue = "false"
	_ = lcf.Value.Set("false")
	lpf := flag.Lookup("logprefix")
	lpf.DefValue = ""
	_ = lpf.Value.Set("")
	flag.CommandLine.Usage = usage
	flag.Parse()
	args := flag.Args()
	nArgs := len(args)
	log.LogVf("got %d arguments: %v", nArgs, args)
	if *versionFlag || (nArgs > 0 && args[0] == "version") {
		fmt.Println(version.Version)
		os.Exit(0)
	}
	qt, exists := dns.StringToType[strings.ToUpper(*queryTypeFlag)]
	if !exists {
		keys := []string{}
		for k := range dns.StringToType {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		log.Errf("Invalid -q type name %q, should be one of %v", *queryTypeFlag, keys)
		os.Exit(1)
	}
	if nArgs != 2 {
		fmt.Fprintf(os.Stderr, "Error: need exactly 2 arguments outside of the flags, got %d\n", nArgs)
		usage()
		os.Exit(1)
	}
	addrStr := fmt.Sprintf("%s:%d", args[1], *portFlag)
	query := args[0]
	if !strings.HasSuffix(query, ".") {
		query = query + "."
		log.LogVf("Adding missing . to query, now %q", query)
	}
	DNSPing(addrStr, query, qt, *countFlag, *intervalFlag, *timeoutFlag)
}

// DNSPing Runs the query howMany times against addrStr server, sleeping for interval time.
func DNSPing(addrStr, queryStr string, queryType uint16, howMany int, interval, timeout time.Duration) {
	m := new(dns.Msg)
	m.SetQuestion(queryStr, queryType)
	qtS := dns.TypeToString[queryType]
	howManyStr := fmt.Sprintf("%d times", howMany)
	if howMany <= 0 {
		howManyStr = "until interrupted"
	}
	log.Infof("Will query %s, sleeping %v in between, the server %s for %s (%d) record for %s",
		howManyStr, interval, addrStr, qtS, queryType, queryStr)
	log.LogVf("Query is: %v", m)
	successCount := 0
	errorCount := 0
	stats := stats.NewHistogram(0, 0.1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	continueRunning := true
	cli := dns.Client{Timeout: timeout}
	format := "%5.1f ms %3d: "
	start := time.Now()
	for i := 1; continueRunning && (howMany <= 0 || i <= howMany); i++ {
		if i != 1 {
			target := time.Duration(i-1) * interval
			elapsedSoFar := time.Since(start)
			waitFor := target - elapsedSoFar
			log.Debugf("target %v, elapsed so far %v -> wait for %v", target, elapsedSoFar, waitFor)
			select {
			case <-ch:
				continueRunning = false
				fmt.Println()
				continue
			case <-time.After(waitFor):
			}
		}
		r, duration, err := cli.Exchange(m, addrStr)
		durationMS := 1000. * duration.Seconds()
		stats.Record(durationMS)
		if err != nil {
			log.Errf(format+"failed call: %v", durationMS, i, err)
			errorCount++
			continue
		}
		if r == nil {
			log.Critf("bug? dns response is nil")
			errorCount++
			continue
		}
		log.LogVf("response is %v", r)
		if r.Rcode != dns.RcodeSuccess {
			log.Errf(format+"server error: %v", durationMS, i, err)
			errorCount++
		} else {
			successCount++
		}
		log.Infof(format+"%v", durationMS, i, r.Answer)
	}
	perc := fmt.Sprintf("%.02f%%", 100.*float64(errorCount)/float64(errorCount+successCount))
	fmt.Printf("%d errors (%s), %d success.\n", errorCount, perc, successCount)
	res := stats.Export()
	res.CalcPercentiles([]float64{50, 90, 99})
	res.Print(os.Stdout, "response time (in ms)")
}
