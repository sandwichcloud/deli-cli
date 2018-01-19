package action

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type DetachCommand struct {
	cmd.Command
	Raw      *bool
	volumeID *string
}

func (c *DetachCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("detach", "Detach a Volume from an Instance.").Action(c.action)
	c.volumeID = command.Arg("volume ID", "The volume ID").Required().String()
}

func (c *DetachCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Volume().ActionDetach(*c.volumeID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("The volume '%s' being detached from the instance", *c.volumeID)
		}
	}
	return nil
}
