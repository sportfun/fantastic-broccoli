package network

type websocket struct {
	LinkId string      `json:"link_id" mapstructure:"link_id"`
	Body   interface{} `json:"body" mapstructure:"body"`
}

const nilClient = "socket.io client not initialised"

func (service *Network) on(method string, f interface{}) bool {
	if service.client == nil {
		service.logger.Errorf(nilClient)
		return false
	}

	return service.checkIf(isListening, nil, service.client.On(method, f))
}

func (service *Network) emit(method string, body interface{}) bool {
	if service.client == nil {
		service.logger.Errorf(nilClient)
		return false
	}

	return service.checkIf(isEmitted, nil, service.client.Emit(method, websocket{service.linkId, body}))
}
