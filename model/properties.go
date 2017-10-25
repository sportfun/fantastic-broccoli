package model

type Properties struct {
	System  SystemDefinition
	Modules []ModuleDefinition
}

type ModuleDefinition struct {
	Name string
	Path string
}

type SystemDefinition struct {
	LinkID     string
	ServerIP   string
	ServerPort int
	ServerSSL  bool
}
