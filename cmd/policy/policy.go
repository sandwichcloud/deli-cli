package policy

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("policy", "Sandwich Cloud policy commands")
	raw := command.Flag("raw", "Show raw json output").Bool()

	inspectCommand := inspectCommand{raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)
}
