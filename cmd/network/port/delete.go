package port

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
	raw    *bool
	portID *string
}

func (c *deleteCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("delete", "Delete a network port").Action(c.action)
	c.portID = command.Arg("network port ID", "The network port ID").Required().String()
}

func (c *deleteCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.NetworkPort().Delete(*c.portID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Network Port with the id of '%s' is being deleted", *c.portID)
		}
	}
	return nil
}
