package logging

import "go.uber.org/zap"

var zLog = zap.NewNop()

func SetLogger(logger *zap.Logger) {
	zLog = logger
}
