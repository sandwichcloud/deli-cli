package instance

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type deleteCommand struct {
	cmd.Command
	raw        *bool
	instanceID *string
}

func (c *deleteCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("delete", "Delete an instance").Action(c.action)
	c.instanceID = command.Arg("instance ID", "The instance ID").Required().String()
}

func (c *deleteCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Instance().Delete(*c.instanceID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Instance with the id of '%s' is being deleted", *c.instanceID)
		}
	}
	return nil
}
