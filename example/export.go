package main

import "github.com/sportfun/gakisitor/module"

func ExportModule() module.Module {
	return &rpmGenerator{}
}

// Fix issue #20312 (https://github.com/golang/go/issues/20312)
func main() {}
