package action

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type StartCommand struct {
	cmd.Command
	Project *string
	Raw     *bool
	name    *string
}

func (c *StartCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("start", "Start an instance").Action(c.action)
	c.name = command.Arg("name", "The instance name").Required().String()
}

func (c *StartCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Instance(*c.Project).ActionStart(*c.name)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Instance '%s' is being started", *c.name)
		}
	}
	return nil
}
