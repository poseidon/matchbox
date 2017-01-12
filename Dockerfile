FROM alpine:3.5
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
COPY bin/matchbox /matchbox
EXPOSE 8080
ENTRYPOINT ["/matchbox"]
