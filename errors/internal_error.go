package errors

var OriginList = struct {
	Core    string
	Service string
	Module  string
}{Core: "core", Service: "service", Module: "module"}

type origin struct {
	OriginType string
	OriginName string
}

type InternalError struct {
	internal error
	Level    string
	Origin   origin
}

func NewInternalError(err error, level, originType, originName string) *InternalError {
	return &InternalError{internal: err, Level: level, Origin: origin{OriginType: originType, OriginName: originName}}
}

func (err *InternalError) Error() string {
	return err.internal.Error()
}
