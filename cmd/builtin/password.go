package builtin

import (
	"fmt"
	"os"

	"errors"

	"github.com/alecthomas/kingpin"
	"github.com/howeyc/gopass"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sirupsen/logrus"
)

type passwordCommand struct {
	cmd.Command
	raw      *bool
	userID   *string
	password *string
}

func (c *passwordCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("password", "Change the password").Action(c.action)
	c.userID = command.Flag("userId", "ID of the user to change. If not given changes your password.").String()
	c.password = command.Flag("password", "Password for the user. If not given, will prompt for input.").Short('p').String()
}

func (c *passwordCommand) action(app *kingpin.Application, element *kingpin.ParseElement, context *kingpin.ParseContext) error {
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

	err = c.Application.APIClient.BuiltInAuth().ChangePassword(*c.userID, password)
	if err != nil {
		if apiError, ok := err.(api.APIErrorInterface); ok && *c.raw {
			err = errors.New(apiError.ToRawJSON())
		}
		return err
	} else {
		if *c.raw {
			fmt.Println("{}")
		} else {
			logrus.Infof("Password has been updated")
		}
	}
	return nil
}
