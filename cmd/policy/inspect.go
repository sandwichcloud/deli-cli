package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/olekukonko/tablewriter"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
)

type inspectCommand struct {
	cmd.Command
	project *string
	raw     *bool
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	cmd.Command("inspect", "Inspect a policy").Action(c.action)
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	var policy *api.Policy
	if c.project != nil {
		policy, err = c.Application.APIClient.ProjectPolicy(*c.project).Get()
	} else {
		policy, err = c.Application.APIClient.SystemPolicy().Get()
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			policyBytes, _ := json.MarshalIndent(policy, "", "  ")
			fmt.Println(string(policyBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Role", "Members"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetAutoMergeCells(true)

			for _, binding := range policy.Bindings {
				for _, member := range binding.Members {
					table.Append([]string{binding.Role, member})
				}
			}
			table.Render()
		}
	}
	return nil
}
