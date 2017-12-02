package auth

import (
	"errors"
	"fmt"

	"encoding/json"

	"os"

	"github.com/alecthomas/kingpin"
	"github.com/howeyc/gopass"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client/auth"
	"github.com/sandwichcloud/deli-cli/cmd"
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

	var token *oauth2.Token

	var authMethod string
	if *c.method == "" {
		authMethod = *apiDiscover.Default
	} else {
		authMethod = *c.method
	}

	switch authMethod {
	case "github":
		log.Info("Using the Github Auth Driver")

		if apiDiscover.Github == nil {
			return errors.New("Github auth method is not enabled on this API Server.")
		}

		if *c.username == "" {
			return errors.New("Username is required for Github Authentication")
		}

		password := *c.password

		for password == "" {
			passwordBytes, err := gopass.GetPasswdPrompt("Please enter your GitHub password: ", true, os.Stdin, os.Stdout)
			if err != nil {
				return err
			}
			password = string(passwordBytes)
		}

		token, err = c.Application.APIClient.Auth().GithubLogin(*apiDiscover.Github, *c.username, password, "")
		if err == auth.ErrOTPRequired {
			var otpBytes []byte
			otpBytes, err = gopass.GetPasswdPrompt("Please enter your GitHub OTP Code: ", true, os.Stdin, os.Stdout)
			if err != nil {
				return err
			}
			otpCode := string(otpBytes)
			token, err = c.Application.APIClient.Auth().GithubLogin(*apiDiscover.Github, *c.username, password, otpCode)
		}
	default:
		return errors.New(fmt.Sprintf("Unknown API Auth Driver %s", authMethod))
	}

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
