package volume

import (
	"encoding/json"
	"errors"
	"fmt"
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
	raw      *bool
	volumeID *string
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("inspect", "Inspect a volume").Action(c.action)
	c.volumeID = command.Arg("volume ID", "The volume ID").String()
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	volume, err := c.Application.APIClient.Volume().Get(*c.volumeID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			volumeBytes, _ := json.MarshalIndent(volume, "", "  ")
			fmt.Println(string(volumeBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetAutoMergeCells(true)

			for _, field := range structs.Fields(volume) {
				if field.Kind() == reflect.Slice {
					v := reflect.ValueOf(field.Value())
					for i := 0; i < v.Len(); i++ {
						table.Append([]string{field.Tag("json"), utils.InterfaceToString(v.Index(i))})
					}
				} else if field.Kind() == reflect.Map {
					v := reflect.ValueOf(field.Value())
					for _, k := range v.MapKeys() {
						table.Append([]string{field.Tag("json"), utils.InterfaceToString(k) + "=" + utils.InterfaceToString(v.MapIndex(k))})
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
