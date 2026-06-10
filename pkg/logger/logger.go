package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	Audit *zap.Logger
}

func NewLogger(devMode bool, auditLogStorage string) (*Logger, error) {
	var cfg zap.Config
	if devMode {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout"}

	logger, err := cfg.Build()
	if err != nil {
		return &Logger{}, fmt.Errorf("mainLogger.Build: %w", err)
	}

	auditCfg := cfg
	auditCfg.OutputPaths = []string{auditLogStorage}
	auditLogger, err := auditCfg.Build()
	if err != nil {
		return &Logger{}, fmt.Errorf("auditLogger.Build: %w", err)

	}

	return &Logger{
		logger,
		auditLogger,
	}, nil
}
