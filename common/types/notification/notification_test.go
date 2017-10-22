package notification

import (
	"testing"
	"fantastic-broccoli/utils"
)

func TestNewNotification(t *testing.T) {
	n := Notification{"Origin", "Destination", struct{}{}}
	p := NewNotification("Origin", "Destination", struct{}{})

	utils.AssertEquals(t, n, *p)
}

func TestNotification(t *testing.T) {
	n := Notification{"Origin", "Destination", struct{}{}}

	utils.AssertEquals(t, Origin("Origin"), n.From())
	utils.AssertEquals(t, Destination("Destination"), n.To())
	utils.AssertEquals(t, struct{}{}, n.Content())
}
