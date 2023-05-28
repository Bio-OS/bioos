FROM golang:1.19.8-bullseye AS builder
ENV GOPROXY=https://goproxy.cn,direct
ARG TARGETOS
ARG TARGETARCH
WORKDIR /go/src/github.com/Bio-OS/bioos
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN BINS=bioctl PLATFORM=${TARGETOS}_${TARGETARCH} make go.build

FROM debian:bullseye
ARG TARGETOS
ARG TARGETARCH
USER    nobody
WORKDIR /app
COPY --from=builder /go/src/github.com/Bio-OS/bioos/_output/platforms/${TARGETOS}/${TARGETARCH}/bioctl .
ENTRYPOINT ["sleep", "1d"]
