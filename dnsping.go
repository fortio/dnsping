// Copyright 2020 Fortio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"fortio.org/cli"
	"fortio.org/fortio/stats"
	"fortio.org/log"
	"fortio.org/safecast"
	"github.com/miekg/dns"
)

// DNSPingConfig is the input configuration for DNSPing().
type DNSPingConfig struct {
	Server        string        // Server to send query to
	Query         string        // Query to send
	HowMany       int           // How many requests to send (0 for until interrupted)
	Interval      time.Duration // Target interval at which to repeat requests
	Timeout       time.Duration // Total timeout
	FixedID       int           // non zero is the message id to use for all requests
	QueryType     uint16        // Type of query (dns.Type)
	SequentialIDs bool          // true means sequential instead of random ids (assuming FixedID is 0)
	Recursion     bool          // DNS recursion requested or not
	TCP           bool          // Use TCP instead of UDP
	EDNS          bool          // Use EDNS0
	DNSSECOK      bool          // Set the DNSSEC OK bit
}

// DNSPingResults is the aggregated results of the DNSPing() call including input. Ready for JSON serialization.
type DNSPingResults struct {
	Config  *DNSPingConfig
	Errors  int
	Success int
	Stats   *stats.HistogramData
}

func main() {
	os.Exit(Main())
}

