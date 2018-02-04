package role

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/olekukonko/tablewriter"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
)

type listCommand struct {
	cmd.Command
	project bool
	raw     *bool
	limit   *int
	marker  *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List roles").Action(c.action)
	c.limit = command.Flag("limit", "Number of projects to show per page").Default("20").Int()
	c.marker = command.Flag("marker", "Marker Token for the next page of results").String()
}

func (c *listCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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
	if c.project {
		roles, err := c.Application.APIClient.ProjectRole().ProjectList(*c.limit, *c.marker)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
				err = errors.New(apiError.ToRawJSON())
			}
			return err
		} else {
			if *c.raw {
				rolesBytes, _ := json.MarshalIndent(roles, "", "  ")
				fmt.Println(string(rolesBytes))
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "ID"})
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				if len(roles.Links) == 1 {
					nextPage := roles.Links[0]
					nextPageUrl, _ := url.Parse(nextPage.HREF)
					table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
				}

				for _, role := range roles.Roles {
					table.Append([]string{role.Name, role.ID.String()})
				}

				table.Render()
			}
		}
	} else {
		roles, err := c.Application.APIClient.GlobalRole().GlobalList(*c.limit, *c.marker)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
				err = errors.New(apiError.ToRawJSON())
			}
			return err
		} else {
			if *c.raw {
				rolesBytes, _ := json.MarshalIndent(roles, "", "  ")
				fmt.Println(string(rolesBytes))
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "ID"})
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				if len(roles.Links) == 1 {
					nextPage := roles.Links[0]
					nextPageUrl, _ := url.Parse(nextPage.HREF)
					table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
				}

				for _, role := range roles.Roles {
					table.Append([]string{role.Name, role.ID.String()})
				}

				table.Render()
			}
		}
	}
	return nil
}
