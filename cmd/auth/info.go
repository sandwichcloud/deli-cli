package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"errors"

	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/utils"
)

type infoCommand struct {
	cmd.Command
	tokenType *string
}

func (c *infoCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("info", "Show information about the current auth tokens").Action(c.action)
	c.tokenType = command.Arg("type", "The type of token to view information from (scoped or unscoped)").Required().Enum("unscoped", "scoped")
}

func (c *infoCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}

	if *c.tokenType == "unscoped" {
		err = c.Application.SetUnScopedToken()
		if err != nil {
			return err
		}
	} else {
		err = c.Application.SetScopedToken()
		if err != nil {
			return err
		}
	}

	tokenInfo, err := c.Application.APIClient.Auth().TokenInfo()
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			projectBytes, _ := json.MarshalIndent(tokenInfo, "", "  ")
			fmt.Println(string(projectBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetAutoMergeCells(true)

			for _, field := range structs.Fields(tokenInfo) {
				name := field.Tag("json")
				if strings.HasSuffix(name, ",omitempty") {
					if utils.IsInterfaceZero(field.Value()) {
						continue
					}
					name = strings.Replace(name, ",omitempty", "", 1)
				}
				if field.Kind() == reflect.Slice {
					v := reflect.ValueOf(field.Value())
					for i := 0; i < v.Len(); i++ {
						table.Append([]string{name, utils.InterfaceToString(v.Index(i))})
					}
				} else {
					table.Append([]string{name, utils.InterfaceToString(field.Value())})
				}
			}
			table.Render()
		}
	}

	return nil
}
