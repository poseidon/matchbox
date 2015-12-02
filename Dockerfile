FROM busybox:latest
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
ADD bin/server /bin/server

EXPOSE 8080
CMD ./bin/server
