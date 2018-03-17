package plugin

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

	IdleState      = State{0x11, "currently idle", nil}
	InSessionState = State{0x12, "currently in session", nil}
	PausedState    = State{0x13, "currently paused", nil}

	GPIODisconnectionState  = State{0xE1, "GPIO has been disconnected", nil}
	GPIOFailureState        = State{0xE2, "GPIO reading has failed", nil}
	AggregationFailureState = State{0xE3, "data aggregation has failed", nil}
	ConversionFailureState  = State{0xE4, "data conversion has failed", nil}

	GPIOPanicState    = State{0xF1, "GPIO critical error", nil}
	HandledPanicState = State{0xF2, "panic handled", nil}
)

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

// Equal compare two states (compare only state codes).
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
