package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"korun.io/shared/config"
)

func NewLogger(cfg config.LoggerConfig) (*zap.Logger, error) {
	var zapConfig zap.Config

	switch cfg.Format {
	case "json":
		zapConfig = zap.NewProductionConfig()
	case "console":
		zapConfig = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("unknown logger format: %s", cfg.Format)
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return logger, nil
}
