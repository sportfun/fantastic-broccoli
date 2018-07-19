package main

import (
	"context"

	"github.com/sportfun/gakisitor/plugin"
	"github.com/sportfun/gakisitor/profile"
)

var Plugin = plugin.Plugin{
	Name: "already_exists",
	Instance: func(ctx context.Context, profile profile.Plugin, channels plugin.Chan) error {
		return nil
	},
}
