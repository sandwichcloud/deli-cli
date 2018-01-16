package member

import (
	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/cmd"
)

type Command struct {
	cmd.Command
	Raw *bool
}

func (c *Command) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("member", "Image member commands")

	addCommand := addCommand{raw: c.Raw}
	addCommand.Application = c.Application
	addCommand.Register(command)

	removeCommand := removeCommand{raw: c.Raw}
	removeCommand.Application = c.Application
	removeCommand.Register(command)

	listCommand := listCommand{raw: c.Raw}
	listCommand.Application = c.Application
	listCommand.Register(command)
}
