package constant

var EntityNames = struct {
	Core string
	Services struct {
		Module  string
		Network string
	}
}{
	Core: "core",
	Services: struct {
		Module  string
		Network string
	}{
		Module:  "module_manager",
		Network: "network_manager",
	},
}
