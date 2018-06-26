package action

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type AttachCommand struct {
	cmd.Command
	Raw          *bool
	Project      *string
	name         *string
	instanceName *string
}

func (c *AttachCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("attach", "Attach a Volume to an Instance.").Action(c.action)
	c.name = command.Arg("name", "The volume name").Required().String()
	c.instanceName = command.Arg("instance ID", "The instance name").Required().String()
}

func (c *AttachCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Volume(*c.Project).ActionAttach(*c.name, *c.instanceName)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("The volume '%s' being attached to the instance '%s'", *c.name, *c.instanceName)
		}
	}
	return nil
}
