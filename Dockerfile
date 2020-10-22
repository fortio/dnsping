# Build the binaries in larger image
FROM docker.io/fortio/fortio.build:v27 as build
WORKDIR /build
COPY . dnsping
# We moved a lot of the logic into the Makefile so it can be reused in brew
RUN make -C dnsping official-build-version OFFICIAL_BIN=/out/dnsping.linux
RUN make -C dnsping official-build OFFICIAL_BIN=/out/dnsping.mac GOOS=darwin
RUN make -C dnsping official-build BUILD_DIR=/build LIB_DIR=. OFFICIAL_BIN=/out/dnsping.exe GOOS=windows
# Minimal image with just the binary and certs
RUN ls -lh /out
FROM scratch as release
# NOTE: the list of files here, if updated, must be changed in release/Dockerfile.in too
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/dnsping.linux /usr/bin/dnsping
ENTRYPOINT ["/usr/bin/dnsping"]
CMD ["-h"]
