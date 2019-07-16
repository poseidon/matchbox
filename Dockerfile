FROM alpine:3.10
LABEL maintainer="Dalton Hubble <dghubble@gmail.com>"
COPY bin/matchbox /matchbox
EXPOSE 8080
ENTRYPOINT ["/matchbox"]
