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
	project        *string
	raw            *bool
	name           *string
	region         *string
	zone           *string
	image          *string
	serviceAccount *string
	network        *string
	flavor         *string
	disk           *int
	keypairs       *[]string
	volumes        *[]int
	tags           *map[string]string
	userData       *string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create an image").Action(c.action)
	c.name = command.Arg("name", "The image name").Required().String()
	c.region = command.Flag("region", "The region to launch the instance in").Required().String()
	c.zone = command.Flag("zone", "The zone to launch the instance in").Default("").String()
	c.image = command.Flag("image", "The image to launch the instance from").Required().String()
	c.serviceAccount = command.Flag("service-account", "The service account to attach to the instance").Default("").String()
	c.network = command.Flag("network", "The network to attach the instance to").Required().String()
	c.flavor = command.Flag("flavor", "The flavor of instance to launch").Required().String()
	c.disk = command.Flag("disk", "The size of the disk to create, this overrides the flavor.").Int()
	c.keypairs = command.Flag("keypair", "A keypair to add to the instance").Strings()
	c.volumes = command.Flag("volume", "The size of a volume to attach to the instance").Ints()
	c.tags = command.Flag("tag", "A metadata tag to add to the instance").StringMap()
	c.userData = command.Flag("user-data", "User data to add to the instance").String()
}

func (c *createCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
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

	userData := *c.userData
	if strings.HasPrefix(userData, "@") {
		userDataBytes, err := ioutil.ReadFile(userData[1:])
		if err != nil {
			return err
		}

		userData = string(userDataBytes)
	}

	var initialVolumes []api.InstanceInitialVolume

	for _, volumeSize := range *c.volumes {
		initialVolumes = append(initialVolumes, api.InstanceInitialVolume{
			Size:       volumeSize,
			AutoDelete: true,
		})
	}

	instance, err := c.Application.APIClient.Instance(*c.project).Create(*c.name, *c.image, *c.region, *c.zone, *c.network, *c.serviceAccount, *c.flavor, *c.disk, *c.keypairs, initialVolumes, tags, userData)
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
			logrus.Infof("Instance '%s' created", instance.Name)
		}
	}
	return nil
}
