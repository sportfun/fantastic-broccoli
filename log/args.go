package log

type argumentBinder interface {
	Bind(a ...interface{}) argumentBinder
	More(index string, data interface{}) argumentBinder
}

func NewArgumentBinder(format string) argumentBinder {
	return &argsBinderImpl{}
}

type argsBinderImpl struct {
}

func (binder *argsBinderImpl) Bind(a ...interface{}) argumentBinder {
	return binder
}

func (binder *argsBinderImpl) More(index string, data interface{}) argumentBinder {
	return binder
}
