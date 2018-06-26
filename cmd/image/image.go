package image

import (
	"github.com/sandwichcloud/deli-cli/cmd"
)

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("image", "Sandwich Cloud image commands")

	project := command.Flag("project", "The project to use for this invocation").Required().String()
	raw := command.Flag("raw", "Show raw json output").Bool()

	importCommand := importCommand{project: project, raw: raw}
	importCommand.Application = c.Application
	importCommand.Register(command)

	inspectCommand := inspectCommand{project: project, raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{project: project, raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{project: project, raw: raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)
}
