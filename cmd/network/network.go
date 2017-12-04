package network

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/network/port"
)

type Command struct {
	cmd.Command
}

var raw *bool

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("network", "Sandwich Cloud network commands")
	raw = command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	portCommand := port.Command{Raw: raw}
	portCommand.Application = c.Application
	portCommand.Register(command)
}
