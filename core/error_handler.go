package core

type errorType int

const (
	ModuleStart     errorType = iota
	ModuleConfigure
	ModuleProcess
	ModuleStop
)

func (c *Core) serviceErrorHandler(t errorType, e error) {
	if e == nil {
		return
	}

	switch t {
	case ModuleStart:
		startErrorHandler(c, e)
	case ModuleConfigure:
		configureErrorHandler(c, e)
	case ModuleProcess:
		processErrorHandler(c, e)
	case ModuleStop:
		stopErrorHandler(c, e)
	}
}

func startErrorHandler(c *Core, e error) {
	//TODO: Start error handler
}

func stopErrorHandler(c *Core, e error) {
	//TODO: Stop error handler
}

func configureErrorHandler(c *Core, e error) {
	//TODO: Configure error handler
}

func processErrorHandler(c *Core, e error) {
	//TODO: Process error handler
}
