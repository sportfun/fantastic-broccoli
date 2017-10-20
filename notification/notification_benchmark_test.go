package notification

import "testing"

func BenchmarkBuilderAll(b *testing.B) {
	bd := Builder{}
	ori := Origin("Origin")
	des := Destination("Destination")
	o := struct {
		a string
		b int
		c bool
	}{"a", 0, false}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n := bd.From(ori).To(des).With(o).Build()
		n.Content()
	}
}

func BenchmarkBuilderBuild(b *testing.B) {
	bd := Builder{}
	bd.From("Origin").To("Destination").With(struct {
		a string
		b int
		c bool
	}{"a", 0, false})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n := bd.Build()
		n.Content()
	}
}

func BenchmarkBuilderNotificationOnly(b *testing.B) {
	ori := Origin("Origin")
	des := Destination("Destination")
	o := struct {
		a string
		b int
		c bool
	}{"a", 0, false}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n := Notification{ori, des, o}
		n.Content()
	}
}
