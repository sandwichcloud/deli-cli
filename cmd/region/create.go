package region

import (
	"encoding/json"
	"fmt"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type createCommand struct {
	cmd.Command
	name           *string
	datacenter     *string
	imageDatastore *string
	imageFolder    *string
	raw            *bool
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a region").Action(c.action)
	c.name = command.Arg("name", "The region name").Required().String()
	c.datacenter = command.Flag("datacenter", "The VMWare Datacenter for this region.").Required().String()
	c.imageDatastore = command.Flag("image-datastore", "The VMWare Datastore to keep images in").Required().String()
	c.imageFolder = command.Flag("image-folder", "The VMWare VM & Templates folder to keep images in").String()
}

func (c *createCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetUnScopedToken()
	if err != nil {
		return err
	}
	region, err := c.Application.APIClient.Region().Create(*c.name, *c.datacenter, *c.imageDatastore, *c.imageFolder)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			regionBytes, _ := json.MarshalIndent(region, "", "  ")
			fmt.Println(string(regionBytes))
		} else {
			logrus.Infof("Region '%s' created with an ID of '%s'", region.Name, region.ID)
		}
	}
	return nil
}