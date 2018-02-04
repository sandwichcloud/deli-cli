package zone

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
	regionID *string
	limit    *int
	marker   *string
	raw      *bool
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List zones").Action(c.action)
	c.regionID = command.Flag("region-id", "The region id to filter by").String()
	c.limit = command.Flag("limit", "Number of zones to show per page").Default("20").Int()
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
	zones, err := c.Application.APIClient.Zone().List(*c.regionID, *c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			zonesBytes, _ := json.MarshalIndent(zones, "", "  ")
			fmt.Println(string(zonesBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(zones.Links) == 1 {
				nextPage := zones.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, zone := range zones.Zones {
				table.Append([]string{zone.Name, zone.ID.String()})
			}

			table.Render()
		}
	}
	return nil
}
