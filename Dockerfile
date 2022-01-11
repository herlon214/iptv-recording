# Build
FROM golang:1.17 as build
ENV GOPATH=/go
WORKDIR /go/src/github.com/herlon214/iptv-recording
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -extldflags '-static'" -o /app/iptv-rec main.go

# Deploy
FROM alpine:3.13.6
WORKDIR "/app"
COPY --from=build /app/iptv-rec .
ENTRYPOINT ["/app/iptv-rec"]