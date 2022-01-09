package config

import (
	"os"
	"fmt"
	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/common"
)

type ConfigContext struct {
	Show  *bool
	Store *bool
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
	
	ctx.Show = c.Cmd.Flag("s", "show", &argparse.Options{
		Required: false,
		Help:     `Show configuration.`,
	})

	ctx.Store = c.Cmd.Flag("S", "store", &argparse.Options{
		Required: false,
		Help:     `Store configuration.`,
	})
}

func Logic(c *cmd.Command) {
	ctx := c.Context.(*ConfigContext)

	if (*ctx.Show) {
		fmt.Print(common.G)
	}

	if (*ctx.Store) {
		file, err := os.Create(common.G.ConfigFile)
		if err != nil {
			fmt.Println("Cannot create config file:", err.Error())
			return
		}
		defer file.Close()

		file.WriteString(common.G.String())
		fmt.Println("Config file saved.")
	}

	if (!*ctx.Show && !*ctx.Store) {
		fmt.Print(c.Cmd.Usage(nil))
	}
}
