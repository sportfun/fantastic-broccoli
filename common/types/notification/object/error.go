package object

type ErrorObject struct {
	origin string
	reason error
}

func NewErrorObject(origin string, reason ...error) *ErrorObject {
	if len(reason) > 0 {
		return &ErrorObject{origin: origin, reason: reason[0]}
	}
	return &ErrorObject{origin: origin, reason: nil}
}

func (e *ErrorObject) Why(reason error) *ErrorObject {
	e.reason = reason
	return e
}

func (e *ErrorObject) Origin() string {
	return e.origin
}

func (e *ErrorObject) Reason() error {
	return e.reason
}
