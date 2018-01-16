package kernel

type origin struct {
	entity string
	name   string
}

type internalError struct {
	internal error
	level    string
	origin   origin
}

func NewInternalError(err error, level, originEntity, originName string) *internalError {
	return &internalError{internal: err, level: level, origin: origin{entity: originEntity, name: originName}}
}

func (err *internalError) Error() string {
	return err.internal.Error()
}
