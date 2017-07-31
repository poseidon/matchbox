FROM alpine:3.6
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
RUN apk -U add dnsmasq curl
COPY tftpboot /var/lib/tftpboot
EXPOSE 53 67 69
ENTRYPOINT ["/usr/sbin/dnsmasq"]
