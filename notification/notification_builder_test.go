package notification

import (
	"testing"
)

func TestBuilderFrom(t *testing.T) {
	b := Builder{}
	b.From("Origin")

	if b._from != "Origin" {
		t.Errorf("Expected 'Origin' (%s)", b._from)
		t.Fail()
	}
}

func TestBuilderTo(t *testing.T) {
	b := Builder{}
	b.To("Destination")

	if b._to != "Destination" {
		t.Errorf("Expected 'Destination' (%s)", b._to)
		t.Fail()
	}
}

func TestBuilderWith(t *testing.T) {
	b := Builder{}
	o := struct {
		a string
		b int
		c bool
	}{"a", 0, false}
	b.With(o)

	if b._content != o {
		t.Errorf("Expected '%v' (%v)", o, b._content)
		t.Fail()
	}
}

func TestBuilderBuild(t *testing.T) {
	b := Builder{}
	o := struct {
		a string
		b int
		c bool
	}{"a", 0, false}
	d1 := Notification{"Origin", "Destination", o}
	d2 := Notification{"", "Destination", o}
	d3 := Notification{"Destination", "Origin", o}
	d4 := Notification{"Destination", "Origin", struct{}{}}

	n := b.From("Origin").To("Destination").With(o).Build()
	if n != d1 {
		t.Errorf("Expected '%v' (%v)", d1, n)
	}

	n = b.From("").Build()
	if n != d2 {
		t.Errorf("Expected '%v' (%v)", d2, n)
	}

	n = b.From("Destination").To("Origin").Build()
	if n != d3 {
		t.Errorf("Expected '%v' (%v)", d3, n)
	}

	n = b.With(struct{}{}).Build()
	if n != d4 {
		t.Errorf("Expected '%v' (%v)", d4, n)
	}
}
