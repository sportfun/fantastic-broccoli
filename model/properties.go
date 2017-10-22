package model

import "fantastic-broccoli/common/types"

type Properties struct {
	Modules []ModuleDefinition
}

type ModuleDefinition struct {
	Name types.Name
	Path types.Path
}
