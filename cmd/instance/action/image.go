package action

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type ImageCommand struct {
	cmd.Command
	Raw        *bool
	instanceID *string
	name       *string
}

func (c *ImageCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("image", "Create an image from an instance").Action(c.action)
	c.instanceID = command.Arg("instance ID", "The instance ID").Required().String()
	c.name = command.Flag("name", "The image name").Required().String()
}

func (c *ImageCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	image, err := c.Application.APIClient.Instance().ActionImage(*c.instanceID, *c.name)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			imageBytes, _ := json.MarshalIndent(image, "", "  ")
			fmt.Println(string(imageBytes))
		} else {
			logrus.Infof("The instance with the ID of '%s' is converting to an image called '%s' with an ID of '%s'", *c.instanceID, image.Name, image.ID)
		}
	}
	return nil
}
