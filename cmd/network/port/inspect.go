package port

import (
	"encoding/json"
	"errors"
	"fmt"

	"os"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/utils"
)

type inspectCommand struct {
	cmd.Command
	project *string
	raw     *bool
	portID  *string
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("inspect", "Inspect a network port").Action(c.action)
	c.portID = command.Arg("id", "The network port ID").Required().String()
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	keypair, err := c.Application.APIClient.NetworkPort(*c.project).Get(*c.portID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			keypairBytes, _ := json.MarshalIndent(keypair, "", "  ")
			fmt.Println(string(keypairBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)

			for _, field := range structs.Fields(keypair) {
				table.Append([]string{field.Tag("json"), utils.InterfaceToString(field.Value())})
			}
			table.Render()
		}
	}
	return nil
}
