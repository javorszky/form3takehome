FROM golang:1.15.6-alpine3.12

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /go/src/github.com/javorszky/form3takehome/accountsclient/

CMD CGO_ENABLED=0 go test -count=1 -cover -v ./...
