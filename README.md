# dnsping
DNS Ping to check packet loss and latency issues with DNS servers

Sample run
```
dnsping -c 8 -i 750ms www.google.com. 8.8.8.8 
20:15:53 I main.go:51> Will query 8 times, sleeping 750ms in between, the server 8.8.8.8:53 for A (1) record for www.google.com.
20:15:53 I main.go:91>   46.7 ms   1: [www.google.com.	288	IN	A	172.217.6.36]
20:15:53 I main.go:91>   86.3 ms   2: [www.google.com.	255	IN	A	216.58.194.196]
20:15:54 I main.go:91>  115.9 ms   3: [www.google.com.	281	IN	A	216.58.194.196]
20:15:55 I main.go:91>  125.4 ms   4: [www.google.com.	43	IN	A	172.217.5.100]
20:15:56 I main.go:91>   10.3 ms   5: [www.google.com.	255	IN	A	216.58.194.196]
20:15:59 E main.go:75> 2001.6 ms   6: failed call: read udp 10.10.50.62:62634->8.8.8.8:53: i/o timeout
20:15:59 I main.go:91>   16.2 ms   7: [www.google.com.	277	IN	A	172.217.6.36]
20:16:00 I main.go:91>    8.3 ms   8: [www.google.com.	251	IN	A	216.58.194.196]
1 errors (12.50%), 7 success.
response time (in ms) : count 8 avg 301.32961 +/- 644.1 min 8.280587 max 2001.589235 sum 2410.6369
# range, mid point, percentile, count
>= 8.28059 <= 9 , 8.64029 , 12.50, 1
> 10 <= 12 , 11 , 25.00, 1
> 16 <= 18 , 17 , 37.50, 1
> 45 <= 50 , 47.5 , 50.00, 1
> 80 <= 90 , 85 , 62.50, 1
> 100 <= 200 , 150 , 87.50, 2
> 2000 <= 2001.59 , 2000.79 , 100.00, 1
# target 50% 50
# target 90% 2000.32
# target 99% 2001.46
```
