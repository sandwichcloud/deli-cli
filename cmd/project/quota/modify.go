package quota

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type modifyCommand struct {
	cmd.Command
	raw  *bool
	vcpu *int
	ram  *int
	disk *int
}

func (c *modifyCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("modify", "Modify a project's quota").Action(c.action)
	c.vcpu = command.Flag("vcpu", "The number of vcpus the project is allowed to use").Default("-1").Int()
	c.ram = command.Flag("ram", "The amount of ram (in MB) the project is allowed to use").Default("-1").Int()
	c.disk = command.Flag("disk", "The amount of disk (in GB) the project is allowed to use").Default("-1").Int()
}

func (c *modifyCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	tokenInfo, err := c.Application.APIClient.Auth().TokenInfo()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Project().SetQuota(*c.vcpu, *c.ram, *c.disk)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Updated the quota for the project '%s'", tokenInfo.ProjectID)
		}
	}
	return nil
}
