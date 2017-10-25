package notification

type Builder struct {
	Notification
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) From(o string) *Builder {
	b.from = o
	return b
}

func (b *Builder) To(d string) *Builder {
	b.to = d
	return b
}

func (b *Builder) With(o interface{}) *Builder {
	b.content = o
	return b
}

func (b *Builder) Build() *Notification {
	return &b.Notification
}
