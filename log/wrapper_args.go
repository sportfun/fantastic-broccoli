package log

import "fmt"

type argumentBinderImpl struct {
	format string
	binds  []interface{}
	more   map[string]interface{}
}

func NewArgumentBinder(format string) argumentBinder {
	return &argumentBinderImpl{format: format, binds: nil, more: make(map[string]interface{})}
}

func (binder *argumentBinderImpl) Bind(a ...interface{}) argumentBinder {
	binder.binds = a
	return binder
}

func (binder *argumentBinderImpl) More(index string, data interface{}) argumentBinder {
	binder.more[index] = data
	return binder
}

func (binder *argumentBinderImpl) getMessage() string {
	return fmt.Sprintf(binder.format, binder.binds...)
}

func (binder *argumentBinderImpl) getMoreInfo() map[string]interface{} {
	return binder.more
}
