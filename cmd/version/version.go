package version

import (
	"fmt"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/config"
)

type VersionContext struct {
}

var Cmd = &cmd.Command{
	Name: "version",
	Help: "Application version",
	Context: &VersionContext{},
	Init: Init,
	Logic: Logic,
}

func Init(c *cmd.Command) {
}

func Logic(c *cmd.Command) {
	fmt.Println(config.Version)
}
