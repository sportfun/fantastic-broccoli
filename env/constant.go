package env

const (
	WarningLevel  = "Warning"
	ErrorLevel    = "Error"
	CriticalLevel = "Critical"
	FatalLevel    = "Fatal"
)

const (
	UndefinedState      = 0
	StartedState   byte = 1 << iota
	IdleState
	WorkingState
	StoppedState
	PanicState
)

const (
	CoreEntity           = "core"
	ModuleServiceEntity  = "module_manager"
	NetworkServiceEntity = "network_manager"
)

const (
	StateCmd          = "state"
	LinkCmd           = "link"
	StartSessionCmd   = "start_session"
	EndSessionCmd     = "end_session"
	RestartServiceCmd = "restart_service"
)
