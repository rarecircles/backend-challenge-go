package main

import (
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
}
