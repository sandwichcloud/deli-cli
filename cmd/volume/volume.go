package volume

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/volume/action"
)

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("volume", "Sandwich Cloud volume commands")
	project := command.Flag("project", "The project to use for this invocation").Required().String()
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{raw: raw, project: project}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{raw: raw, project: project}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw: raw, project: project}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{raw: raw, project: project}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	attachCommand := action.AttachCommand{Raw: raw, Project: project}
	attachCommand.Application = c.Application
	attachCommand.Register(command)

	detachCommand := action.DetachCommand{Raw: raw, Project: project}
	detachCommand.Application = c.Application
	detachCommand.Register(command)

	growCommand := action.GrowCommand{Raw: raw, Project: project}
	growCommand.Application = c.Application
	growCommand.Register(command)

	cloneCommand := action.CloneCommand{Raw: raw, Project: project}
	cloneCommand.Application = c.Application
	cloneCommand.Register(command)

}
