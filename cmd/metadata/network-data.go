package metadata

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/metadata"
)

type networkDataCommand struct {
	cmd.Command
}

func (c *networkDataCommand) Register(cmd *kingpin.CmdClause) {
	cmd.Command("network-data", "View instance network-data").Action(c.action)
}

func (c *networkDataCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {

	mClient := metadata.MetaDataClient{}
	err := mClient.Connect("/dev/ttyS0")
	if err != nil {
		return err
	}

	data, err := mClient.GetNetworkData()
	if err != nil {
		return err
	}

	fmt.Println(data)

	return nil
}
