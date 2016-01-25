FROM alpine:latest
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
RUN apk -U add dnsmasq curl
RUN mkdir -p /var/lib/tftpboot
RUN curl -s -o /var/lib/tftpboot/undionly.kpxe http://boot.ipxe.org/undionly.kpxe
RUN ln -s /var/lib/tftpboot/undionly.kpxe /var/lib/tftpboot/undionly.kpxe.0
EXPOSE 53
ENTRYPOINT ["/usr/sbin/dnsmasq", "-d"]