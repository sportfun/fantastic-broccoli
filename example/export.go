package main

import "github.com/xunleii/fantastic-broccoli/common/types/module"

func ExportModule() module.Module {
	return &rpmGenerator{}
}
