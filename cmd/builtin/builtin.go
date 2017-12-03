package builtin

import (
	"github.com/sandwichcloud/deli-cli/cmd"
)

type Command struct {
	cmd.Command
	raw *bool
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("builtin", "Sandwich Cloud built-in auth commands")
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

	passwordCommand := passwordCommand{raw: c.raw}
	passwordCommand.Application = c.Application
	passwordCommand.Register(command)

	roleAddCommand := roleAddCommand{raw: c.raw}
	roleAddCommand.Application = c.Application
	roleAddCommand.Register(command)

	roleRemoveCommand := roleRemoveCommand{raw: c.raw}
	roleRemoveCommand.Application = c.Application
	roleRemoveCommand.Register(command)
}
