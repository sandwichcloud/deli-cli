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
	Raw        *bool
	instanceID *string
	hard       *bool
	timeout    *int
}

func (c *StopCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("stop", "Stop an instance").Action(c.action)
	c.instanceID = command.Arg("instance ID", "The instance ID").Required().String()
	c.hard = command.Flag("hard", "Hard stop the instance").Default("false").Bool()
	c.timeout = command.Flag("timeout", "Time in seconds until the instance is hard stopped").Default("60").Int()
}

func (c *StopCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Instance().ActionStop(*c.instanceID, *c.hard, *c.timeout)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Instance with the id of '%s' is being stopped", *c.instanceID)
		}
	}
	return nil
}
