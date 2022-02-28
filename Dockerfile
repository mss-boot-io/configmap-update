FROM golang:alpine as builder

MAINTAINER lwnmengjing <991154416@qq.com>

#ENV GOPROXY https://goproxy.io/

WORKDIR /go/release
RUN apk update && apk add tzdata && apk add curl unzip procps ca-certificates

COPY go.mod ./go.mod
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o configmap-update .

FROM alpine

COPY --from=builder /go/release/configmap-update /

COPY --from=builder /go/release/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]