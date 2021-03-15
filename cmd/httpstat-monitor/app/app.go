package app

import "github.com/rdeusser/cli"

var Debug = &cli.BoolFlag{
	Name:    "debug",
	Desc:    "Show debug logs",
	Default: false,
}
