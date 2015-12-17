FROM alpine:latest
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
RUN apk -U add dnsmasq
EXPOSE 53
ENTRYPOINT ["dnsmasq", "-d"]
CMD ["--dhcp-range=172.17.0.43,172.17.0.99,30m"]
