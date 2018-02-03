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
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/profile"
)

// State represent the current plugin state. State are immutable
// (to prevent predefined state edition).
type State struct {
	code byte        // State code
	desc string      // State name
	raw  interface{} // Additional and optional content (error or string)
}

// Predefined states.
var (
	// State code formats
	// +---------+-----------------+--------------------------------------------------------------+
	// |  Range  |   Description   |                         Behaviours                           |
	// +---------+-----------------+--------------------------------------------------------------+
	// | 1x ~ 1F | OK states       | Nothing                                                      |
	// | 2x ~ DF | Custom states   | Nothing                                                      |
	// | Ex ~ EF | Error states    | The scheduler notify the server                              |
	// | Fx ~ FF | Critical states | The scheduler try to restart the plugin && notify the server |
	// +---------+-----------------+--------------------------------------------------------------+
	NilState State

	RunningState   = State{0x11, "currently running", nil}
	InSessionState = State{0x12, "currently in session", nil}
	PausedState    = State{0x13, "currently paused", nil}
	StoppedState   = State{0x14, "currently stopped", nil}

	GPIODisconnectionState  = State{0xE1, "GPIO has been disconnected", nil}
	GPIOFailureState        = State{0xE2, "GPIO reading has failed", nil}
	AggregationFailureState = State{0xE3, "data aggregation has failed", nil}
	ConversionFailureState  = State{0xE4, "data conversion has failed", nil}

	GPIOPanicState    = State{0xF1, "GPIO critical error", nil}
	HandledPanicState = State{0xF2, "panic handled", nil}
)

type (
	// Instruction type (only accessed in this package to prevent custom instructions).
	instruction byte

	// Chan struct containing all channels used by the plugins.
	Chan struct {
		Data        chan<- interface{} // Used to send data (only JSON serializable data)
		Status      chan<- State       // Used to send status
		Instruction <-chan instruction // Used to read instruction from the gakisitor
	}
)

// Instruction list (instructions are immutable to prevent custom instructions)
const (
	StatusPluginInstruction instruction = 0x01 // Send a the current status
	StopPluginInstruction   instruction = 0x0F // Stop the plugin
	StartSessionInstruction instruction = 0x10 // Start a game session (you MUST retrieve user input during this session)
	StopSessionInstruction  instruction = 0x11 // Stop the game session (you MUST stop your retrieving user input)
)

// Plugin structure used to extend the Gakisitor functionality
type Plugin struct {
	// The plugin name. It will be used by the server/game
	// engine to know which plugin the data comes from.
	Name string

	// Start the plugin instance with the plugin profile, logger and channels. You
	// MUST check the profile before starting the process.
	//
	// For more information about plugin, see the package description.
	// For more information about plugin channels, see the Chan structure above.
	Instance func(profile profile.Plugin, log log.Log, channels Chan) error
}

// NewState creates a new custom state.
func NewState(code byte, desc string, raw ...interface{}) State {
	var sraw interface{}
	if len(raw) > 0 {
		sraw = raw[0]
	}
	return State{code, desc, sraw}
}

// Code return the state code.
func (state State) Code() byte { return state.code }

// Code return the state description.
func (state State) Desc() string { return state.desc }

// Code return the state raw.
func (state State) Raw() interface{} { return state.raw }

// Equal compare two states.
func (state State) Equal(o interface{}) bool {
	if s, ok := o.(State); ok {
		return s.code == state.code
	}
	return false
}

// AddRaw add a raw to an existing state.
func (state State) AddRaw(raw interface{}) State {
	return State{state.code, state.desc, raw}
}
