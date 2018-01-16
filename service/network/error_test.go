package network

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/service"
	"testing"
)

func TestNetwork_ModuleError(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	network := &Network{logger: log.NewTest(&buffer), notifications: service.NewNotificationQueue()}

	testCases := []struct {
		Fnc    netError
		Error  error
		Reason string
	}{
		{Fnc: isEmitted, Error: fmt.Errorf("err#1"), Reason: "ERROR	failed to emit message: err#1"},
		{Fnc: isListening, Error: fmt.Errorf("err#2"), Reason: "ERROR	failed to create channel handler: err#2"},
	}

	for _, tc := range testCases {
		buffer = ""
		tc.Fnc(network, nil, tc.Error)
		Expect(buffer).Should(Equal(tc.Reason))
	}
}

func TestNetwork_CheckIf(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	network := &Network{logger: log.NewTest(&buffer), notifications: service.NewNotificationQueue()}

	testCases := []struct {
		Fnc    netError
		Error  error
		Reason string
	}{
		{Fnc: nil, Error: nil, Reason: ""},
		{Fnc: isEmitted, Error: fmt.Errorf("err#1"), Reason: "ERROR	failed to emit message: err#1"},
		{Fnc: isListening, Error: fmt.Errorf("err#2"), Reason: "ERROR	failed to create channel handler: err#2"},
	}

	for _, tc := range testCases {
		buffer = ""
		network.checkIf(tc.Fnc, nil, tc.Error)
		Expect(buffer).Should(Equal(tc.Reason))
	}
}
