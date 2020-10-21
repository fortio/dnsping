# dnsping
DNS Ping to check packet loss and latency issues with DNS servers

If you have golang, easiest install is `go get -u fortio.org/dnsping`

Usage:
`dnsping [flags] query server`

```Shell
$ dnsping -h
Usage:	dnsping [flags] query server
eg:	dnsping www.google.com. 127.0.0.1
With flags:
  -c requests
    	How many requests to make. Default is to run until ^C
  -i wait
    	How long to wait between requests (default 1s)
  -loglevel value
    	loglevel, one of [Debug Verbose Info Warning Error Critical Fatal] (default Info)
  -p Port
    	Port to connect to (default 53)
  -q type
    	Query type to use (A, AAAA, SOA, CNAME...) (default "A")
  -t Timeout
    	Timeout for each query (default 700ms)
```

Sample run
```
dnsping  -c 8 www.google.com. 8.8.8.8
Wed Oct 21 16:04:09 PDT 2020
16:04:10 I Will query 8 times, sleeping 1s in between, the server 8.8.8.8:53 for A (1) record for www.google.com.
16:04:10 I   8.6 ms   1: [www.google.com.	179	IN	A	172.217.5.100]
16:04:11 I  23.1 ms   2: [www.google.com.	290	IN	A	172.217.0.36]
16:04:12 I  15.7 ms   3: [www.google.com.	293	IN	A	172.217.164.100]
16:04:13 I   9.9 ms   4: [www.google.com.	159	IN	A	172.217.6.36]
16:04:14 E 700.1 ms   5: failed call: read udp 10.10.50.62:59250->8.8.8.8:53: i/o timeout
16:04:15 I  24.5 ms   6: [www.google.com.	255	IN	A	172.217.6.68]
16:04:16 I  12.4 ms   7: [www.google.com.	179	IN	A	172.217.6.36]
16:04:17 I  17.8 ms   8: [www.google.com.	260	IN	A	172.217.164.100]
1 errors (12.50%), 7 success.
response time (in ms) : count 8 avg 101.49609 +/- 226.3 min 8.553971 max 700.101971 sum 811.968711
# range, mid point, percentile, count
>= 8.55397 <= 9 , 8.77699 , 12.50, 1
> 9 <= 10 , 9.5 , 25.00, 1
> 12 <= 14 , 13 , 37.50, 1
> 14 <= 16 , 15 , 50.00, 1
> 16 <= 18 , 17 , 62.50, 1
> 20 <= 25 , 22.5 , 87.50, 2
> 500 <= 700.102 , 600.051 , 100.00, 1
# target 50% 16
# target 90% 540.02
# target 99% 684.094
```

Made thanks to https://github.com/miekg/dns (and https://github.com/fortio/fortio stats and logger)
