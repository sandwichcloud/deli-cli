package image

import (
	"encoding/json"
	"errors"
	"fmt"

	"os"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/utils"
)

type inspectCommand struct {
	cmd.Command
	name    *string
	project *string
	raw     *bool
}

func (c *inspectCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("inspect", "Inspect an image").Action(c.action)
	c.name = command.Arg("image name", "The image name").String()
}

func (c *inspectCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	image, err := c.Application.APIClient.Image(*c.project).Get(*c.name)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			imageBytes, _ := json.MarshalIndent(image, "", "  ")
			fmt.Println(string(imageBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Property", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)

			for _, field := range structs.Fields(image) {
				table.Append([]string{field.Tag("json"), utils.InterfaceToString(field.Value())})
			}
			table.Render()
		}
	}
	return nil
}
