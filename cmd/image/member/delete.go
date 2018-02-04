package member

import (
	"fmt"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type removeCommand struct {
	cmd.Command
	imageID   *string
	projectID *string
	raw       *bool
}

func (c *removeCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("remove", "Remove a member").Action(c.action)
	c.imageID = command.Arg("image ID", "The image ID").Required().String()
	c.projectID = command.Arg("project id", "The project ID to remove").Required().String()
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
	err = c.Application.APIClient.Image().MemberRemove(*c.imageID, *c.projectID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Project %s was removed from image %s", *c.projectID, *c.imageID)
		}
	}
	return nil
}
