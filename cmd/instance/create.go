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
	raw              *bool
	name             *string
	regionID         *string
	zoneID           *string
	imageID          *string
	serviceAccountID *string
	networkID        *string
	flavorId         *string
	disk             *int
	keypairIDs       *[]string
	volumes          *[]int
	tags             *map[string]string
	userData         *string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create an image").Action(c.action)
	c.name = command.Arg("name", "The image name").Required().String()
	c.regionID = command.Flag("region-id", "The region to launch the instance in").Required().String()
	c.zoneID = command.Flag("zone-id", "The zone to launch the instance in").Default("").String()
	c.imageID = command.Flag("image-id", "The image to launch the instance from").Required().String()
	c.serviceAccountID = command.Flag("service-account-id", "The service account to attach to the instance").Default("").String()
	c.networkID = command.Flag("network-id", "The network to attach the instance to").Required().String()
	c.flavorId = command.Flag("flavor-id", "The flavor of instance to launch").Required().String()
	c.disk = command.Flag("disk", "The size of the disk to create, this overrides the flavor.").Int()
	c.keypairIDs = command.Flag("keypair-id", "An ID of a keypair to add to the instance").Strings()
	c.volumes = command.Flag("volume", "A size of a volume to attach to the instance").Ints()
	c.tags = command.Flag("tag", "A metadata tag to add to the instance").StringMap()
	c.userData = command.Flag("user-data", "User data to add to the instance").String()
}

func (c *createCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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

	var initialVolumes []api.InstanceInitialVolume

	for _, volumeSize := range *c.volumes {
		initialVolumes = append(initialVolumes, api.InstanceInitialVolume{
			Size:       volumeSize,
			AutoDelete: true,
		})
	}

	instance, err := c.Application.APIClient.Instance().Create(*c.name, *c.imageID, *c.regionID, *c.zoneID, *c.networkID, *c.serviceAccountID, *c.flavorId, *c.disk, *c.keypairIDs, initialVolumes, tags, *c.userData)
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
