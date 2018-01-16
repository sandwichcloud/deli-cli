package action

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type VisibilityCommand struct {
	cmd.Command
	Raw     *bool
	imageID *string
	public  *bool
}

func (c *VisibilityCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("visibility", "Change the visibility of an image.").Action(c.action)
	c.imageID = command.Arg("image ID", "The image ID").Required().String()
	c.public = command.Flag("public", "Enable or disable public visibility of an image").Required().NegatableBool()
}

func (c *VisibilityCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	err = c.Application.APIClient.Image().ActionSetVisibility(*c.imageID, *c.public)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Public visibility of the Image '%s' has been changed to '%s'", *c.imageID, strconv.FormatBool(*c.public))
		}
	}
	return nil
}
