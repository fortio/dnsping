# dnsping
DNS Ping to check packet loss and latency issues with DNS servers

Sample run
```
dnsping www.google.com. 8.8.8.8
19:19:51 I main.go:40> Will query server: 8.8.8.8:53 for A (1) record for www.google.com.
19:19:51 I main.go:70>   12.3 ms   1: [www.google.com.	297	IN	A	172.217.5.100]
19:19:52 I main.go:70>   15.0 ms   2: [www.google.com.	72	IN	A	216.58.194.196]
19:19:53 I main.go:70>    8.5 ms   3: [www.google.com.	285	IN	A	172.217.5.100]
19:19:56 E main.go:54> 2005.2 ms   4: failed call: read udp 10.10.50.62:63568->8.8.8.8:53: i/o timeout
19:19:59 E main.go:54> 2001.1 ms   5: failed call: read udp 10.10.50.62:52814->8.8.8.8:53: i/o timeout
19:20:00 I main.go:70>    8.7 ms   6: [www.google.com.	154	IN	A	172.217.5.100]
19:20:01 I main.go:70>    9.2 ms   7: [www.google.com.	191	IN	A	216.58.194.196]
19:20:02 I main.go:70>   15.6 ms   8: [www.google.com.	286	IN	A	172.217.5.100]
19:20:03 I main.go:70>   15.9 ms   9: [www.google.com.	64	IN	A	172.217.5.100]
19:20:06 E main.go:54> 2005.2 ms  10: failed call: read udp 10.10.50.62:64919->8.8.8.8:53: i/o timeout
3 errors (30.00%), 7 success.
response time (in ms) : count 10 avg 609.695 +/- 912.7 min 8.548663 max 2005.2117090000002 sum 6096.94996
# range, mid point, percentile, count
>= 8.54866 <= 9 , 8.77433 , 20.00, 2
> 9 <= 10 , 9.5 , 30.00, 1
> 12 <= 14 , 13 , 40.00, 1
> 14 <= 16 , 15 , 70.00, 3
> 2000 <= 2005.21 , 2002.61 , 100.00, 3
# target 50% 14.6667
# target 90% 2003.47
# target 99% 2005.04
```
