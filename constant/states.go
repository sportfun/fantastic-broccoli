package constant

var States = struct {
	Started byte
	Idle    byte
	Working byte
	Stopped byte
	Panic   byte
}{
	Started: 0x1,
	Idle:    0x2,
	Working: 0x4,
	Stopped: 0x8,
	Panic:   0x10,
}
