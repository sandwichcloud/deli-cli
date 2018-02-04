package serviceAccount

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type updateCommand struct {
	cmd.Command
	raw              *bool
	serviceAccountID *string
	roles            *[]string
}

func (c *updateCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("update", "Update a service account's roles").Action(c.action)
	c.serviceAccountID = command.Arg("service account ID", "The service account ID").Required().String()
	c.roles = command.Flag("role-id", "The role to give the service account").Required().Strings()
}

func (c *updateCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.ServiceAccount().Update(*c.serviceAccountID, *c.roles)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Updated the roles of service account %s", *c.serviceAccountID)
		}
	}
	return nil
}
