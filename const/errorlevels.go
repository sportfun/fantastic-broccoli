package _const

import "fantastic-broccoli/common/types"

const (
	WARNING  types.ErrorLevel = 1 << iota
	ERROR
	CRITICAL
	FATAL
)
