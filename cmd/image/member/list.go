package member

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
	imageID *string
	raw     *bool
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List image members").Action(c.action)
	c.imageID = command.Arg("image ID", "The image ID").Required().String()
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
	imageMembers, err := c.Application.APIClient.Image().MemberList(*c.imageID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			imageMembersBytes, _ := json.MarshalIndent(imageMembers, "", "  ")
			fmt.Println(string(imageMembersBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Project ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(imageMembers.Links) == 1 {
				nextPage := imageMembers.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, imageMember := range imageMembers.Members {
				table.Append([]string{imageMember.ProjectID.String()})
			}

			table.Render()
		}
	}
	return nil
}
