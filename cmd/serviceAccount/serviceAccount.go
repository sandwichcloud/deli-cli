package serviceAccount

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/serviceAccount/key"
)

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	c.global(app)
	c.project(app)

}

func (c *Command) global(app *cmd.Application) {
	command := app.CLIApp.Command("global-service-account", "Sandwich Cloud global service account commands")
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{raw: raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{raw: raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	updateCommand := updateCommand{raw: raw}
	updateCommand.Application = c.Application
	updateCommand.Register(command)

	keyCommand := key.Command{Raw: raw}
	keyCommand.Application = c.Application
	keyCommand.Register(command)
}

func (c *Command) project(app *cmd.Application) {
	command := app.CLIApp.Command("project-service-account", "Sandwich Cloud project service account commands")
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{raw: raw, project: true}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{raw: raw, project: true}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw: raw, project: true}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{raw: raw, project: true}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	updateCommand := updateCommand{raw: raw, project: true}
	updateCommand.Application = c.Application
	updateCommand.Register(command)

	keyCommand := key.Command{Raw: raw, Project: true}
	keyCommand.Application = c.Application
	keyCommand.Register(command)
}
