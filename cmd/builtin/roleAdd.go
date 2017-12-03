package builtin

import (
	"fmt"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type roleAddCommand struct {
	cmd.Command
	raw    *bool
	userID *string
	role   *string
}

func (c *roleAddCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("role-add", "Add a role to a user").Action(c.action)
	c.userID = command.Flag("userId", "ID of the user").Required().String()
	c.role = command.Arg("role name", "The role to add").Required().String()
}

func (c *roleAddCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetUnScopedToken()
	if err != nil {
		return err
	}

	err = c.Application.APIClient.BuiltInAuth().AddRole(*c.userID, *c.role)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Role has been added to the user %s", *c.userID)
		}
	}
	return nil
}
