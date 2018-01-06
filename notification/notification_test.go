package notification

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestNotification(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		From    string
		To      string
		Content interface{}
	}{
		{From: "from", To: "to", Content: 5},
		{From: "", To: "", Content: nil},
		{From: "from", To: "to", Content: struct{ A int }{A: 0}},
	}

	for _, tc := range testCases {
		notification := NewNotification(tc.From, tc.To, tc.Content)

		Expect(notification.From()).Should(Equal(tc.From))
		Expect(notification.To()).Should(Equal(tc.To))
		if tc.Content == nil {
			Expect(notification.Content()).Should(BeNil())
		} else {
			Expect(notification.Content()).Should(Equal(tc.Content))
		}
	}
}

func TestNotificationBuilder(t *testing.T) {
	RegisterTestingT(t)

	builder := NewBuilder()
	testCases := []struct {
		From    string
		To      string
		Content interface{}
	}{
		{From: "from", To: "to", Content: 5},
		{From: "", To: "", Content: nil},
		{From: "from", To: "to", Content: struct{ A int }{A: 0}},
	}

	for _, tc := range testCases {
		builder.From(tc.From).To(tc.To).With(tc.Content)
		Expect(builder.Build()).Should(Equal(NewNotification(tc.From, tc.To, tc.Content)))
	}
}
