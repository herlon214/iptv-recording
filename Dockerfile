FROM alpine:3.13.6
WORKDIR "/app"
COPY iptv-rec .
ENTRYPOINT ["/app/iptv-rec"]