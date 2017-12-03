package builtin

import (
	"encoding/json"
	"fmt"
	"os"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/howeyc/gopass"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type createCommand struct {
	cmd.Command
	raw      *bool
	username *string
	password *string
}

func (c *createCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("create", "Create a user").Action(c.action)
	c.username = command.Flag("username", "Username of the user").Short('u').String()
	c.password = command.Flag("password", "Password for the user. If not given, will prompt for input.").Short('p').String()
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

	password := *c.password
	for password == "" {
		passwordBytes, err := gopass.GetPasswdPrompt("Please enter a password: ", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		password = string(passwordBytes)
	}

	user, err := c.Application.APIClient.BuiltInAuth().Create(*c.username, password)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			userBytes, _ := json.MarshalIndent(user, "", "  ")
			fmt.Println(string(userBytes))
		} else {
			logrus.Infof("User '%s' created with an ID of '%s'", user.Username, user.ID)
		}
	}
	return nil
}
