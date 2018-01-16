package log

type argumentBinder interface {
	Bind(a ...interface{}) argumentBinder
	More(index string, data interface{}) argumentBinder
	getMessage() string
	getMoreInfo() map[string]interface{}
}

type Logger interface {
	Debug(argumentBinder)
	Debugf(format string, a ...interface{})
	Info(argumentBinder)
	Infof(format string, a ...interface{})
	Warn(argumentBinder)
	Warnf(format string, a ...interface{})
	Error(argumentBinder)
	Errorf(format string, a ...interface{})
}
