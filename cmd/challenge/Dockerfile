FROM golang:1.16
RUN apt update -y && apt upgrade -y

COPY . ./app
RUN mkdir /home/app

WORKDIR app

RUN env GO111MODULE=on go get github.com/cortesi/modd/cmd/modd
RUN go mod vendor