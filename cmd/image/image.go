package image

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/image/action"
	"github.com/sandwichcloud/deli-cli/cmd/image/member"
)

type Command struct {
	cmd.Command
}

var raw *bool

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("image", "Sandwich Cloud image commands")
	raw = command.Flag("raw", "Show raw json output").Bool()

	importCommand := importCommand{}
	importCommand.Application = c.Application
	importCommand.Register(command)

	inspectCommand := inspectCommand{}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	memberCommand := member.Command{Raw: raw}
	memberCommand.Application = c.Application
	memberCommand.Register(command)

	visibilityCommand := action.VisibilityCommand{Raw: raw}
	visibilityCommand.Application = c.Application
	visibilityCommand.Register(command)
}
