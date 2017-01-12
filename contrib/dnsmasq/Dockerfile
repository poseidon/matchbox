FROM alpine:3.5
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
RUN apk -U add dnsmasq curl
COPY tftpboot /var/lib/tftpboot
EXPOSE 53
ENTRYPOINT ["/usr/sbin/dnsmasq"]