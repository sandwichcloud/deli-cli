package image

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
	limit      *int
	marker     *string
	visibility *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List images").Action(c.action)
	c.visibility = command.Flag("visibility", "The visibility state to filter by (PUBLIC, PRIVATE)").Default("PRIVATE").Enum("PUBLIC", "PRIVATE")
	c.limit = command.Flag("limit", "Number of projects to show per page").Default("20").Int()
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
	images, err := c.Application.APIClient.Image().List(*c.visibility, *c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *raw {
			imageBytes, _ := json.MarshalIndent(images, "", "  ")
			fmt.Println(string(imageBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(images.Links) == 1 {
				nextPage := images.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, image := range images.Images {
				table.Append([]string{image.Name, image.ID.String()})
			}

			table.Render()
		}
	}
	return nil
}
