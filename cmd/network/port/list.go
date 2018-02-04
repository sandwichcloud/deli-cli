package port

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
	raw    *bool
	limit  *int
	marker *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List network ports").Action(c.action)
	c.limit = command.Flag("limit", "Number of network ports to show per page").Default("20").Int()
	c.marker = command.Flag("marker", "Marker Token for the next page of results").String()
}

func (c *listCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	networkPorts, err := c.Application.APIClient.NetworkPort().List(*c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			keypairsBytes, _ := json.MarshalIndent(networkPorts, "", "  ")
			fmt.Println(string(keypairsBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(networkPorts.Links) == 1 {
				nextPage := networkPorts.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, networkPort := range networkPorts.NetworkPorts {
				table.Append([]string{networkPort.ID.String()})
			}

			table.Render()
		}
	}
	return nil
}
