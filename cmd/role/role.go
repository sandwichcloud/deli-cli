package role

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("role", "Sandwich Cloud role commands")
	command.Flag("raw", "Show raw json output").Bool()
}
