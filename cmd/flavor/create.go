package flavor

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
	name  *string
	vcpus *int
	ram   *int
	disk  *int
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a flavor").Action(c.action)
	c.name = command.Arg("name", "The flavor name").Required().String()
	c.vcpus = command.Flag("vcpus", "Number of VCPUs for this flavor").Required().Int()
	c.ram = command.Flag("ram", "Amount of ram in megabytes for this flavor").Required().Int()
	c.disk = command.Flag("disk", "Size in gigabytes of the root disk for this flavor").Required().Int()
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
	flavor, err := c.Application.APIClient.Flavor().Create(*c.name, *c.vcpus, *c.ram, *c.disk)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			flavorBytes, _ := json.MarshalIndent(flavor, "", "  ")
			fmt.Println(string(flavorBytes))
		} else {
			logrus.Infof("Flavor '%s' created with an ID of '%s'", flavor.Name, flavor.ID)
		}
	}
	return nil
}
