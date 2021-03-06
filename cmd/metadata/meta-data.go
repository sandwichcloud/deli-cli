package metadata

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/metadata"
)

type metaDataCommand struct {
	cmd.Command
}

func (c *metaDataCommand) Register(cmd *kingpin.CmdClause) {
	cmd.Command("meta-data", "View instance meta-data").Action(c.action)
}

func (c *metaDataCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {

	mClient := metadata.MetaDataClient{}
	err := mClient.Connect("/dev/ttyS0")
	if err != nil {
		return err
	}

	data, err := mClient.GetMetaData()
	if err != nil {
		return err
	}

	fmt.Println(data)

	return nil
}
