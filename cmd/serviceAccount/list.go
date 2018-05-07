package serviceAccount

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
	project bool
	raw     *bool
	limit   *int
	marker  *string
}

func (c *listCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("list", "List service accounts").Action(c.action)
	c.limit = command.Flag("limit", "Number of instances to show per page").Default("20").Int()
	c.marker = command.Flag("marker", "Marker Token for the next page of results").String()
}

func (c *listCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	if c.project {
		err = c.Application.SetScopedToken()
	} else {
		err = c.Application.SetUnScopedToken()
	}
	if err != nil {
		return err
	}
	if c.project {
		serviceAccounts, err := c.Application.APIClient.ProjectServiceAccount().ProjectList(*c.limit, *c.marker)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
				err = errors.New(apiError.ToRawJSON())
			}
			return err
		} else {
			if *c.raw {
				serviceAccountsBytes, _ := json.MarshalIndent(serviceAccounts, "", "  ")
				fmt.Println(string(serviceAccountsBytes))
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "ID"})
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				if len(serviceAccounts.Links) == 1 {
					nextPage := serviceAccounts.Links[0]
					nextPageUrl, _ := url.Parse(nextPage.HREF)
					table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
				}

				for _, serviceAccount := range serviceAccounts.ServiceAccounts {
					table.Append([]string{serviceAccount.Name, serviceAccount.ID.String()})
				}

				table.Render()
			}
		}
	} else {
		serviceAccounts, err := c.Application.APIClient.GlobalServiceAccount().GlobalList(*c.limit, *c.marker)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
				err = errors.New(apiError.ToRawJSON())
			}
			return err
		} else {
			if *c.raw {
				serviceAccountsBytes, _ := json.MarshalIndent(serviceAccounts, "", "  ")
				fmt.Println(string(serviceAccountsBytes))
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "ID"})
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				if len(serviceAccounts.Links) == 1 {
					nextPage := serviceAccounts.Links[0]
					nextPageUrl, _ := url.Parse(nextPage.HREF)
					table.SetCaption(true, fmt.Sprintf("Next Page Marker %s", nextPageUrl.Query().Get("marker")))
				}

				for _, serviceAccount := range serviceAccounts.ServiceAccounts {
					table.Append([]string{serviceAccount.Name, serviceAccount.ID.String()})
				}

				table.Render()
			}
		}
	}
	return nil
}
