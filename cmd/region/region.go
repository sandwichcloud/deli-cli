package region

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("region", "Sandwich Cloud region commands")
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{raw:raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{raw:raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw:raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{raw:raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	updateCommand := updateCommand{raw:raw}
	updateCommand.Application = c.Application
	updateCommand.Register(command)

}
