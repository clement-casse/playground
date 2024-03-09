FROM golang:1.22 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=$GOPKG/pkg \
    go mod download && go mod verify

COPY . .
ARG mod="github.com/clement-casse/playground/webservice-go"
ARG version="develop"
ARG builddate

RUN --mount=type=cache,target=$GOPKG/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -v \
      -o /usr/local/bin/my-service \
      -ldflags="-X ${mod}/cmd/internal.version=${version} -X ${mod}/cmd/internal.date=${builddate}" \
      ./cmd/my-service

#-----
FROM gcr.io/distroless/static-debian12

WORKDIR /
COPY --from=builder /usr/local/bin/my-service /my-service

USER nonroot:nonroot
ENTRYPOINT ["/my-service"]
