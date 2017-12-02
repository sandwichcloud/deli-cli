package action

import (
	"errors"
	"fmt"

	"strconv"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type ResetStateCommand struct {
	cmd.Command
	Raw        *bool
	instanceID *string
	active     *bool
}

func (c *ResetStateCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("reset-state", "Reset the state of an instance").Action(c.action)
	c.instanceID = command.Arg("instance ID", "The instance ID").Required().String()
	c.active = command.Flag("active", "Set the state to ACTIVE instead of ERROR").Default(strconv.FormatBool(false)).Bool()
}

func (c *ResetStateCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Instance().ActionResetState(*c.instanceID, *c.active)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("The state of the instance with the id of '%s' has been reset.", *c.instanceID)
		}
	}
	return nil
}
