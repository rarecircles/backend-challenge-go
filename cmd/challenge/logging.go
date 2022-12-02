package main

import (
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/data"
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/http"
	endpoint "github.com/rarecircles/backend-challenge-go/cmd/challenge/http/endpoints"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/logging"
	"go.uber.org/zap"
)

var zlog *zap.Logger

func init() {
	zlog = logging.MustCreateLoggerWithServiceName("challenge")
	rpc.SetLogger(zlog)
	eth.SetLogger(zlog)
	endpoint.SetLogger(zlog)
	data.SetLogger(zlog)
	http.SetLogger(zlog)
}
