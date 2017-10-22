package core

type ServiceState int

const (
	START     ServiceState = iota
	CONFIGURE
	PROCESS
	STOP
)

func (c *Core) serviceErrorHandler(s ServiceState, e error) {
	if e == nil {
		return
	}

	switch s {
	case START:
		startErrorHandler(c, e)
	case CONFIGURE:
		configureErrorHandler(c, e)
	case PROCESS:
		processErrorHandler(c, e)
	case STOP:
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
