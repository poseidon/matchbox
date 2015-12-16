FROM alpine:latest
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
RUN apk -U add dnsmasq curl
COPY dnsmasq.conf /etc/dnsmasq.conf
RUN mkdir -p /var/lib/tftpboot
RUN curl -s -o /var/lib/tftpboot/undionly.kpxe http://boot.ipxe.org/undionly.kpxe
EXPOSE 53 53/udp 67/udp
ENTRYPOINT ["dnsmasq", "-d"]