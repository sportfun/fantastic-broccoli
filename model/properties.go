package model

import "fantastic-broccoli/common/types"

type Properties struct {
	System  SystemDefinition
	Modules []ModuleDefinition
}

type ModuleDefinition struct {
	Name types.Name
	Path types.Path
}

type SystemDefinition struct {
	LinkID types.ID
}
