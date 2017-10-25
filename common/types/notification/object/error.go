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

func (errorObject *ErrorObject) Why(reason error) *ErrorObject {
	errorObject.reason = reason
	return errorObject
}

func (errorObject *ErrorObject) Origin() string {
	return errorObject.origin
}

func (errorObject *ErrorObject) Reason() error {
	return errorObject.reason
}
