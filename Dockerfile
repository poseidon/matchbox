FROM alpine:3.9
LABEL maintainer="Dalton Hubble <dghubble@gmail.com>"
RUN apk update && apk add ca-certificates
COPY bin/matchbox /matchbox
EXPOSE 8080
ENTRYPOINT ["/matchbox"]
