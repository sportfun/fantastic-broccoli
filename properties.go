package fantastic_broccoli

type Properties struct {
	Modules []ModuleDefinition
}

type ModuleDefinition struct {
	Name Name
	Path Path
}