package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger with clean console output
func NewLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	
	// Clean console output: colored level + message + fields
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("")  // No timestamp
	config.EncoderConfig.EncodeCaller = nil                             // No caller info
	config.DisableCaller = true                                         // Disable caller completely
	config.DisableStacktrace = true                                     // Disable stacktraces
	
	// Use console encoder for better readability
	config.Encoding = "console"
	config.EncoderConfig.ConsoleSeparator = " "
	
	logger, _ := config.Build()
	return logger
}

// NewProductionLogger creates a production logger with JSON output
func NewProductionLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	
	logger, _ := config.Build()
	return logger
}
