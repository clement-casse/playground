FROM golang:1.22 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=$GOPKG/pkg \
    go mod download && go mod verify

COPY main.go main.go

RUN --mount=type=cache,target=$GOPKG/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -v -o /usr/local/bin/entrypoint .


FROM alpine:latest

RUN apk update \
 && apk add ca-certificates iptables ip6tables \
 && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy Tailscale binaries from the tailscale image on Docker Hub.
COPY --from=docker.io/tailscale/tailscale:stable \
     --chown=nobody:nobody \
     --chmod=700 \
     /usr/local/bin/tailscaled /app/tailscaled

# Copy blocky binary from its own image on Docker Hub.
COPY --from=docker.io/spx01/blocky:latest \
     --chown=nobody:nobody \
     --chmod=700 \
     /app/blocky /app/blocky


COPY --from=builder \
     --chown=nobody:nobody \
     --chmod=700 \
     /usr/local/bin/entrypoint entrypoint

COPY --chown=nobody:nobody \
     --chmod=400 \
     ./blocky/config.yml config.yml

EXPOSE 53 443 4000

ENTRYPOINT ["/app/entrypoint"]

HEALTHCHECK --start-period=1m --timeout=3s CMD ["/app/blocky", "healthcheck"]