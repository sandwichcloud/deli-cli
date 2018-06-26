package permission

import (
	"encoding/json"

	"net/url"

	"os"

	"fmt"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/olekukonko/tablewriter"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
)

type listCommand struct {
	cmd.Command
	limit  *int
	marker *string
	raw    *bool
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List permissions").Action(c.action)
	c.limit = command.Flag("limit", "Number of permissions to show per page").Default("20").Int()
	c.marker = command.Flag("marker", "Marker Token for the next page of results").String()
}

func (c *listCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	permissions, err := c.Application.APIClient.Permission().List(*c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			permissionsBytes, _ := json.MarshalIndent(permissions, "", "  ")
			fmt.Println(string(permissionsBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Description"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(permissions.Links) == 1 {
				nextPage := permissions.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, permission := range permissions.Permissions {
				table.Append([]string{permission.Name, permission.Description})
			}

			table.Render()
		}
	}
	return nil
}
