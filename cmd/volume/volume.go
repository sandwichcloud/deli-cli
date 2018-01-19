package volume

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/volume/action"
)

type Command struct {
	cmd.Command
	raw *bool
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("volume", "Sandwich Cloud volume commands")
	c.raw = command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{raw: c.raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{raw: c.raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw: c.raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{raw: c.raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	attachCommand := action.AttachCommand{Raw: c.raw}
	attachCommand.Application = c.Application
	attachCommand.Register(command)

	detachCommand := action.DetachCommand{Raw: c.raw}
	detachCommand.Application = c.Application
	detachCommand.Register(command)

	growCommand := action.GrowCommand{Raw: c.raw}
	growCommand.Application = c.Application
	growCommand.Register(command)

	cloneCommand := action.CloneCommand{Raw: c.raw}
	cloneCommand.Application = c.Application
	cloneCommand.Register(command)

}
