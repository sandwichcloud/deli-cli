package member

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
	raw    *bool
	limit  *int
	marker *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List project members").Action(c.action)
	c.limit = command.Flag("limit", "Number of project memberss to show per page").Default("20").Int()
	c.marker = command.Flag("marker", "Marker Token for the next page of results").String()
}

func (c *listCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetScopedToken()
	if err != nil {
		return err
	}
	projectMembers, err := c.Application.APIClient.Project().ListMembers(*c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			projectsBytes, _ := json.MarshalIndent(projectMembers, "", "  ")
			fmt.Println(string(projectsBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Driver", "ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(projectMembers.Links) == 1 {
				nextPage := projectMembers.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, projectMember := range projectMembers.ProjectMemberss {
				table.Append([]string{projectMember.Username, projectMember.Driver, projectMember.ID.String()})
			}

			table.Render()
		}
	}
	return nil
}
