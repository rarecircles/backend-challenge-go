package main

import (
	"github.com/rarecircles/backend-challenge-go/logging"
	"go.uber.org/zap"
)

func getLogger(name string) *zap.Logger {
	zlog := logging.MustCreateLoggerWithServiceName(name)
	return zlog
}
