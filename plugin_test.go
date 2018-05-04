package main

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/profile"
)

func TestPlugin_Task(t *testing.T) {
	RegisterTestingT(t)

	Profile = &profile.Profile{
		Plugins: []profile.Plugin{
			{
				Name: "Example plugin",
				Path: "./.resources/plugin_example.so",
				Config: map[string]interface{}{
					"ManyItems": map[string]interface{}{
						"ThisItem": 0,
					},
				},
			},
		},
	}

}
