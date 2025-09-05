FROM golang:1.25-alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS builder

COPY . /go/src/github.com/Luzifer/go-latestver
WORKDIR /go/src/github.com/Luzifer/go-latestver

RUN set -ex \
 && apk add --update \
      git \
      make \
      nodejs-current \
 && make frontend_prod \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly


FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

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
