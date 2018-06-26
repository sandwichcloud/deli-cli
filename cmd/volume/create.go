package volume

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
	raw     *bool
	project *string
	name    *string
	zone    *string
	size    *int
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create an volume").Action(c.action)
	c.name = command.Arg("name", "The volume name").Required().String()
	c.zone = command.Flag("zone", "The zone to create the volume in").Required().String()
	c.size = command.Flag("size", "The size of the volume in gigabytes").Required().Int()
}

func (c *createCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}

	volume, err := c.Application.APIClient.Volume(*c.project).Create(*c.name, *c.zone, *c.size)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			volumeBytes, _ := json.MarshalIndent(volume, "", "  ")
			fmt.Println(string(volumeBytes))
		} else {
			logrus.Infof("Volume '%s' created", volume.Name)
		}
	}
	return nil
}
