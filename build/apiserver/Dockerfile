FROM golang:1.19.8-bullseye AS builder
ARG TARGETOS
ARG TARGETARCH
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /go/src/github.com/Bio-OS/bioos
COPY go.mod .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN go env -w CGO_ENABLED=1 && BINS=apiserver PLATFORM=${TARGETOS}_${TARGETARCH} make go.build && make tools.install.cfssl && make tools.install.womtool && make generate.certs

FROM debian:bullseye
RUN apt update \
     && apt install -y --no-install-recommends ca-certificates openjdk-11-jre-headless \
     && apt clean \
     && rm -rf /var/lib/apt/lists/*
ARG TARGETOS
ARG TARGETARCH
USER    nobody
WORKDIR /app
COPY --from=builder --chown=nobody:nogroup /go/src/github.com/Bio-OS/bioos/conf conf
COPY --from=builder /go/src/github.com/Bio-OS/bioos/_output/platforms/${TARGETOS}/${TARGETARCH}/apiserver .
COPY --from=builder /go/src/github.com/Bio-OS/bioos/womtool.jar .
ENTRYPOINT ["/app/apiserver"]
