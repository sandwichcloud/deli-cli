package action

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type DetachCommand struct {
	cmd.Command
	Raw     *bool
	Project *string
	name    *string
}

func (c *DetachCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("detach", "Detach a Volume from an Instance.").Action(c.action)
	c.name = command.Arg("volume name", "The volume name").Required().String()
}

func (c *DetachCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Volume(*c.Project).ActionDetach(*c.name)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("The volume '%s' being detached from the instance", *c.name)
		}
	}
	return nil
}
