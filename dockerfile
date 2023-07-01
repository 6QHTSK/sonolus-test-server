# sonolus-test-server Builder
#
# VERSION 0.7.0-rc.3

FROM golang:1.20 as builder
MAINTAINER 6QHTSK <psk2019@qq.com>

ENV GO111MODULE=on

WORKDIR /go/src/sonolus-test-server
COPY . /go/src/sonolus-test-server

RUN GOOS=linux GOARCH=amd64 go build -o sonolus-test-server .

# sonolus-test-server
#
# VERSION 0.7.0-rc.3
FROM ubuntu:latest

MAINTAINER 6QHTSK <psk2019@qq.com>

ENV GIN_MODE=release
ENV PORT=8000
# 安装ca-certificates来进行http服务
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl

# 安装tzdata来进行时区转化服务
ENV TZ=Asia/Shanghai
RUN echo "${TZ}" > /etc/timezone \
  && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
  && apt-get -qq install tzdata \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /sonolus-test-server

COPY --from=builder /go/src/sonolus-test-server/sonolus-test-server .

EXPOSE 8000

ENTRYPOINT ["./sonolus-test-server"]