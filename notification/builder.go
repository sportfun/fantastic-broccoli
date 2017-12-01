package notification

type builder struct {
	Notification
}

func NewBuilder() *builder {
	return &builder{}
}

func (builder *builder) From(o string) *builder {
	builder.from = o
	return builder
}

func (builder *builder) To(d string) *builder {
	builder.to = d
	return builder
}

func (builder *builder) With(o interface{}) *builder {
	builder.content = o
	return builder
}

func (builder *builder) Build() *Notification {
	n := builder.Notification
	return &n
}
