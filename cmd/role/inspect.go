package role

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

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
	name    *string
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("inspect", "Inspect a role").Action(c.action)
	c.name = command.Arg("name", "The role name").String()
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	var role *api.Role
	if c.project != nil {
		role, err = c.Application.APIClient.ProjectRole(*c.project).Get(*c.name)
	} else {
		role, err = c.Application.APIClient.SystemRole().Get(*c.name)
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			roleBytes, _ := json.MarshalIndent(role, "", "  ")
			fmt.Println(string(roleBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)

			for _, field := range structs.Fields(role) {
				if field.Kind() == reflect.Slice {
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
