package zone

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
	raw         *bool
	name        *string
	schedulable *bool
}

func (c *updateCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("update", "").Action(c.action)
	c.name = command.Arg("name", "The zone ID").Required().String()
	c.schedulable = command.Flag("schedulable", "Enable or disable the ability to schedule workloads in the zone").Required().NegatableBool()
}

func (c *updateCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Zone().ActionSchedule(*c.name, *c.schedulable)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Zone '%s' has been updated.", *c.name)
		}
	}
	return nil
}
