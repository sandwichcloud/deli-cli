package instance

import (
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/cmd/instance/action"
)

type Command struct {
	cmd.Command
}

func (c *Command) Register(app *cmd.Application) {
	c.Application = app
	command := app.CLIApp.Command("instance", "Sandwich Cloud instance commands")

	project := command.Flag("project", "The project to use for this invocation").Required().String()
	raw := command.Flag("raw", "Show raw json output").Bool()

	createCommand := createCommand{project: project, raw: raw}
	createCommand.Application = c.Application
	createCommand.Register(command)

	inspectCommand := inspectCommand{project: project, raw: raw}
	inspectCommand.Application = c.Application
	inspectCommand.Register(command)

	listCommand := listCommand{project: project, raw: raw}
	listCommand.Application = c.Application
	listCommand.Register(command)

	deleteCommand := deleteCommand{project: project, raw: raw}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)

	actionImage := action.ImageCommand{Project: project, Raw: raw}
	actionImage.Application = c.Application
	actionImage.Register(command)

	actionRestart := action.RestartCommand{Project: project, Raw: raw}
	actionRestart.Application = c.Application
	actionRestart.Register(command)

	actionStart := action.StartCommand{Project: project, Raw: raw}
	actionStart.Application = c.Application
	actionStart.Register(command)

	actionStop := action.StopCommand{Project: project, Raw: raw}
	actionStop.Application = c.Application
	actionStop.Register(command)

}
