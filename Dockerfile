FROM docker.io/alpine:3.12
LABEL maintainer="Dalton Hubble <dghubble@gmail.com>"
COPY bin/matchbox /matchbox
EXPOSE 8080
ENTRYPOINT ["/matchbox"]
