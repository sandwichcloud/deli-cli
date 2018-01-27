package instance

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/instance/action"
)

type Command struct {
	cmd.Command
	raw *bool
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("instance", "Sandwich Cloud instance commands")
	c.raw = command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{raw: c.raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{raw: c.raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{raw: c.raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{raw: c.raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	actionImage := action.ImageCommand{Raw: c.raw}
	actionImage.Application = c.Application
	actionImage.Register(command)

	actionRestart := action.RestartCommand{Raw: c.raw}
	actionRestart.Application = c.Application
	actionRestart.Register(command)

	actionStart := action.StartCommand{Raw: c.raw}
	actionStart.Application = c.Application
	actionStart.Register(command)

	actionStop := action.StopCommand{Raw: c.raw}
	actionStop.Application = c.Application
	actionStop.Register(command)

}
