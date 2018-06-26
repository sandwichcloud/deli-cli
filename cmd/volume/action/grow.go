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
	Raw     *bool
	Project *string
	name    *string
	size    *int
}

func (c *GrowCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("grow", "Increase the size of a volume.").Action(c.action)
	c.name = command.Arg("volume name", "The volume name").Required().String()
	c.size = command.Arg("size", "The new size of the volume in gigabytes").Required().Int()
}

func (c *GrowCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Volume(*c.Project).ActionGrow(*c.name, *c.size)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("The size of the volume '%s' being expanded to '%s'", *c.name, *c.size)
		}
	}
	return nil
}
