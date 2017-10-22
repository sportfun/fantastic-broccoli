package fantastic_broccoli

type Path string
type Name string

type Properties struct {
	Modules []ModuleDefinition
}

type ModuleDefinition struct {
	Name Name
	Path Path
}