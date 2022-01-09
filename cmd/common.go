package cmd

import (
	"github.com/akamensky/argparse"
)

type InitHandler func(cmd *Command)
type LogicHandler func(cmd *Command)

type Command struct {
	Name		string
	Help		string
	Cmd			*argparse.Command
	Context		interface{}
	Init		InitHandler
	Logic		LogicHandler
}

func (c *Command) Register(parser *argparse.Parser) {
	c.Cmd = parser.NewCommand(c.Name, c.Help)
	c.Init(c)
}

func (c *Command) IsInvoked() bool {
	return c.Cmd.Happened()
}

func (c *Command) Service() {
	c.Logic(c)
}
