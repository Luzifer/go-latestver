FROM golang:1.27.0-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v12.1.0@sha256:c5ab77a93cb7faeca4b5c5477114d2ceded271adfbc609986c5901d9d3aebb79 . /

COPY . /go/src/github.com/Luzifer/go-latestver
WORKDIR /go/src/github.com/Luzifer/go-latestver

RUN set -ex \
 && apk add --update \
      git \
      make \
      nodejs \
 && make frontend_prod \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly


FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

LABEL org.opencontainers.image.source="https://github.com/Luzifer/go-latestver" \
      org.opencontainers.image.name="go-latestver" \
      org.opencontainers.image.description="Monitor Software Versions in a single place" \
      org.opencontainers.image.authors="Knut Ahlers <knut@ahlers.me>" \
      org.opencontainers.image.url="https://github.com/Luzifer/go-latestver" \
      org.opencontainers.image.licenses="Apache-2.0"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/go-latestver /usr/local/bin/go-latestver

EXPOSE 3000
VOLUME ["/data"]

ENTRYPOINT ["/usr/local/bin/go-latestver"]
CMD ["--"]

# vim: set ft=Dockerfile:
