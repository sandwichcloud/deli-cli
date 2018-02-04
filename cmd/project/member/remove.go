package member

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type removeCommand struct {
	cmd.Command
	raw             *bool
	projectMemberID *string
}

func (c *removeCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("remove", "Remove a member from a project").Action(c.action)
	c.projectMemberID = command.Arg("project member ID", "The project member ID").Required().String()
}

func (c *removeCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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
	err = c.Application.APIClient.Project().RemoveMember(*c.projectMemberID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Removed member %s from project %s", *c.projectMemberID, tokenInfo.ProjectID)
		}
	}
	return nil
}
