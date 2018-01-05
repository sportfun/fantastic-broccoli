package main

import (
	"flag"

	"github.com/sportfun/gakisitor/kernel"

	_ "github.com/sportfun/gakisitor/service/network"

	_ "github.com/sportfun/gakisitor/service/module"
)

var Core kernel.Core

func init() {
	fileConf := Core.Parameter("config").(*string)
	flag.StringVar(fileConf, "c", "/etc/gakisitor/config.json", "path where file is configured (shorthand)")
	flag.StringVar(fileConf, "config", "/etc/gakisitor/config.json", "path where file is configured")

	retryMax := Core.Parameter("retry_max").(*int)
	flag.IntVar(retryMax, "m", 5, "number max of retry before failure (shorthand)")
	flag.IntVar(retryMax, "max-retry", 5, "number max of retry before failure")
}

func main() {
	flag.Parse()
	Core.Init()
	Core.Run()
}
