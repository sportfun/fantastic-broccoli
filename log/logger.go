package log

import "github.com/xunleii/fantastic-broccoli/properties"

type argumentBinder interface {
	Bind(a ...interface{}) argumentBinder
	More(index string, data interface{}) argumentBinder
	getMessage() string
	getMoreInfo() map[string]interface{}
}

type Logger interface {
	Debug(argumentBinder)
	Info(argumentBinder)
	Warn(argumentBinder)
	Error(argumentBinder)
}

type loggerFactory func(*properties.Properties) Logger

var NewLogger = struct {
	Prod loggerFactory
	Dev  loggerFactory
}{
	Prod: newProdLogger,
	Dev:  newDevLogger,
}