FROM golang:1.26.4-alpine@sha256:7a3e50096189ad57c9f9f865e7e4aa8585ed1585248513dc5cda498e2f41812c AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.5.2@sha256:46d369a8379ee3cd52fca61147e5be818073896c25127407df2b2aca1ec10a9f . /

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


FROM alpine:3.24@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4

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
