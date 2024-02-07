FROM docker.io/golang:1.22.0 AS builder
COPY . src
RUN cd src && make build

FROM docker.io/alpine:3.19.1
LABEL maintainer="Dalton Hubble <dghubble@gmail.com>"
LABEL org.opencontainers.image.title="Matchbox",
LABEL org.opencontainers.image.source="https://github.com/poseidon/matchbox"
LABEL org.opencontainers.image.documentation="https://matchbox.psdn.io"
LABEL org.opencontainers.image.vendor="Poseidon Labs"
COPY --from=builder /go/src/bin/matchbox /matchbox
EXPOSE 8080
ENTRYPOINT ["/matchbox"]
