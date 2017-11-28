package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/xunleii/fantastic-broccoli/properties"
)

var errorLevelMapping = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

var fileNameMapping = map[string]*os.File{
	"stdout": os.Stdout,
	"stderr": os.Stderr,
}

func newZapCore(definition properties.LogDefinition) zapcore.Core {
	var encoder zapcore.Encoder
	switch definition.Encoding {
	case "json":
		encoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	default:
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	file, ok := fileNameMapping[strings.ToLower(definition.File)]
	if !ok {
		var err error
		file, err = os.Open(definition.File)
		if err != nil {
			file = os.Stdout
		}
	}
	writer := zapcore.Lock(file)

	level, ok := errorLevelMapping[strings.ToLower(definition.Level)]
	if !ok {
		level = zapcore.WarnLevel
	}
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	return zapcore.NewCore(encoder, writer, levelEnabler)
}

func newProdLogger(properties *properties.Properties) Logger {
	if len(properties.Log) == 0 {
		logger := newDevLogger(properties)
		logger.Error(&argumentBinderImpl{format: "No log configuration for production, switch to development logger"})
		return logger
	}

	var cores []zapcore.Core
	for _, definition := range properties.Log {
		cores = append(cores, newZapCore(definition))
	}

	logger := zap.New(zapcore.NewTee(cores...))

	return &loggerImpl{instance: logger}
}

func newDevLogger(*properties.Properties) Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
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

	return &loggerImpl{instance: logger}
}
