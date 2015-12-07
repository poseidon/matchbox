FROM busybox:latest
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
ADD bin/server /bin/server

EXPOSE 8081
CMD ./bin/server
