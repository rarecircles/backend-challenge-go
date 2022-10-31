FROM golang:1.19-alpine3.16 AS base

ENV CGO_ENABLED 0
ENV GOPROXY https://proxy.golang.org

RUN apk update && apk add bash

SHELL ["/bin/sh", "-o", "pipefail", "-c"]
#-------------DEPENDENCIES-----------
FROM base AS deps

COPY go.mod go.sum $GOPATH/src/github.com/degarajesh/backend-challenge-go/
WORKDIR $GOPATH/src/github.com/degarajesh/backend-challenge-go
COPY . $GOPATH/src/github.com/degarajesh/backend-challenge-go


#------------BUILD----------------------------
FROM deps AS build
WORKDIR $GOPATH/src/github.com/degarajesh/backend-challenge-go/
RUN go build -a -o /challenge cmd/challenge/main.go


#------------APP------------------------------
FROM busybox:1.30 AS app
WORKDIR /
COPY data/addresses.jsonl /data/addresses.jsonl
COPY app.env .
COPY --from=build /challenge .
CMD ["./challenge"]