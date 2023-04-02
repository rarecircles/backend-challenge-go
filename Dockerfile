FROM golang:1.20 AS GO_BUILD
ENV CGO_ENABLED 0
WORKDIR /go-app
COPY . .
RUN go build -o server
FROM alpine:3.15
WORKDIR /go-app
COPY data/addresses_i.jsonl data/addresses_i.jsonl
COPY --from=GO_BUILD /go-app/server /go-app/server
EXPOSE 10000
CMD ["./server"]