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
	project          bool
	raw              *bool
	serviceAccountID *string
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("inspect", "Inspect a service account").Action(c.action)
	c.serviceAccountID = command.Arg("service account ID", "The service account ID").String()
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if c.project {
		err = c.Application.SetScopedToken()
	} else {
		err = c.Application.SetUnScopedToken()
	}
	if err != nil {
		return err
	}

	var serviceAccount *api.ServiceAccount
	if c.project {
		serviceAccount, err = c.Application.APIClient.ProjectServiceAccount().Get(*c.serviceAccountID)
	} else {
		serviceAccount, err = c.Application.APIClient.GlobalServiceAccount().Get(*c.serviceAccountID)
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
			table.SetAutoMergeCells(true)

			for _, field := range structs.Fields(serviceAccount) {
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
