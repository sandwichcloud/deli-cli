package project

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
	all    *bool
	limit  *int
	marker *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List projects").Action(c.action)
	c.all = command.Flag("all", "Include projects you are not a member of").Bool()
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
	projects, err := c.Application.APIClient.Project().List(*c.all, *c.limit, *c.marker)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			projectsBytes, _ := json.MarshalIndent(projects, "", "  ")
			fmt.Println(string(projectsBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "ID"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			if len(projects.Links) == 1 {
				nextPage := projects.Links[0]
				nextPageUrl, _ := url.Parse(nextPage.HREF)
				table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
			}

			for _, project := range projects.Projects {
				table.Append([]string{project.Name, project.ID.String()})
			}

			table.Render()
		}
	}
	return nil
}
