FROM golang:1.17.5 AS GO_BUILD
ENV CGO_ENABLED 0
COPY . /go-app
WORKDIR /go-app
RUN go build -o server
FROM alpine:3.15
WORKDIR /go-app
COPY --from=GO_BUILD /go-app/server /go-app/server
EXPOSE 8080
CMD ["./server"]