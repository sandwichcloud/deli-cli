package member

import (
	"fmt"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type addCommand struct {
	cmd.Command
	imageID   *string
	projectID *string
	raw       *bool
}

func (c *addCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("add", "Add a member").Action(c.action)
	c.imageID = command.Arg("image ID", "The image ID").Required().String()
	c.projectID = command.Arg("project id", "The project ID to add").Required().String()
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
	err = c.Application.APIClient.Image().MemberAdd(*c.imageID, *c.projectID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Project %s was added as a member to image %s", *c.projectID, *c.imageID)
		}
	}
	return nil
}
