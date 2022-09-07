FROM docker.io/golang:1.19.1 AS builder
COPY . src
RUN cd src && make build

FROM docker.io/alpine:3.16.2
LABEL maintainer="Dalton Hubble <dghubble@gmail.com>"
COPY --from=builder /go/src/bin/matchbox /matchbox
EXPOSE 8080
ENTRYPOINT ["/matchbox"]
