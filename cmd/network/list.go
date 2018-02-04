package network

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
	name     *string
	regionID *string
	limit    *int
	marker   *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List networks").Action(c.action)
	c.name = command.Flag("name", "The name to filter by").Default("").String()
	c.regionID = command.Flag("regionID", "The region to filter by").Default("").String()
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
	networks, err := c.Application.APIClient.Network().List(*c.name, *c.regionID, *c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			imageBytes, _ := json.MarshalIndent(networks, "", "  ")
			fmt.Println(string(imageBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(networks.Links) == 1 {
				nextPage := networks.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, network := range networks.Networks {
				table.Append([]string{network.Name, network.ID.String()})
			}

			table.Render()
		}
	}
	return nil
}
