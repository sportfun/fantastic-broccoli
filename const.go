package fantastic_broccoli

type State int

const (
	STARTED State = 1 << iota
	STOPPED
	IDLE
	WORKING
)

type ErrorType int

const (
	WARNING  ErrorType = 1 << iota
	ERROR
	CRITICAL
	FATAL
)

const (
	MODULE_MANAGER = "ModuleManager"
)
