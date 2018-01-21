package member

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type addCommand struct {
	cmd.Command
	raw      *bool
	username *string
	driver   *string
}

func (c *addCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("add", "Add a member to a project").Action(c.action)
	c.username = command.Arg("username", "The username of the user to add").Required().String()
	c.driver = command.Arg("driver", "The driver of the user to add").Required().String()
}

func (c *addCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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
	projectMember, err := c.Application.APIClient.Project().AddMember(*c.username, *c.driver)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			projectMemberBytes, _ := json.MarshalIndent(projectMember, "", "  ")
			fmt.Println(string(projectMemberBytes))
		} else {
			logrus.Infof("Added user (%s/%s) to project %s", projectMember.Username, projectMember.Driver, tokenInfo.ProjectID)
		}
	}
	return nil
}
