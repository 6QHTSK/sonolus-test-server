# sonolus-test-server Builder
#
# VERSION 0.7.0-rc.1

FROM golang:1.20 as builder
MAINTAINER 6QHTSK <psk2019@qq.com>

ENV GO111MODULE=on

WORKDIR /go/src/sonolus-test-server
COPY . /go/src/sonolus-test-server

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sonolus-test-server .

# sonolus-test-server
#
# VERSION 0.7.0-rc.1
FROM alpine:latest

MAINTAINER 6QHTSK <psk2019@qq.com>

ENV GIN_MODE=release
ENV PORT=8000
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /sonolus-test-server

COPY --from=builder /go/src/sonolus-test-server/sonolus-test-server .

EXPOSE 8000

ENTRYPOINT ["./sonolus-test-server"]