func Main() int {
	jsonFlag := flag.String("json", "", "Json output to provided file `path` or '-' for stdout (empty = no json output)")
	portFlag := flag.Int("p", 53, "`Port` to connect to")
	intervalFlag := flag.Duration("i", 1*time.Second, "How long to `wait` between requests")
	timeoutFlag := flag.Duration("t", 700*time.Millisecond, "`Timeout` for each query")
	countFlag := flag.Int("c", 0, "How many `requests` to make. Default is to run until ^C")
	queryTypeFlag := flag.String("q", "A", "Query `type` to use (A, AAAA, SOA, CNAME...)")
	seqIDFlag := flag.Bool("sequential-id", false, "Use sequential ids instead of random.")
	sameIDFlag := flag.Int("fixed-id", 0, "Non 0 id to use instead of random or sequential")
	recursionFlag := flag.Bool("no-recursion", false, "Pass to disable (default) recursion.")
	ednsFlag := flag.Bool("edns", false, "Set EDNS0 4096")
	dnssecFlag := flag.Bool("dnssec", false, "Set the DNSSEC OK bit")
	tcpFlag := flag.Bool("tcp", false, "Use TCP instead of normal UDP")
	cli.MinArgs = 2
	cli.ArgsHelp = "query server\neg:\tdnsping www.google.com. 8.8.8.8"
	cli.Main()
	qt, exists := dns.StringToType[strings.ToUpper(*queryTypeFlag)]
	if !exists {
		keys := []string{}
		for k := range dns.StringToType {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return log.FErrf("Invalid -q type name %q, should be one of %v", *queryTypeFlag, keys)
	}
	server := flag.Arg(1)
	if strings.Contains(server, ":") && !strings.HasPrefix(server, "[") {
		server = "[" + server + "]"
		log.Infof("Adding [] around detected input IPv6 server ip info: %s", server)
	}
	addrStr := fmt.Sprintf("%s:%d", server, *portFlag)
	query := flag.Arg(0)
	if !strings.HasSuffix(query, ".") {
		query += "."
		log.LogVf("Adding missing . to query, now %q", query)
	}
	cfg := DNSPingConfig{
		Server:        addrStr,
		Query:         query,
		QueryType:     qt,
		HowMany:       *countFlag,
		Interval:      *intervalFlag,
		Timeout:       *timeoutFlag,
		SequentialIDs: *seqIDFlag,
		FixedID:       *sameIDFlag,
		Recursion:     !*recursionFlag,
		TCP:           *tcpFlag,
		EDNS:          *ednsFlag,
		DNSSECOK:      *dnssecFlag,
	}
	r := DNSPing(&cfg)
	if *jsonFlag == "" {
		return 0
	}
	return JSONSave(r, *jsonFlag)
}

// JSONSave exports a result into a json file (or stdpout if - is passed).
// TODO refactor from fortio's main.
func JSONSave(res *DNSPingResults, jsonFileName string) int {
	var j []byte
	j, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return log.FErrf("Unable to json serialize result: %v", err)
	}
	var f *os.File
	if jsonFileName == "-" {
		f = os.Stdout
		jsonFileName = "stdout"
	} else {
		f, err = os.Create(jsonFileName)
		if err != nil {
			return log.FErrf("Unable to create %s: %v", jsonFileName, err)
		}
	}
	n, err := f.Write(append(j, '\n'))
	if err != nil {
		return log.FErrf("Unable to write json to %s: %v", jsonFileName, err)
	}
	if f != os.Stdout {
		err := f.Close()
		if err != nil {
			return log.FErrf("Close error for %s: %v", jsonFileName, err)
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "Successfully wrote %d bytes of Json data to %s\n", n, jsonFileName)
	return 0
}

// DNSPing Runs the query howMany times against addrStr server, sleeping for interval time.
func DNSPing(cfg *DNSPingConfig) *DNSPingResults { //nolint:funlen // yes it's long
	m := new(dns.Msg)
	if cfg.EDNS || cfg.DNSSECOK {
		m.SetEdns0(4096, cfg.DNSSECOK)
	}
	m.SetQuestion(cfg.Query, cfg.QueryType)
	m.RecursionDesired = cfg.Recursion
	qtS := dns.TypeToString[cfg.QueryType]
	howMany := cfg.HowMany
	howManyStr := fmt.Sprintf("%d times", howMany)
	if howMany <= 0 {
		howManyStr = "until interrupted"
	}
	log.Infof("dnsping %s: will query %s, sleeping %v in between, the server %s for %s (%d) record for %s",
		cli.ShortVersion, howManyStr, cfg.Interval, cfg.Server, qtS, cfg.QueryType, cfg.Query)
	log.LogVf("Query is: %v", m)
	successCount := 0
	errorCount := 0
	stats := stats.NewHistogram(0, 0.1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	continueRunning := true
	cli := dns.Client{Timeout: cfg.Timeout}
	if cfg.TCP {
		log.Infof("Using TCP instead of UDP for DNS queries")
		cli.Net = "tcp"
	}
	format := "%5.1f ms %3d: "
	start := time.Now()
	for i := 1; continueRunning && (howMany <= 0 || i <= howMany); i++ {
		if i != 1 {
			target := time.Duration(i-1) * cfg.Interval
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
		switch {
		case cfg.FixedID != 0:
			m.Id = safecast.MustConvert[uint16](cfg.FixedID)
		case cfg.SequentialIDs:
			m.Id = uint16(i)
		default:
			m.Id = dns.Id()
		}
		r, duration, err := cli.Exchange(m, cfg.Server)
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
			continue
		}
		successCount++
		extra := ""
		if r.Truncated {
			extra = " (truncated: use tcp)"
		}
		if edns := r.IsEdns0(); edns != nil {
			dnssec := ""
			if edns.Do() {
				dnssec = " and DNSSEC OK"
			}
			extra += fmt.Sprintf(" (with EDNS0 %d%s)", edns.UDPSize(), dnssec)
		}
		log.Printf(format+"%v%s (%d bytes)", durationMS, i, r.Answer, extra, r.Len())
	}
	perc := fmt.Sprintf("%.02f%%", 100.*float64(errorCount)/float64(errorCount+successCount))
	plural := "s" // 0 errors 1 error 2 errors...
	if errorCount == 1 {
		plural = ""
	}
	fmt.Printf("%d error%s (%s), %d success.\n", errorCount, plural, perc, successCount)
	res := stats.Export()
	res.CalcPercentiles([]float64{50, 90, 99})
	res.Print(os.Stdout, "response time (in ms)")
	return &DNSPingResults{
		Config:  cfg,
		Errors:  errorCount,
		Success: successCount,
		Stats:   res,
	}
}
