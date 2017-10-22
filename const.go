package fantastic_broccoli

const (
	CORE            Name = "b2fe86de"
	MODULE_SERVICE  Name = "8c7a5db4"
	NETWORK_SERVICE Name = "97c403c4"
)

const (
	STARTED State = 1 << iota
	STOPPED
	IDLE
	WORKING
)

const (
	WARNING  ErrorType = 1 << iota
	ERROR
	CRITICAL
	FATAL
)