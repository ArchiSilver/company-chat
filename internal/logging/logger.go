package logging

import (
	"go.uber.org/zap"
)

var L *zap.SugaredLogger

// Init initializes the global logger. Call with isProd=true in production.
func Init(isProd bool) {
	var logger *zap.Logger
	var err error
	if isProd {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	L = logger.Sugar()
}

// Sync flushes any buffered log entries
func Sync() error {
	if L == nil {
		return nil
	}
	return L.Sync()
}
