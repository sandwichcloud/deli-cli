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
	raw    *bool
	name   *string
	zoneID *string
	size   *int
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create an volume").Action(c.action)
	c.name = command.Arg("name", "The volume name").Required().String()
	c.zoneID = command.Flag("zone-id", "The zone to create the volume in").Required().String()
	c.size = command.Flag("size", "The size of the volume").Required().Int()
}

func (c *createCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}

	volume, err := c.Application.APIClient.Volume().Create(*c.name, *c.zoneID, *c.size)
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
			logrus.Infof("Volume '%s' created with an ID of '%s'", volume.Name, volume.ID)
		}
	}
	return nil
}
