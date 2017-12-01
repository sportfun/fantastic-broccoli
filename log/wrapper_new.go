package log

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sportfun/gakisitor/config"
)

type writeSyncedImpl struct {
	*string
}

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

func NewProduction(properties ...config.LogDefinition) Logger {
	if len(properties) == 0 {
		logger := NewDevelopment()
		logger.Error(&argumentBinderImpl{format: "No log configuration for production, switch to development logger"})
		return logger
	}

	var cores []zapcore.Core
	for _, definition := range properties {
		cores = append(cores, newProductionZapCore(definition))
	}

	return &loggerImpl{instance: zap.New(zapcore.NewTee(cores...)), forProduction: true}
}

func NewDevelopment(properties ...config.LogDefinition) Logger {
	if len(properties) > 0 && properties[0].Raw != nil {
		if buffer, ok := properties[0].Raw.(*string); ok {
			return &loggerImpl{instance: zap.New(newTestZapCore(buffer)), forProduction: false}
		}
	}

	return &loggerImpl{instance: zap.New(newDevelopmentZapCore()), forProduction: false}
}

func newProductionZapCore(definition config.LogDefinition) zapcore.Core {
	var encoder zapcore.Encoder
	switch definition.Encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	default:
		encoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
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

func newDevelopmentZapCore() zapcore.Core {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
}

func newTestZapCore(buffer *string) zapcore.Core {
	enab := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})
	enc := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("")
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("")
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
	})
	ws := zapcore.Lock(&writeSyncedImpl{string: buffer})
	return zapcore.NewCore(enc, ws, enab)
}

func (ws *writeSyncedImpl) Write(p []byte) (n int, err error) {
	*ws.string = strings.TrimSpace(string(p))
	return len(p), nil
}

func (ws *writeSyncedImpl) Sync() error {
	return nil
}
