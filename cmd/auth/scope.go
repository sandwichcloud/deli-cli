package auth

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	log "github.com/sirupsen/logrus"
)

type scopeCommand struct {
	cmd.Command
	projectID *string
}

func (c *scopeCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("scope", "Scope the current auth token to a project").Action(c.action)
	c.projectID = command.Arg("project ID", "The project ID to scope to").Required().String()
}

func (c *scopeCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	err := c.Application.LoadCreds()
	if err != nil {
		return err
	}
	err = c.Application.SetUnScopedToken()
	if err != nil {
		return err
	}

	project, err := c.Application.APIClient.Project().Get(*c.projectID)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	}

	token, err := c.Application.APIClient.Auth().ScopeToken(project)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	}

	c.Application.AuthTokens.Scoped = token
	c.Application.SaveCreds()

	if *raw {
		tokenBytes, _ := json.MarshalIndent(token, "", "  ")
		fmt.Println(string(tokenBytes))
	} else {
		log.Info("You are now scoped to the requested project.")
	}

	return nil
}
