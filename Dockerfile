FROM alpine:latest
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
COPY bin/server /server
EXPOSE 8080
ENTRYPOINT ["./server"]
