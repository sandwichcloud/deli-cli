package action

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type StopCommand struct {
	cmd.Command
	Project *string
	Raw     *bool
	name    *string
	hard    *bool
	timeout *int
}

func (c *StopCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("stop", "Stop an instance").Action(c.action)
	c.name = command.Arg("name", "The instance name").Required().String()
	c.hard = command.Flag("hard", "Hard stop the instance").Default("false").Bool()
	c.timeout = command.Flag("timeout", "Time in seconds until the instance is hard stopped").Default("60").Int()
}

func (c *StopCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Instance(*c.Project).ActionStop(*c.name, *c.hard, *c.timeout)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Instance '%s' is being stopped", *c.name)
		}
	}
	return nil
}
