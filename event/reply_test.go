package event

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestSyncReplyHandler(t *testing.T) {
	RegisterTestingT(t)

	dch := make(chan interface{})
	ech := make(chan error)

	cases := []struct {
		err error

		fnc          func()
		dataMatcher  OmegaMatcher
		errorMatcher OmegaMatcher
	}{
		{ErrChannelNotFound, func() {}, BeNil(), MatchError(ErrChannelNotFound)},
		{nil, func() { dch <- "data" }, Equal("data"), BeNil()},
		{nil, func() { ech <- ErrChannelNotFound }, BeNil(), MatchError(ErrChannelNotFound)},
		{nil, func() {}, BeNil(), MatchError(ErrReplyTimeout)},
	}

	for _, test := range cases {
		go test.fnc()

		SyncReplyHandler(func(data interface{}, err error) {
			Expect(data).Should(test.dataMatcher)
			Expect(err).Should(test.errorMatcher)
		}).consume(dch, ech, test.err, time.Millisecond)
	}
}
