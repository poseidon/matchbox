FROM alpine:3.9
LABEL maintainer="Dalton Hubble <dghubble@gmail.com>"
# ca-certificates needed in case s3 is used as storage backend
RUN apk add --no-cache ca-certificates
COPY bin/matchbox /matchbox
EXPOSE 8080
ENTRYPOINT ["/matchbox"]
