package flavor

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type deleteCommand struct {
	cmd.Command
	flavorID *string
}

func (c *deleteCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("delete", "Delete a flavor").Action(c.action)
	c.flavorID = command.Arg("flavor ID", "The flavor ID").Required().String()
}

func (c *deleteCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Flavor().Delete(*c.flavorID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Flavor with the id of '%s' is being deleted", *c.flavorID)
		}
	}
	return nil
}
