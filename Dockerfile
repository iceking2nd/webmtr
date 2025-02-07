FROM golang:alpine AS builder
ARG VERSION="0.0.0"
ARG BUILD_TIME="Thu Jan 01 1970 00:00:00 GMT+0000"
ARG SHA="e15d5c7bc909a54f53f98c8984a8d321"

COPY . /go/src/github.com/iceking2nd/webmtr
WORKDIR /go/src/github.com/iceking2nd/webmtr

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-extldflags \"-static\" -X 'github.com/iceking2nd/webmtr/global.Version=${VERSION}-docker' -X 'github.com/iceking2nd/webmtr/global.BuildTime=${BUILD_TIME}' -X github.com/iceking2nd/webmtr/global.GitCommit=${SHA}" \
    -o /bin/lsf
RUN go test -cover -v ./...

FROM alpine:latest
COPY --from=builder /bin/webmtr /bin/webmtr
ENTRYPOINT ["/bin/webmtr"]