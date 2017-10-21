package model

type State int

const (
	STARTED State = 0x001
	STOPPED State = 0x002
	IDLE    State = 0x004
	WORKING State = 0x008
)

type ErrorType int

const (
	WARNING  ErrorType = 0x001
	ERROR    ErrorType = 0x002
	CRITICAL ErrorType = 0x004
	FATAL    ErrorType = 0x008
)

const (
	MODULE_MANAGER = "ModuleManager"
)