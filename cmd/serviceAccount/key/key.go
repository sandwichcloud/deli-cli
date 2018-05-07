package key

import (
	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/cmd"
)

type Command struct {
	cmd.Command
	Project bool
	Raw     *bool
}

func (c *Command) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("key", "Service Account Key commands")

	createCommand := createCommand{raw: c.Raw, project: c.Project}
	createCommand.Application = c.Application
	createCommand.Register(command)

	deleteCommand := deleteCommand{raw: c.Raw, project: c.Project}
	deleteCommand.Application = c.Application
	deleteCommand.Register(command)
}
