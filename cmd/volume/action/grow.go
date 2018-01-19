package action

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type GrowCommand struct {
	cmd.Command
	Raw      *bool
	volumeID *string
	size     *int
}

func (c *GrowCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("grow", "Increase the size of a volume.").Action(c.action)
	c.volumeID = command.Arg("volume ID", "The volume ID").Required().String()
	c.size = command.Arg("size", "The new size of the volume").Required().Int()
}

func (c *GrowCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Volume().ActionGrow(*c.volumeID, *c.size)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("The size of the volume '%s' being expanded to '%s'", *c.volumeID, *c.size)
		}
	}
	return nil
}
