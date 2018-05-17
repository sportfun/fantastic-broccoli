package main

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/plugin/plugin_test"
)

func TestPlugin(t *testing.T) {
	plugin_test.PluginValidityChecker(t, &Plugin, plugin_test.PluginTestDesc{
		ConfigJSON: `{"ManyItems": {"ThisItem": 0}}`,
		ValueChecker: gomega.WithTransform(func(v interface{}) time.Time {
			return v.(struct {
				A int       `json:"a"`
				B time.Time `json:"b"`
			}).B
		}, gomega.BeTemporally("~", time.Now(), 2*time.Second)),
	})
}
