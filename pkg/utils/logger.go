package utils

import (
	"os"
	"strings"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger with clean console output
func NewLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	
	// Detect if we should use colors (only in terminals that support it)
	useColors := shouldUseColors()
	
	// Clean console output: level + message + fields
	if useColors {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	
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

// shouldUseColors determines if the terminal supports ANSI colors
func shouldUseColors() bool {
	// Check if running in Windows Terminal, PowerShell, or other modern terminals
	term := strings.ToLower(os.Getenv("TERM"))
	wtSession := os.Getenv("WT_SESSION")
	
	// Windows Terminal or terminals with TERM set (Git Bash, WSL, etc.)
	if wtSession != "" || term != "" {
		return true
	}
	
	// Disable colors for Windows CMD
	return false
}

// NewProductionLogger creates a production logger with JSON output
func NewProductionLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	
	logger, _ := config.Build()
	return logger
}
