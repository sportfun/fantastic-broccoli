package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sportfun/gakisitor/config"
)

func NewProduction(properties ...config.LogDefinition) Logger {
	if len(properties) == 0 {
		logger := NewDevelopment()
		logger.Error(&argumentBinderImpl{format: "No log configuration for production, switch to development logger"})
		return logger
	}

	var cores []zapcore.Core
	for _, definition := range properties {
		cores = append(cores, newProdCore(definition))
	}

	return &loggerImpl{instance: zap.New(zapcore.NewTee(cores...)), forProduction: true}
}

func NewDevelopment(...config.LogDefinition) Logger {
	return &loggerImpl{instance: zap.New(newDevCore()), forProduction: false}
}

func NewTest(buffer *string) Logger {
	return &loggerImpl{instance: zap.New(newTestCore(buffer)), forProduction: false}
}
