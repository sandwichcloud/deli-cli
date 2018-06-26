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
	project     *string
	raw         *bool
	name        *string
	permissions *[]string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a role").Action(c.action)
	c.name = command.Arg("name", "The role name").Required().String()
	c.permissions = command.Flag("permission", "Permission to add to the role").Strings()
}

func (c *createCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	permissions := make([]string, 0)
	if *c.permissions != nil {
		permissions = *c.permissions
	}
	var role *api.Role
	if c.project != nil {
		role, err = c.Application.APIClient.ProjectRole(*c.project).Create(*c.name, permissions)
	} else {
		role, err = c.Application.APIClient.SystemRole().Create(*c.name, permissions)
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
			logrus.Infof("Role '%s' created", role.Name)
		}
	}
	return nil
}
