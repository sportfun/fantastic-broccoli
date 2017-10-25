package network

type webPacket struct {
	linkId string
	body   interface{}
}

func (s *Service) On(method string, f interface{}) bool {
	if err := s.client.On(method, f); err != nil {
		s.errorHandler(SocketOn, err)
		return false
	}
	return true
}

func (s *Service) Emit(method string, body interface{}) bool {
	if err := s.client.Emit(method, webPacket{s.linkId, body}); err != nil {
		s.errorHandler(SocketEmit, err)
		return false
	}
	return true
}
