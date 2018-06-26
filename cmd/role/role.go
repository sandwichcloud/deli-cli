package role

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	c.system(app)
	c.project(app)
}

func (c *Command) system(app *cmd.Application) {
	command := app.CLIApp.Command("system-role", "Sandwich Cloud system role commands")
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{raw: raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	listCommand := listCommand{raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	inspectCommand := inspectCommand{raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	deleteCommand := deleteCommand{raw: raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	updateCommand := updateCommand{raw: raw}
	updateCommand.Application = c.Application
	updateCommand.Register(command)
}

func (c *Command) project(app *cmd.Application) {
	command := app.CLIApp.Command("project-role", "Sandwich Cloud project role commands")

	project := command.Flag("project", "The project to use for this invocation").Required().String()
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{project: project, raw: raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	listCommand := listCommand{project: project, raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	inspectCommand := inspectCommand{project: project, raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	deleteCommand := deleteCommand{project: project, raw: raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	updateCommand := updateCommand{project: project, raw: raw}
	updateCommand.Application = c.Application
	updateCommand.Register(command)
}
