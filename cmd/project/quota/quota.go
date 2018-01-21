package quota

import (
	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/cmd"
)

type Command struct {
	cmd.Command
	Raw *bool
}

func (c *Command) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("quota", "Project quota commands")

	inspectCommand := inspectCommand{raw: c.Raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	modifyCommand := modifyCommand{raw: c.Raw}
	modifyCommand.Application = c.Application
	modifyCommand.Register(command)

}
