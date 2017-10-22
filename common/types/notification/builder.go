package notification

type Builder struct {
	Notification
}

func (b *Builder) From(o Origin) *Builder  {
	b.from = o
	return b
}

func (b *Builder) To(d Destination) *Builder  {
	b.to = d
	return b
}

func (b *Builder) With(o Object) *Builder  {
	b.content = o
	return b
}

func (b *Builder) Build() Notification  {
	return b.Notification
}