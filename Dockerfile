FROM golang:1.25-alpine@sha256:06cdd34bd531b810650e47762c01e025eb9b1c7eadd191553b91c9f2d549fae8 AS builder

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


FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

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
