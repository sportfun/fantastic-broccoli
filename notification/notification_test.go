package notification

import (
	"testing"
	"fantastic-broccoli/utils"
)

type casterTest struct{}

func (c *casterTest) cast(o Origin, obj Object) (casted Object, err error) {
	return "Casted", nil
}

func TestNewNotification(t *testing.T) {
	n := Notification{"Origin", "Destination", struct{}{}}
	p := NewNotification("Origin", "Destination", struct{}{})

	utils.AssertEquals(t, n, *p)
}

func TestNotification(t *testing.T) {
	n := Notification{"Origin", "Destination", struct{}{}}
	c := Caster(new(casterTest))

	utils.AssertEquals(t, Origin("Origin"), n.From())
	utils.AssertEquals(t, Destination("Destination"), n.To())
	utils.AssertEquals(t, struct{}{}, n.Content())

	a, _ := c.cast("", nil)
	b, _ := n.Cast(c)
	utils.AssertEquals(t, a, b)
}
