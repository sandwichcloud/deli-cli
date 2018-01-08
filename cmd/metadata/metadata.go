package metadata

import "github.com/sandwichcloud/deli-cli/cmd"

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("metadata", "Sandwich Cloud metadata commands")

	metaDataCommand := metaDataCommand{}
	metaDataCommand.Application = c.Application
	metaDataCommand.Register(command)

	networkDataCommand := networkDataCommand{}
	networkDataCommand.Application = c.Application
	networkDataCommand.Register(command)

	userDataCommand := userDataCommand{}
	userDataCommand.Application = c.Application
	userDataCommand.Register(command)
}
