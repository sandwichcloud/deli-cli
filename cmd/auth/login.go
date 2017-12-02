package auth

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/utils"
	log "github.com/sirupsen/logrus"
)

type loginCommand struct {
	cmd.Command
	method   *string
	username *string
	password *string
}

func (c *loginCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("login", "Login to the Sandwich Cloud API").Action(c.action)
	c.method = command.Flag("method", "Method to use for auth, if not given uses the API default.").String()
	c.username = command.Flag("username", "Username to auth with").Short('u').String()
	c.password = command.Flag("password", "User password to auth with. If not given, will prompt for input.").Short('p').String()
}

func (c *loginCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	apiDiscover, err := c.Application.APIClient.Auth().DiscoverAuth()
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	}

	var authMethod string
	if *c.method == "" {
		authMethod = *apiDiscover.Default
	} else {
		authMethod = *c.method
	}

	log.Infof("Using the %s Auth Method", authMethod)
	token, err := utils.Login(c.Application.APIClient.Auth(), *c.username, *c.password, authMethod, true)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	}

	c.Application.AuthTokens = &cmd.AuthTokens{
		Unscoped: token,
	}
	c.Application.SaveCreds()

	if *raw {
		tokenBytes, _ := json.MarshalIndent(token, "", "  ")
		fmt.Println(string(tokenBytes))
	} else {
		log.Info("You are now logged into Sandwich Cloud!")
	}

	return nil
}
