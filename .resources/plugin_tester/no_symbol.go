package main

import (
	"context"

	"github.com/sportfun/gakisitor/plugin"
	"github.com/sportfun/gakisitor/profile"
)

var InvalidPlugin = plugin.Plugin{
	Name: "Plugin Example",
	Instance: func(ctx context.Context, profile profile.Plugin, channels plugin.Chan) error {
		return nil
	},
}
