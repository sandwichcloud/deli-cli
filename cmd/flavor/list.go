package flavor

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
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List flavors").Action(c.action)
	c.limit = command.Flag("limit", "Number of projects to show per page").Default("20").Int()
	c.marker = command.Flag("marker", "Marker Token for the next page of results").String()
}

func (c *listCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetUnScopedToken()
	if err != nil {
		return err
	}
	flavors, err := c.Application.APIClient.Flavor().List(*c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			flavorsBytes, _ := json.MarshalIndent(flavors, "", "  ")
			fmt.Println(string(flavorsBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(flavors.Links) == 1 {
				nextPage := flavors.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, flavor := range flavors.Flavors {
				table.Append([]string{flavor.Name, flavor.ID.String()})
			}

			table.Render()
		}
	}
	return nil
}
