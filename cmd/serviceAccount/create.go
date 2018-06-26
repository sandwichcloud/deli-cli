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
	project *string
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
	if err != nil {
		return err
	}

	var serviceAccount *api.ServiceAccount
	if c.project != nil {
		serviceAccount, err = c.Application.APIClient.ProjectServiceAccount(*c.project).Create(*c.name)
	} else {
		serviceAccount, err = c.Application.APIClient.SystemServiceAccount().Create(*c.name)
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
			logrus.Infof("Service Account '%s' created", serviceAccount.Name)
		}
	}
	return nil
}
