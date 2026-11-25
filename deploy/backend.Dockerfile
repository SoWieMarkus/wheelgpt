# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
ARG VERSION=dev
# modernc.org/sqlite is pure Go: a static binary without cgo.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/wheelgpt ./cmd/wheelgpt
# Empty data directory for the database, owned by the distroless nonroot
# user. Named volumes mounted here inherit this ownership.
RUN mkdir -p /out/data

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/wheelgpt /wheelgpt
COPY --from=build --chown=65532:65532 /out/data /data
ENV DATA_DIR=/data
VOLUME /data
EXPOSE 8080 9091
USER nonroot:nonroot
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s CMD ["/wheelgpt", "healthcheck"]
ENTRYPOINT ["/wheelgpt"]
