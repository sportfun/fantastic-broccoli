package core

type ErrorType int

const (
	ModuleStart     ErrorType = iota
	ModuleConfigure
	ModuleProcess
	ModuleStop
)

func (c *Core) serviceErrorHandler(t ErrorType, e error) {
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
		processErrorHandler(c, e)
	}
}

func startErrorHandler(c *Core, e error) {
}

func stopErrorHandler(c *Core, e error) {
}

func configureErrorHandler(c *Core, e error) {
}

func processErrorHandler(c *Core, e error) {
}
