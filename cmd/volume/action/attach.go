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
	Raw        *bool
	volumeID   *string
	instanceID *string
}

func (c *AttachCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("attach", "Attach a Volume to an Instance.").Action(c.action)
	c.volumeID = command.Arg("volume ID", "The volume ID").Required().String()
	c.instanceID = command.Arg("instance ID", "The instance ID").Required().String()
}

func (c *AttachCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Volume().ActionAttach(*c.volumeID, *c.instanceID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("The volume '%s' being attached to the instance '%s'", *c.volumeID, *c.instanceID)
		}
	}
	return nil
}
