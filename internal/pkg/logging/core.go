package logging

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var AutoStartServer = true

func MustCreateLoggerWithServiceName(serviceName string, opts ...zap.Option) *zap.Logger {
	return mustCreateLoggerWithLevel(serviceName, levelFromEnvironment(), opts...)
}

func mustCreateLoggerWithLevel(serviceName string, atomicLevel zap.AtomicLevel, opts ...zap.Option) *zap.Logger {
	logger, err := createLoggerWithLevel(serviceName, atomicLevel, opts...)
	if err != nil {
		panic(fmt.Errorf("unable to create logger (in production: %t): %s", isProductionEnvironment(), err))
	}

	return logger
}

func createLoggerWithLevel(serviceName string, atomicLevel zap.AtomicLevel, opts ...zap.Option) (*zap.Logger, error) {
	config := BasicLoggingConfig(serviceName, atomicLevel, opts...)

	zlog, err := config.Build(opts...)
	if err != nil {
		return nil, err
	}

	if AutoStartServer {
		go func() {
			zlog.Info("starting atomic level switcher, port :1065")
			if err := http.ListenAndServe(":1065", atomicLevel); err != nil {
				zlog.Info("failed listening on :1065 to switch log level:", zap.Error(err))
			}
		}()
	}

	return zlog, nil
}

func BasicLoggingConfig(serviceName string, atomicLevel zap.AtomicLevel, opts ...zap.Option) *zap.Config {
	var config zap.Config

	if isProductionEnvironment() || os.Getenv("ZAP_PRETTY") != "" {
		config = zapdriver.NewProductionConfig()
		opts = append(opts, zapdriver.WrapCore(
			zapdriver.ReportAllErrors(true),
			zapdriver.ServiceName(serviceName),
		))
	} else {
		config = zap.NewDevelopmentConfig()
	}

	if os.Getenv("ZAP_PRETTY") != "" {
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stdout"}
	}

	config.Level = atomicLevel
	return &config
}

func levelFromEnvironment() zap.AtomicLevel {
	zapPrettyValue := os.Getenv("ZAP_PRETTY")
	if zapPrettyValue != "" {
		return zap.NewAtomicLevelAt(zapLevelFromString(zapPrettyValue))
	}

	if isProductionEnvironment() {
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return zap.NewAtomicLevelAt(zap.DebugLevel)
}

func isProductionEnvironment() bool {
	_, err := os.Stat("/.dockerenv")
	if !os.IsNotExist(err) {
		return true
	}

	goEnv := os.Getenv("GO_ENV")
	if goEnv != "" && ((goEnv == "production") || (goEnv == "staging")) {
		return true
	}
	return false
}

func zapLevelFromString(input string) zapcore.Level {
	switch strings.ToLower(input) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warning", "warn":
		return zap.WarnLevel
	case "error", "err":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.DebugLevel
	}
}
