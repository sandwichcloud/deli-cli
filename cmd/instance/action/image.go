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
	Project   *string
	Raw       *bool
	name      *string
	imageName *string
}

func (c *ImageCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("image", "Create an image from an instance").Action(c.action)
	c.name = command.Arg("name", "The instance name").Required().String()
	c.imageName = command.Flag("image-name", "The image name").Required().String()
}

func (c *ImageCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	image, err := c.Application.APIClient.Instance(*c.Project).ActionImage(*c.name, *c.imageName)
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
			logrus.Infof("The instance '%s' is converting to an image '%s'", *c.name, image.Name)
		}
	}
	return nil
}
