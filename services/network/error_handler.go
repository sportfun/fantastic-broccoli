package network

import "go.uber.org/zap"

type errorType int

const (
	SocketOn   errorType = iota
	SocketEmit
)

func (s *Service) errorHandler(t errorType, err error, a ...interface{}) {
	if err == nil {
		return
	}

	switch t {
	case SocketEmit:
		s.logger.Error("Failed to emit message", zap.String("reason", err.Error()))
	case SocketOn:
		s.logger.Error("Failed to create channel handler", zap.String("reason", err.Error()))
	}
}
