package key

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type createCommand struct {
	cmd.Command
	project          bool
	raw              *bool
	serviceAccountID *string
	keyName          *string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a key").Action(c.action)
	c.serviceAccountID = command.Arg("service account ID", "The service account ID").Required().String()
	c.keyName = command.Arg("key name", "The key name").Required().String()
}

func (c *createCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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

	var token *oauth2.Token
	if c.project {
		token, err = c.Application.APIClient.ProjectServiceAccount().CreateKey(*c.serviceAccountID, *c.keyName)
	} else {
		token, err = c.Application.APIClient.GlobalServiceAccount().CreateKey(*c.serviceAccountID, *c.keyName)
	}
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			tokenBytes, _ := json.MarshalIndent(token, "", "  ")
			fmt.Println(string(tokenBytes))
		} else {
			logrus.Infof("Key '%s' has been created on the service account '%s'", *c.keyName, *c.serviceAccountID)
			logrus.Infof("Token: %s", token.AccessToken)
		}
	}
	return nil
}
