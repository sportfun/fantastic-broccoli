package log

type Logger interface {
	Debug(argumentBinder)
	Info(argumentBinder)
	Warn(argumentBinder)
	Error(argumentBinder)
	Fatal(argumentBinder)
}