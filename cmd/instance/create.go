package instance

import (
	"encoding/json"
	"fmt"

	"errors"

	"strings"

	"io/ioutil"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type createCommand struct {
	cmd.Command
	raw        *bool
	name       *string
	imageID    *string
	networkID  *string
	keypairIDs *[]string
	tags       *map[string]string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create an image").Action(c.action)
	c.name = command.Arg("name", "The image name").Required().String()
	c.imageID = command.Flag("imageID", "The image to launch the instance from").Required().String()
	c.networkID = command.Flag("networkID", "The network to attach the instance to").Required().String()
	c.keypairIDs = command.Flag("keypairID", "An ID of a keypair to add to the instance").Strings()
	c.tags = command.Flag("tag", "A metadata tag to add to the instance").StringMap()
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

	tags := *c.tags
	for k, v := range tags {
		if strings.HasPrefix(v, "@") {
			tagBytes, err := ioutil.ReadFile(v[1:])
			if err != nil {
				return err
			}

			tags[k] = string(tagBytes)
		}
	}

	instance, err := c.Application.APIClient.Instance().Create(*c.name, *c.imageID, *c.networkID, *c.keypairIDs, tags)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			instanceBytes, _ := json.MarshalIndent(instance, "", "  ")
			fmt.Println(string(instanceBytes))
		} else {
			logrus.Infof("Instance '%s' created with an ID of '%s'", instance.Name, instance.ID)
		}
	}
	return nil
}
