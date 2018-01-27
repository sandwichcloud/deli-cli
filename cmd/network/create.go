package network

import (
	"encoding/json"
	"errors"
	"fmt"

	"net"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type createCommand struct {
	cmd.Command
	name       *string
	regionID   *string
	portGroup  *string
	cidr       *string
	gateway    *net.IP
	dnsServers *[]net.IP
	poolStart  *net.IP
	poolEnd    *net.IP
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a network").Action(c.action)
	c.name = command.Arg("name", "The network name").Required().String()
	c.regionID = command.Flag("region-id", "The region to create the network in").Required().String()
	c.portGroup = command.Flag("port-group", "The port group for the network").Required().String()
	c.cidr = command.Flag("cidr", "The network cidr").Required().String()
	c.gateway = command.Flag("gateway", "The network gateway").Required().IP()
	c.dnsServers = command.Flag("dns-server", "DNS Servers for the network").Required().IPList()
	c.poolStart = command.Flag("pool-start", "The address for the start of the IP pool").Required().IP()
	c.poolEnd = command.Flag("pool-end", "The address for the end of the IP pool").Required().IP()
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
	network, err := c.Application.APIClient.Network().Create(*c.name, *c.regionID, *c.portGroup, *c.cidr, *c.gateway, *c.poolStart, *c.poolEnd, *c.dnsServers)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			networkBytes, _ := json.MarshalIndent(network, "", "  ")
			fmt.Println(string(networkBytes))
		} else {
			logrus.Infof("Network '%s' created with an ID of '%s'", network.Name, network.ID)
		}
	}
	return nil
}
