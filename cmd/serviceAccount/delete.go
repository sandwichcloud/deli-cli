package serviceAccount

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
	project          bool
	raw              *bool
	serviceAccountID *string
}

func (c *deleteCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("delete", "Delete a service account").Action(c.action)
	c.serviceAccountID = command.Arg("service account ID", "The service account ID").Required().String()
}

func (c *deleteCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if c.project {
		err = c.Application.SetScopedToken()
	} else {
		err = c.Application.SetUnScopedToken()
	}
	if err != nil {
		return err
	}

	if c.project {
		err = c.Application.APIClient.ProjectServiceAccount().Delete(*c.serviceAccountID)
	} else {
		err = c.Application.APIClient.GlobalServiceAccount().Delete(*c.serviceAccountID)
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Service Account with the id of '%s' is being deleted", *c.serviceAccountID)
		}
	}
	return nil
}
