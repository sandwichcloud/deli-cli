package project

import (
	"encoding/json"
	"fmt"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type createCommand struct {
	cmd.Command
	raw  *bool
	name *string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a project").Action(c.action)
	c.name = command.Arg("name", "The project name").Required().String()
}

func (c *createCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetUnScopedToken()
	if err != nil {
		return err
	}
	project, err := c.Application.APIClient.Project().Create(*c.name)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			projectBytes, _ := json.MarshalIndent(project, "", "  ")
			fmt.Println(string(projectBytes))
		} else {
			logrus.Infof("Project '%s' created with an ID of '%s'", project.Name, project.ID)
		}
	}
	return nil
}
