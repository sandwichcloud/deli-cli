package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"

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
	networkID *string
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("inspect", "Inspect a network").Action(c.action)
	c.networkID = command.Arg("network ID", "The network ID").String()
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetUnScopedToken()
	if err != nil {
		return err
	}
	project, err := c.Application.APIClient.Network().Get(*c.networkID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			projectBytes, _ := json.MarshalIndent(project, "", "  ")
			fmt.Println(string(projectBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetAutoMergeCells(true)

			for _, field := range structs.Fields(project) {
				if field.Kind() == reflect.Slice && reflect.TypeOf(net.IP{}) != reflect.TypeOf(field.Value()) {
					v := reflect.ValueOf(field.Value())
					for i := 0; i < v.Len(); i++ {
						table.Append([]string{field.Tag("json"), utils.InterfaceToString(v.Index(i))})
					}
				} else {
					table.Append([]string{field.Tag("json"), utils.InterfaceToString(field.Value())})
				}
			}
			table.Render()
		}
	}
	return nil
}
