package config

import (
	"fmt"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/common"
)

type ConfigContext struct {
	Dump *bool
}

var Cmd = &cmd.Command{
	Name: "config",
	Help: "Configuration management",
	Context: &ConfigContext{},
	Init: Init,
	Logic: Logic,
}

func Init(c *cmd.Command) {
	ctx := c.Context.(*ConfigContext)
	
	ctx.Dump = c.Cmd.Flag("d", "dump", &argparse.Options{
		Required: false,
		Help:     `Dump configuration file.`,
	})
}

func Logic(c *cmd.Command) {
	ctx := c.Context.(*ConfigContext)

	if (*ctx.Dump) {
		fmt.Print(common.G.Settings)
		return
	}

	fmt.Print(c.Cmd.Usage(nil))
}
