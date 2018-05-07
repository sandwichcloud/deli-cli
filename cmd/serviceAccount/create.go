package serviceAccount

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type createCommand struct {
	cmd.Command
	project bool
	raw     *bool
	name    *string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a service account").Action(c.action)
	c.name = command.Arg("name", "The service account name").Required().String()
}

func (c *createCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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

	var serviceAccount *api.ServiceAccount
	if c.project {
		serviceAccount, err = c.Application.APIClient.ProjectServiceAccount().Create(*c.name)
	} else {
		serviceAccount, err = c.Application.APIClient.GlobalServiceAccount().Create(*c.name)
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			serviceAccountBytes, _ := json.MarshalIndent(serviceAccount, "", "  ")
			fmt.Println(string(serviceAccountBytes))
		} else {
			logrus.Infof("Service Account '%s' created with an ID of '%s'", serviceAccount.Name, serviceAccount.ID)
		}
	}
	return nil
}
