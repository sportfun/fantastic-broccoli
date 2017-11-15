package network

type webPacket struct {
	LinkId string      `json:"link_id" mapstructure:"link_id"`
	Body   interface{} `json:"body" mapstructure:"body"`
}

func (s *Service) on(method string, f interface{}) bool {
	if err := s.client.On(method, f); err != nil {
		s.errorHandler(SocketOn, err)
		return false
	}
	return true
}

func (s *Service) emit(method string, body interface{}) bool {
	if err := s.client.Emit(method, webPacket{s.linkId, body}); err != nil {
		s.errorHandler(SocketEmit, err)
		return false
	}
	return true
}
