package role

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
	project     *string
	raw         *bool
	name        *string
	permissions *[]string
}

func (c *updateCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("update", "Update a role").Action(c.action)
	c.name = command.Arg("role name", "The role name").String()
	c.permissions = command.Flag("permission", "The permission to give the role").Required().Strings()
}

func (c *updateCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if c.project != nil {
		err = c.Application.APIClient.ProjectRole(*c.project).Update(*c.name, *c.permissions)
	} else {
		err = c.Application.APIClient.SystemRole().Update(*c.name, *c.permissions)
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
			logrus.Infof("Role '%s' has been updated.", *c.name)
		}
	}
	return nil
}
