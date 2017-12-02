package project

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
	raw       *bool
	projectID *string
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("inspect", "Inspect a project").Action(c.action)
	c.projectID = command.Arg("projectID", "The project ID").String()
}

func (c *inspectCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetUnScopedToken()
	if err != nil {
		return err
	}
	project, err := c.Application.APIClient.Project().Get(*c.projectID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			projectBytes, _ := json.MarshalIndent(project, "", "  ")
			fmt.Println(string(projectBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)

			for _, field := range structs.Fields(project) {
				table.Append([]string{field.Tag("json"), utils.InterfaceToString(field.Value())})
			}
			table.Render()
		}
	}
	return nil
}
