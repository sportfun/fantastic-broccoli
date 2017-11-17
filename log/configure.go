package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	"fantastic-broccoli/properties"
)

//Todo: Need implementation
func Configure(props *properties.Properties) *zap.Logger {
	// Add property parameter if needed in the future
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	logger := zap.New(core)
	defer logger.Sync()
	return logger}
