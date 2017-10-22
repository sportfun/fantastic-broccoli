package _const

import "fantastic-broccoli/common/types"

const (
	WARNING  types.ErrorType = 1 << iota
	ERROR
	CRITICAL
	FATAL
)
