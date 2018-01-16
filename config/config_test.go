package config

import (
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/utils"
	"io/ioutil"
	"testing"
	"time"
)

const (
	resourceConfigPath        = "../.resources/gakisitor.conf"
	invalidResourceConfigPath = "../.resources/invalid.conf"
	iofsResourceConfigPath    = "../.resources/iofs.conf"
)

func TestGAkisitorConfig_Load(t *testing.T) {
	RegisterTestingT(t)

	config := GAkisitorConfig{}
	testCases := []struct {
		File    string
		Error   string
		Matcher OmegaMatcher
	}{
		{File: "", Error: "impossible to read the configuration file: open :", Matcher: BeFalse()},
		{File: "./none", Error: "impossible to read the configuration file: open ./none:", Matcher: BeFalse()},
		{File: invalidResourceConfigPath, Error: "impossible to unmarshal the configuration file: invalid character 'u' looking for beginning of value", Matcher: BeFalse()},
		{File: resourceConfigPath, Error: "", Matcher: BeTrue()},
	}

	for _, tc := range testCases {
		*config.FilePtr() = tc.File
		Expect(config.file).Should(Equal(tc.File))

		if tc.Error == "" {
			Expect(config.Load()).Should(Succeed())
		} else {
			Expect(config.Load()).Should(MatchError(MatchRegexp(tc.Error)))
		}
		Expect(config.IsLoaded()).Should(tc.Matcher)
	}
}

func TestGAkisitorConfig_WaitReconfiguration(t *testing.T) {
	RegisterTestingT(t)

	config := GAkisitorConfig{}
	testCases := []struct {
		File     string
		Duration time.Duration
		Timeout  time.Duration
		Writer   func(file string)
	}{
		{File: iofsResourceConfigPath, Duration: 250 * time.Millisecond, Timeout: 300 * time.Millisecond, Writer: func(file string) { time.Sleep(50 * time.Millisecond); ioutil.WriteFile(file, []byte("{}"), 0644) }},
		{File: iofsResourceConfigPath, Duration: 250 * time.Millisecond, Timeout: 300 * time.Millisecond, Writer: func(string) {}},
	}

	for _, tc := range testCases {
		config.file = tc.File

		go tc.Writer(tc.File)
		utils.ReleaseIfTimeout(t, tc.Timeout, func(testing.TB) {
			Expect(config.WaitReconfiguration(tc.Duration)).Should(Succeed())
		})
	}
}
