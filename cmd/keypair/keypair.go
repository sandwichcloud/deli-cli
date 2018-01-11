package keypair

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

var raw *bool

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("keypair", "Sandwich Cloud keypair commands")
	raw = command.Flag("raw", "Show raw json output").Bool()

	generateCommand := generateCommand{raw: raw}
	generateCommand.Application = c.Application
	generateCommand.Register(command)

	importCommand := importCommand{raw: raw}
	importCommand.Application = c.Application
	importCommand.Register(command)

	inspectCommand := inspectCommand{raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{raw: raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)
}
