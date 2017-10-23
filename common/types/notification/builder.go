package notification

import "fantastic-broccoli/common/types"

type Builder struct {
	Notification
}

func (b *Builder) From(o types.Name) *Builder {
	b.from = o
	return b
}

func (b *Builder) To(d types.Name) *Builder {
	b.to = d
	return b
}

func (b *Builder) With(o Object) *Builder {
	b.content = o
	return b
}

func (b *Builder) Build() *Notification {
	return &b.Notification
}
