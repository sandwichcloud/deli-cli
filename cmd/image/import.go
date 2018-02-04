package image

import (
	"encoding/json"
	"fmt"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type importCommand struct {
	cmd.Command
	name     *string
	regionID *string
	fileName *string
}

func (c *importCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("import", "Import an image").Action(c.action)
	c.name = command.Arg("name", "The image name").Required().String()
	c.regionID = command.Flag("region-id", "The region to create the image in").Required().String()
	c.fileName = command.Arg("file name", "The image's file name").Required().String()
}

func (c *importCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	image, err := c.Application.APIClient.Image().Create(*c.name, *c.regionID, *c.fileName)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			imageBytes, _ := json.MarshalIndent(image, "", "  ")
			fmt.Println(string(imageBytes))
		} else {
			logrus.Infof("Image '%s' imported with an ID of '%s'", image.Name, image.ID)
		}
	}
	return nil
}
