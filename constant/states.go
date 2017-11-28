package constant

import "github.com/xunleii/fantastic-broccoli/common/types"

var States = struct {
	Started types.StateType
	Idle    types.StateType
	Working types.StateType
	Stopped types.StateType
	Failed  types.StateType
	Panic   types.StateType
}{
	Started: 0x1,
	Idle:    0x2,
	Working: 0x4,
	Stopped: 0x8,
	Failed:  0x10,
	Panic:   0x12,
}
