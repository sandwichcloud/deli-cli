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

type CloneCommand struct {
	cmd.Command
	Raw      *bool
	volumeID *string
	name     *string
}

func (c *CloneCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("clone", "Clone a volume").Action(c.action)
	c.volumeID = command.Arg("volume ID", "The volume ID").Required().String()
	c.name = command.Arg("name", "The name of the new volume").Required().String()
}

func (c *CloneCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	volume, err := c.Application.APIClient.Volume().ActionClone(*c.volumeID, *c.name)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.Raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.Raw {
			volumeBytes, _ := json.MarshalIndent(volume, "", "  ")
			fmt.Println(string(volumeBytes))
		} else {
			logrus.Infof("Volume '%s' is being cloned to '%s'", *c.volumeID, volume.ID.String())
		}
	}
	return nil
}
