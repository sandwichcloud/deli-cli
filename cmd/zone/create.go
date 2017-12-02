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
	regionID             *string
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
	c.regionID = command.Flag("region-id", "The region this zone belongs to").Required().String()
	c.vmCluster = command.Flag("vm-cluster", "The VMware cluster for this zone").Required().String()
	c.vmDatastore = command.Flag("vm-datastore", "The VMWare datastore for this cluster").Required().String()
	c.vmFolder = command.Flag("vm-folder", "The VMWare VM & Templates folder to keep vms in").String()
	c.coreProvisionPercent = command.Flag("core-provision-percent", "The percentage of cores to provision from the VMWare cluster").Default("100").Int()
	c.ramProvisionPercent = command.Flag("ram-provision-percent", "The percentage of ram to provision from the VMWare cluster").Default("100").Int()

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
	region, err := c.Application.APIClient.Zone().Create(*c.name, *c.regionID, *c.vmCluster, *c.vmDatastore, *c.vmFolder, *c.coreProvisionPercent, *c.ramProvisionPercent)
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
			logrus.Infof("Zone '%s' created with an ID of '%s'", region.Name, region.ID)
		}
	}
	return nil
}
