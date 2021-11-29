FROM golang:1.16.10-alpine3.14 AS base

ENV CGO_ENABLED 0
ENV GOPROXY https://proxy.golang.org

RUN apk update && apk add bash

SHELL ["/bin/sh", "-o", "pipefail", "-c"]
#-------------DEPENDENCIES-----------
FROM base AS deps

COPY go.mod go.sum $GOPATH/src/github.com/jose-camilo/backend-challenge-go/
WORKDIR $GOPATH/src/github.com/jose-camilo/backend-challenge-go
COPY . $GOPATH/src/github.com/jose-camilo/backend-challenge-go


#------------BUILD----------------------------
FROM deps AS build
WORKDIR $GOPATH/src/github.com/jose-camilo/backend-challenge-go/
RUN go build -a -o /challenge cmd/challenge/main.go


#------------APP------------------------------
FROM busybox:1.30 AS app
WORKDIR /
COPY data/addresses.jsonl /data/addresses.jsonl
COPY --from=build /challenge .
CMD ["./challenge"]
