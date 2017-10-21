package model

type Path string
type Name string

type Properties struct {
	Modules []Module
}

type Module struct {
	Name Name
	Path Path
}