package project

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/project/quota"
)

type Command struct {
	cmd.Command
	raw *bool
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("project", "Sandwich Cloud project commands")
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

	quotaCommand := quota.Command{Raw: c.raw}
	quotaCommand.Application = c.Application
	quotaCommand.Register(command)
}
