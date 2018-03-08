// Package plugin provides types used to create plugins.
//
// A plugin is a "code" used to extend functionality of the
// gakisitor, like new user input. In fact, because the gakisitor
// must be working with different user interface (like training
// bike or simple button), we need to easily extend its functionality,
// working for several user input through GPIO (button, distance
// sensor, ...).
package plugin

import (
	"context"

	"github.com/sportfun/gakisitor/profile"
)

// Instruction type (only accessed in this package to prevent custom instructions).
type instruction byte

// Chan struct containing all channels used by the plugins.
type Chan struct {
	Data        chan<- interface{} // Used to send data (only JSON serializable data)
	Status      chan<- State       // Used to send status
	Instruction <-chan instruction // Used to read instruction from the gakisitor
}

// Instruction list (instructions are immutable to prevent custom instructions)
const (
	StatusPluginInstruction instruction = 0x01 // Send a the current status
	StartSessionInstruction instruction = 0x10 // Start a game session (you MUST retrieve user input during this session)
	StopSessionInstruction  instruction = 0x1F // Stop the game session (you MUST stop your retrieving user input)
)

// Plugin structure used to extend the Gakisitor functionality
type Plugin struct {
	// The plugin name. It will be used by the server/game
	// engine to know which plugin the data comes from.
	Name string

	// Start the plugin instance with the plugin profile and channels. You
	// MUST check the profile before starting the process.
	//
	// For more information about plugin, see the package description.
	// For more information about plugin channels, see the Chan structure above.
	Instance func(ctx context.Context, profile profile.Plugin, channels Chan) error
}
