package module

import (
	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/constant"
)

type Error interface {
	error
	ErrorLevel() types.ErrorLevel
}

type moduleError struct {
	err error
	level types.ErrorLevel
}


func (e *moduleError) Error() string {
	return e.err.Error()
}

func (e *moduleError) ErrorLevel() types.ErrorLevel {
	return e.level
}


func Warned(reason error) Error {
	return &moduleError{err:reason, level:constant.ErrorLevels.Warning}
}

func Failed(reason error) Error {
	return &moduleError{err:reason, level:constant.ErrorLevels.Error}
}

func Crashed(reason error) Error {
	return &moduleError{err:reason, level:constant.ErrorLevels.Critical}
}

func Panicked(reason error) Error {
	return &moduleError{err:reason, level:constant.ErrorLevels.Fatal}
}


