package log

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
