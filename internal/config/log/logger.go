package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	defaultLog "log"
	"os"
	"strings"
	"sync"
)

var (
	once   = new(sync.Once)
	Logger *zap.SugaredLogger
)

func SetupLogger() {
	once.Do(func() {
		log, err := newLogger()
		if err != nil {
			defaultLog.Fatalf("Can't initialize logger: %v", err)
		}
		Logger = log.Sugar()
		Logger.Debug("Created new logger")
	})
}

func newLogger() (*zap.Logger, error) {
	env := os.Getenv("ENV")

	var config zap.Config
	if env == "development" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	if lvl, exists := os.LookupEnv("LOG_LEVEL"); exists {
		lvl = strings.ToLower(lvl)
		switch lvl {
		case "debug":
			config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		case "info":
			config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		case "warn":
			config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		case "error":
			config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		case "panic":
			config.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
		case "fatal":
			config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
		}
	}

	return config.Build()
}
