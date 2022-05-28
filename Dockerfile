FROM scratch
COPY dnsping /usr/bin/dnsping
ENTRYPOINT ["/usr/bin/dnsping"]
CMD ["-h"]
