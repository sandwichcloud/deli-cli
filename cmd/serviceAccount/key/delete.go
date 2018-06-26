package key

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
	project *string
	raw     *bool
	name    *string
	keyName *string
}

func (c *deleteCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("delete", "Delete a key").Action(c.action)
	c.name = command.Arg("service account name", "The service account name").Required().String()
	c.keyName = command.Arg("key name", "The key name").Required().String()
}

func (c *deleteCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if c.project != nil {
		err = c.Application.APIClient.ProjectServiceAccount(*c.project).DeleteKey(*c.name, *c.keyName)
	} else {
		err = c.Application.APIClient.SystemServiceAccount().DeleteKey(*c.name, *c.keyName)
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
			logrus.Infof("Key '%s' has been delete from the service account '%s'", *c.keyName, *c.name)
		}
	}
	return nil
}
