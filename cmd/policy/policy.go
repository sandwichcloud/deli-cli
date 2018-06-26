package policy

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
	command := app.CLIApp.Command("system-policy", "Sandwich Cloud system policy commands")
	raw := command.Flag("raw", "Show raw json output").Bool()

	inspectCommand := inspectCommand{raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	setCommand := setCommand{raw: raw}
	setCommand.Application = c.Application
	setCommand.Register(command)
}

func (c *Command) project(app *cmd.Application) {
	command := app.CLIApp.Command("project-policy", "Sandwich Cloud project policy commands")

	project := command.Flag("project", "The project to use for this invocation").Required().String()
	raw := command.Flag("raw", "Show raw json output").Bool()

	inspectCommand := inspectCommand{project: project, raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	setCommand := setCommand{project: project, raw: raw}
	setCommand.Application = c.Application
	setCommand.Register(command)
}
