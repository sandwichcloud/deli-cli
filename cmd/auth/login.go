package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"encoding/json"

	"github.com/alecthomas/kingpin"
	"github.com/howeyc/gopass"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/cmd"
	"github.com/sandwichcloud/deli-cli/metadata"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type loginCommand struct {
	cmd.Command
	method   *string
	username *string
	password *string
}

func (c *loginCommand) Register(cmd *kingpin.CmdClause) {
	command := cmd.Command("login", "Login to the Sandwich Cloud API").Action(c.action)
	c.method = command.Flag("method", "Method to use for auth, if not given uses the API to auth.").Enum("manual", "metadata")
	c.username = command.Flag("username", "Username to auth with").Short('u').String()
	c.password = command.Flag("password", "User password to auth with. If not given and required, will prompt for input.").Short('p').String()
}

func (c *loginCommand) action(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	var token *oauth2.Token
	switch *c.method {
	case "manual":
		tokenString := *c.password
		for tokenString == "" {
			tokenBytes, err := gopass.GetPasswdPrompt("Please enter your token: ", true, os.Stdin, os.Stdout)
			if err != nil {
				return err
			}
			tokenString = string(tokenBytes)
		}
		token = &oauth2.Token{
			AccessToken: tokenString,
			TokenType:   "Bearer",
		}
	case "metadata":
		mClient := metadata.MetaDataClient{}
		err := mClient.Connect("/dev/ttyS0")
		if err != nil {
			return err
		}

		tokenString, err := mClient.GetSecurityData()
		if err != nil {
			return err
		}

		token = &oauth2.Token{
			AccessToken: tokenString,
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(30 * time.Minute),
		}
	default:
		var err error
		password := *c.password
		for password == "" {
			passwordBytes, err := gopass.GetPasswdPrompt("Please enter your password: ", true, os.Stdin, os.Stdout)
			if err != nil {
				return err
			}
			password = string(passwordBytes)
		}
		token, err = c.Application.APIClient.Auth().Login(*c.username, password)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok && *raw {
				err = errors.New(apiError.ToRawJSON())
			}
			return err
		}
	}

	c.Application.AuthToken = token
	c.Application.SaveCreds()

	if *raw {
		tokenBytes, _ := json.MarshalIndent(token, "", "  ")
		fmt.Println(string(tokenBytes))
	} else {
		log.Info("You are now logged into Sandwich Cloud!")
	}

	return nil
}
