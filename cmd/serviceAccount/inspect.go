package serviceAccount

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
	command := cmd.Command("inspect", "Inspect a service account").Action(c.action)
	c.name = command.Arg("name", "The service account name").String()
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	var serviceAccount *api.ServiceAccount
	if c.project != nil {
		serviceAccount, err = c.Application.APIClient.ProjectServiceAccount(*c.project).Get(*c.name)
	} else {
		serviceAccount, err = c.Application.APIClient.SystemServiceAccount().Get(*c.name)
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			serviceAccountBytes, _ := json.MarshalIndent(serviceAccount, "", "  ")
			fmt.Println(string(serviceAccountBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)

			for _, field := range structs.Fields(serviceAccount) {
				if field.Tag("json") == "keys" {
					continue
				}
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
