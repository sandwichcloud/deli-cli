package port

import (
	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/cmd"
)

type Command struct {
	cmd.Command
	Raw *bool
}

func (c *Command) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("port", "Sandwich Cloud network port commands")

	project := command.Flag("project", "The project to use for this invocation").Required().String()

	inspectCommand := inspectCommand{project: project, raw: c.Raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{project: project, raw: c.Raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{project: project, raw: c.Raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)
}
