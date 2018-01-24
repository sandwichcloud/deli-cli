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
	project  bool
	raw      *bool
	roleID   *string
	policies *[]string
}

func (c *updateCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("update", "Update a role").Action(c.action)
	c.roleID = command.Arg("role ID", "The role ID").String()
	c.policies = command.Flag("policy", "The policy to give the role").Required().Strings()
}

func (c *updateCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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
		err = c.Application.APIClient.ProjectRole().Update(*c.roleID, *c.policies)
	} else {
		err = c.Application.APIClient.GlobalRole().Update(*c.roleID, *c.policies)
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
			logrus.Infof("Role with the id of %s has been updated.", *c.roleID)
		}
	}
	return nil
}
