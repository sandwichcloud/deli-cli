package zone

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
	name                 *string
	region               *string
	vmCluster            *string
	vmDatastore          *string
	vmFolder             *string
	coreProvisionPercent *int
	ramProvisionPercent  *int
	raw                  *bool
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a zone").Action(c.action)
	c.name = command.Arg("name", "The zone name").Required().String()
	c.region = command.Flag("region", "The region this zone belongs to").Required().String()
	c.vmCluster = command.Flag("vm-cluster", "The VMware cluster for this zone").Required().String()
	c.vmDatastore = command.Flag("vm-datastore", "The VMware datastore for this zone").Required().String()
	c.vmFolder = command.Flag("vm-folder", "The VMware VM & Templates folder to keep vms in").String()
	c.coreProvisionPercent = command.Flag("core-provision-percent", "The percentage of cores to provision from the VMware cluster").Default("1600").Int()
	c.ramProvisionPercent = command.Flag("ram-provision-percent", "The percentage of ram to provision from the VMware cluster").Default("150").Int()

}

func (c *createCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	zone, err := c.Application.APIClient.Zone().Create(*c.name, *c.region, *c.vmCluster, *c.vmDatastore, *c.vmFolder, *c.coreProvisionPercent, *c.ramProvisionPercent)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			regionBytes, _ := json.MarshalIndent(zone, "", "  ")
			fmt.Println(string(regionBytes))
		} else {
			logrus.Infof("Zone '%s' created", zone.Name)
		}
	}
	return nil
}
