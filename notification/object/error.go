package object

type ErrorObject struct {
	Origin string `json:"origin" mapstructure:"origin"`
	Reason string `json:"reason" mapstructure:"reason"`
}

func NewErrorObject(origin string, reason ...error) *ErrorObject {
	if len(reason) > 0 && reason[0] != nil {
		return &ErrorObject{Origin: origin, Reason: reason[0].Error()}
	}
	return &ErrorObject{Origin: origin, Reason: ""}
}

func (errorObject *ErrorObject) Why(reason error) *ErrorObject {
	errorObject.Reason = reason.Error()
	return errorObject
}
