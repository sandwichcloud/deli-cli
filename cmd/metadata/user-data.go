package metadata

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/metadata"
)

type userDataCommand struct {
	cmd.Command
}

func (c *userDataCommand) Register(cmd *kingpin.CmdClause) {
	cmd.Command("user-data", "View instance user-data").Action(c.action)
}

func (c *userDataCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {

	mClient := metadata.MetaDataClient{SerialPort: "/dev/ttyS0"}

	err := mClient.Connect()
	if err != nil {
		return err
	}

	data, err := mClient.GetUserData()
	if err != nil {
		return err
	}

	fmt.Println(data)

	return nil
}
