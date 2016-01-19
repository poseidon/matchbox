FROM alpine:latest
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
COPY bin/bootcfg /bootcfg
EXPOSE 8080
ENTRYPOINT ["./bootcfg"]
