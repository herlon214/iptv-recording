FROM jrottenberg/ffmpeg:4.3-alpine

WORKDIR "/app"
COPY iptv-rec .
ENTRYPOINT ["/app/iptv-rec"]