FROM golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.12.0@sha256:9ae65fe485d12e216f537c5bf16054ad785ac14fba2425b2452bcb9ab501c15b . /

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
