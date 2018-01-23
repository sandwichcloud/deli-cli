package role

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
	project  bool
	raw      *bool
	name     *string
	policies *[]string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a role").Action(c.action)
	c.name = command.Arg("name", "The role name").Required().String()
	c.policies = command.Flag("policy", "Policy to add to the role").Strings()
}

func (c *createCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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
	policies := make([]string, 0)
	if *c.policies != nil {
		policies = *c.policies
	}
	var role *api.Role
	if c.project {
		role, err = c.Application.APIClient.ProjectRole().Create(*c.name, policies)
	} else {
		role, err = c.Application.APIClient.GlobalRole().Create(*c.name, policies)
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			roleBytes, _ := json.MarshalIndent(role, "", "  ")
			fmt.Println(string(roleBytes))
		} else {
			logrus.Infof("Role '%s' created with an ID of '%s'", role.Name, role.ID)
		}
	}
	return nil
}
