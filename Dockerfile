FROM golang:1.14.7-alpine AS builder

COPY . /root/app

RUN cd /root/app && go install .

FROM alpine:3.12

RUN apk add --update --no-cache ca-certificates \
 && rm -rf /var/cache/apk/*

COPY --from=builder /go/bin/lang-tool-versions /usr/local/bin/

WORKDIR /root
