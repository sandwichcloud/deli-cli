package auth

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

var raw *bool

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("auth", "Sandwich Cloud authentication commands")
	raw = command.Flag("raw", "Show raw json output").Bool()

	loginCommand := loginCommand{}
	loginCommand.Application = c.Application
	loginCommand.Register(command)

}
