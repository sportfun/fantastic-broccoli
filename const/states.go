package _const

import "fantastic-broccoli/common/types"

const (
	STARTED types.State = 1 << iota
	STOPPED
	IDLE
	WORKING
)
