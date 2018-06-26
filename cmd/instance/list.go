package instance

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
	project *string
	raw     *bool
	image   *string
	limit   *int
	marker  *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List instances").Action(c.action)
	c.image = command.Flag("image", "The image name to filter instances by").String()
	c.limit = command.Flag("limit", "Number of instances to show per page").Default("20").Int()
	c.marker = command.Flag("marker", "Marker Token for the next page of results").String()
}

func (c *listCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	instancess, err := c.Application.APIClient.Instance(*c.project).List(*c.image, *c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			instanceBytes, _ := json.MarshalIndent(instancess, "", "  ")
			fmt.Println(string(instanceBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(instancess.Links) == 1 {
				nextPage := instancess.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, instance := range instancess.Instances {
				table.Append([]string{instance.Name})
			}

			table.Render()
		}
	}
	return nil
}
