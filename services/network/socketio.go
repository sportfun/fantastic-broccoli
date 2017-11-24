package network

type webPacket struct {
	LinkId string      `json:"link_id" mapstructure:"link_id"`
	Body   interface{} `json:"body" mapstructure:"body"`
}

func (service *Service) on(method string, f interface{}) bool {
	return service.checkIf(nil, service.client.On(method, f), IsListening)
}

func (service *Service) emit(method string, body interface{}) bool {
	return service.checkIf(nil, service.client.Emit(method, webPacket{service.linkId, body}), IsEmitted)
}
