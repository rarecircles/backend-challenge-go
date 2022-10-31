package logging

import (
	"context"
	"time"

	"github.com/teris-io/shortid"
	"go.uber.org/zap"
)

type loggerKeyType int

var shortIDGenerator *shortid.Shortid

const loggerKey loggerKeyType = iota

func init() {
	// A new generator using the default alphabet set
	shortIDGenerator = shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().UnixNano()))
}

// WithLogger is used to create a new context with a logging added to it
// So it can be later retrieved using `Logger`.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Logger is used to retrieve the logging from the context. If no logging
// is present in the context, the `fallbackLogger` received in parameter
// is returned instead.
func Logger(ctx context.Context, fallbackLogger *zap.Logger) *zap.Logger {
	if ctx == nil {
		return fallbackLogger
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}

	return fallbackLogger
}
