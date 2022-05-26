# DNSping
[![PkgGoDev](https://pkg.go.dev/badge/fortio.org/dnsping?tab=overview)](https://pkg.go.dev/fortio.org/dnsping?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/fortio.org/dnsping)](https://goreportcard.com/report/fortio.org/dnsping)
[![Docker Build](https://img.shields.io/docker/cloud/build/fortio/dnsping.svg)](https://hub.docker.com/r/fortio/dnsping)
[![Docker Pulls](https://img.shields.io/docker/pulls/fortio/dnsping.svg)](https://hub.docker.com/r/fortio/dnsping)

DNS Ping checks packet loss and latency issues with DNS servers

## Installation

If you have golang, easiest install is `go install fortio.org/dnsping@latest`

Or with docker `docker run fortio/dnsping ...`

Or brew custom tap source build `brew install fortio/dnsping/dnsping` (please star the project so it can go in core and get binary bottles built)

Otherwise head over to https://github.com/fortio/dnsping/releases for binary releases

## Usage:
`dnsping [flags] query server`

```Shell
$ dnsping -h
dnsping v1.1.5 usage:
	dnsping [flags] query server
eg:	dnsping www.google.com. 127.0.0.1
with flags:
  -c requests
    	How many requests to make. Default is to run until ^C
  -fixed-id int
    	Non 0 id to use instead of random or sequential
  -i wait
    	How long to wait between requests (default 1s)
  -json path
    	Json output to provided file path or '-' for stdout (empty = no json output)
  -loglevel value
    	loglevel, one of [Debug Verbose Info Warning Error Critical Fatal] (default Info)
  -no-recursion
    	Pass to disable (default) recursion.
  -p Port
    	Port to connect to (default 53)
  -q type
    	Query type to use (A, AAAA, SOA, CNAME...) (default "A")
  -sequential-id
    	Use sequential ids instead of random.
  -t Timeout
    	Timeout for each query (default 700ms)
  -v	Display version and exit.
```

Sample run
```
dnsping -fixed-id 42 -json sampleResult.json -c 8  www.google.com 8.8.4.4
16:08:03 I Will query 8 times, sleeping 1s in between, the server 8.8.4.4:53 for A (1) record for www.google.com.
16:08:03 I   8.7 ms   1: [www.google.com.	298	IN	A	172.217.6.68]
16:08:04 I  16.5 ms   2: [www.google.com.	229	IN	A	172.217.6.36]
16:08:05 I  14.1 ms   3: [www.google.com.	179	IN	A	216.58.194.196]
16:08:06 E 700.3 ms   4: failed call: read udp 10.10.50.62:65456->8.8.4.4:53: i/o timeout
16:08:07 I  15.0 ms   5: [www.google.com.	195	IN	A	216.58.194.196]
16:08:08 I  13.5 ms   6: [www.google.com.	196	IN	A	216.58.194.196]
16:08:09 I  14.8 ms   7: [www.google.com.	179	IN	A	216.58.194.196]
16:08:10 I  15.5 ms   8: [www.google.com.	285	IN	A	172.217.6.68]
1 error (12.50%), 7 success.
response time (in ms) : count 8 avg 99.792926 +/- 227 min 8.684216 max 700.257965 sum 798.343406
# range, mid point, percentile, count
>= 8.68422 <= 9 , 8.84211 , 12.50, 1
> 12 <= 14 , 13 , 25.00, 1
> 14 <= 16 , 15 , 75.00, 4
> 16 <= 18 , 17 , 87.50, 1
> 500 <= 700.258 , 600.129 , 100.00, 1
# target 50% 15
# target 90% 540.052
# target 99% 684.237
Successfully wrote 1212 bytes of Json data to sampleResult.json
```

Which also produces the json:
```Json
{
  "Config": {
    "Server": "8.8.4.4:53",
    "Query": "www.google.com.",
    "HowMany": 8,
    "Interval": 1000000000,
    "Timeout": 700000000,
    "FixedID": 42,
    "QueryType": 1,
    "SequentialIDs": false,
    "Recursion": true
  },
  "Errors": 1,
  "Success": 7,
  "Stats": {
    "Count": 8,
    "Min": 8.684216,
    "Max": 700.257965,
    "Sum": 798.3434060000001,
    "Avg": 99.79292575000001,
    "StdDev": 226.96508473843934,
    "Data": [
      {
        "Start": 8.684216,
        "End": 9,
        "Percent": 12.5,
        "Count": 1
      },
      {
        "Start": 12,
        "End": 14,
        "Percent": 25,
        "Count": 1
      },
      {
        "Start": 14,
        "End": 16,
        "Percent": 75,
        "Count": 4
      },
      {
        "Start": 16,
        "End": 18,
        "Percent": 87.5,
        "Count": 1
      },
      {
        "Start": 500,
        "End": 700.257965,
        "Percent": 100,
        "Count": 1
      }
    ],
    "Percentiles": [
      {
        "Percentile": 50,
        "Value": 15
      },
      {
        "Percentile": 90,
        "Value": 540.051593
      },
      {
        "Percentile": 99,
        "Value": 684.2373278
      }
    ]
  }
}
```

Made thanks to https://github.com/miekg/dns (and https://github.com/fortio/fortio stats and logger)
