package log

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerImpl struct {
	instance      *zap.Logger
	forProduction bool
}

type infoWrapper func(string, interface{}) zapcore.Field

var typeWrapping = map[reflect.Kind]infoWrapper{
	reflect.String:     func(i string, v interface{}) zapcore.Field { return zap.String(i, v.(string)) },
	reflect.Bool:       func(i string, v interface{}) zapcore.Field { return zap.Bool(i, v.(bool)) },
	reflect.Int:        func(i string, v interface{}) zapcore.Field { return zap.Int(i, v.(int)) },
	reflect.Int8:       func(i string, v interface{}) zapcore.Field { return zap.Int8(i, v.(int8)) },
	reflect.Int16:      func(i string, v interface{}) zapcore.Field { return zap.Int16(i, v.(int16)) },
	reflect.Int32:      func(i string, v interface{}) zapcore.Field { return zap.Int32(i, v.(int32)) },
	reflect.Int64:      func(i string, v interface{}) zapcore.Field { return zap.Int64(i, v.(int64)) },
	reflect.Uint:       func(i string, v interface{}) zapcore.Field { return zap.Uint(i, v.(uint)) },
	reflect.Uint8:      func(i string, v interface{}) zapcore.Field { return zap.Uint8(i, v.(uint8)) },
	reflect.Uint16:     func(i string, v interface{}) zapcore.Field { return zap.Uint16(i, v.(uint16)) },
	reflect.Uint32:     func(i string, v interface{}) zapcore.Field { return zap.Uint32(i, v.(uint32)) },
	reflect.Uint64:     func(i string, v interface{}) zapcore.Field { return zap.Uint64(i, v.(uint64)) },
	reflect.Float32:    func(i string, v interface{}) zapcore.Field { return zap.Float32(i, v.(float32)) },
	reflect.Float64:    func(i string, v interface{}) zapcore.Field { return zap.Float64(i, v.(float64)) },
	reflect.Complex64:  func(i string, v interface{}) zapcore.Field { return zap.Complex64(i, v.(complex64)) },
	reflect.Complex128: func(i string, v interface{}) zapcore.Field { return zap.Complex128(i, v.(complex128)) },
}

func toFields(binder argumentBinder) []zapcore.Field {
	var fields []zapcore.Field

	for i, v := range binder.getMoreInfo() {
		if v == nil {
			continue
		}

		if fnc, ok := typeWrapping[reflect.TypeOf(v).Kind()]; ok {
			fields = append(fields, fnc(i, v))
		} else {
			fields = append(fields, zap.Reflect(i, v))
		}
	}

	return fields
}

func (logger *loggerImpl) Debug(a argumentBinder) {
	logger.instance.Debug(a.getMessage(), toFields(a)...)
}

func (logger *loggerImpl) Debugf(f string, a ...interface{}) {
	logger.instance.Debug(fmt.Sprintf(f, a...))
}

func (logger *loggerImpl) Info(a argumentBinder) { logger.instance.Info(a.getMessage(), toFields(a)...) }

func (logger *loggerImpl) Infof(f string, a ...interface{}) {
	logger.instance.Info(fmt.Sprintf(f, a...))
}

func (logger *loggerImpl) Warn(a argumentBinder) { logger.instance.Warn(a.getMessage(), toFields(a)...) }

func (logger *loggerImpl) Warnf(f string, a ...interface{}) {
	logger.instance.Warn(fmt.Sprintf(f, a...))
}

func (logger *loggerImpl) Error(a argumentBinder) {
	logger.instance.Error(a.getMessage(), toFields(a)...)
}

func (logger *loggerImpl) Errorf(f string, a ...interface{}) {
	logger.instance.Error(fmt.Sprintf(f, a...))
}
