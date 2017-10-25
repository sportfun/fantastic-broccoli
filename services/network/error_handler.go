package network

type errorType int

const (
	SocketOn   errorType = iota
	SocketEmit
)

func (s *Service) errorHandler(t errorType, e error, p ...interface{}) {
	if e == nil {
		return
	}
}
