package serviceAccount

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("service-account", "Sandwich Cloud service account commands")
	command.Flag("raw", "Show raw json output").Bool()
}
