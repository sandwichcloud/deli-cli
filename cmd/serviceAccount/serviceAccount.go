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
	c.system(app)
	c.project(app)

}

func (c *Command) system(app *cmd.Application) {
	command := app.CLIApp.Command("system-service-account", "Sandwich Cloud system service account commands")
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

	keyCommand := key.Command{Raw: raw}
	keyCommand.Application = c.Application
	keyCommand.Register(command)
}

func (c *Command) project(app *cmd.Application) {
	command := app.CLIApp.Command("project-service-account", "Sandwich Cloud project service account commands")

	project := command.Flag("project", "The project to use for this invocation").Required().String()
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{project: project, raw: raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{project: project, raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{project: project, raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{project: project, raw: raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	keyCommand := key.Command{Project: project, Raw: raw}
	keyCommand.Application = c.Application
	keyCommand.Register(command)
}
